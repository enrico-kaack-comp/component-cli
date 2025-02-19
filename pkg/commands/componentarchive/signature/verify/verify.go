package verify

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	cdv2Sign "github.com/gardener/component-spec/bindings-go/apis/v2/signatures"
	cdoci "github.com/gardener/component-spec/bindings-go/oci"
	"github.com/go-logr/logr"
	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/gardener/component-cli/ociclient"
	ociopts "github.com/gardener/component-cli/ociclient/options"
	"github.com/gardener/component-cli/pkg/commands/constants"
	"github.com/gardener/component-cli/pkg/logger"
	"github.com/gardener/component-cli/pkg/signatures"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewVerifyCommand creates a new command to verify signatures.
func NewVerifyCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "command to verify the signature of a component descriptor",
	}

	cmd.AddCommand(NewRSAVerifyCommand(ctx))
	return cmd
}

type GenericVerifyOptions struct {
	// BaseUrl is the oci registry where the component is stored.
	BaseUrl string
	// ComponentName is the unique name of the component in the registry.
	ComponentName string
	// Version is the component version in the oci registry.
	Version string

	// SignatureName selects the matching signature to verify
	SignatureName string

	// SkipAccessTypes defines the access types that will be ignored for verification
	SkipAccessTypes []string

	// OciOptions contains all exposed options to configure the oci client.
	OciOptions ociopts.Options
}

//Complete validates the arguments and flags from the command line
func (o *GenericVerifyOptions) Complete(args []string) error {
	o.BaseUrl = args[0]
	o.ComponentName = args[1]
	o.Version = args[2]

	cliHomeDir, err := constants.CliHomeDir()
	if err != nil {
		return err
	}

	o.OciOptions.CacheDir = filepath.Join(cliHomeDir, "components")
	if err := os.MkdirAll(o.OciOptions.CacheDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create cache directory %s: %w", o.OciOptions.CacheDir, err)
	}

	if len(o.BaseUrl) == 0 {
		return errors.New("the base url must be defined")
	}
	if len(o.ComponentName) == 0 {
		return errors.New("a component name must be defined")
	}
	if len(o.Version) == 0 {
		return errors.New("a component's version must be defined")
	}
	if o.SignatureName == "" {
		return errors.New("a signature name must be provided")
	}
	return nil
}

func (o *GenericVerifyOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.SignatureName, "signature-name", "", "name of the signature to verify")
	fs.StringSliceVar(&o.SkipAccessTypes, "skip-access-types", []string{}, "comma separated list of access types that will be ignored for verification")
	o.OciOptions.AddFlags(fs)
}

func (o *GenericVerifyOptions) VerifyWithVerifier(ctx context.Context, log logr.Logger, fs vfs.FileSystem, verifier cdv2Sign.Verifier) error {
	repoCtx := cdv2.NewOCIRegistryRepository(o.BaseUrl, "")

	ociClient, _, err := o.OciOptions.Build(log, fs)
	if err != nil {
		return fmt.Errorf("unable to build oci client: %s", err.Error())
	}

	cdresolver := cdoci.NewResolver(ociClient)
	cd, err := cdresolver.Resolve(ctx, repoCtx, o.ComponentName, o.Version)
	if err != nil {
		return fmt.Errorf("unable to to fetch component descriptor %s:%s: %w", o.ComponentName, o.Version, err)
	}

	// check componentReferences and resources
	if err := CheckCdDigests(cd, *repoCtx, ociClient, context.TODO(), o.SkipAccessTypes); err != nil {
		return fmt.Errorf("failed checking cd: %w", err)
	}

	// check if digest is correctly signed
	if err = cdv2Sign.VerifySignedComponentDescriptor(cd, verifier, o.SignatureName); err != nil {
		return fmt.Errorf("signature invalid for digest: %w", err)
	}

	// check if digest matches the normalised component descriptor
	hasher, err := cdv2Sign.HasherForName(cdv2Sign.SHA256)
	if err != nil {
		return fmt.Errorf("failed creating hasher: %w", err)
	}
	hashCd, err := cdv2Sign.HashForComponentDescriptor(*cd, *hasher)
	if err != nil {
		return fmt.Errorf("failed hashing cd %s:%s: %w", cd.Name, cd.Version, err)
	}

	matchingSignature, err := cdv2Sign.SelectSignatureByName(cd, o.SignatureName)
	if err != nil {
		return fmt.Errorf("failed selecting signature %s: %w", o.SignatureName, err)
	}

	if hashCd.HashAlgorithm != matchingSignature.Digest.HashAlgorithm || hashCd.NormalisationAlgorithm != matchingSignature.Digest.NormalisationAlgorithm || hashCd.Value != matchingSignature.Digest.Value {
		return fmt.Errorf("failed verifiying signature: signed normalised digest does not match calculated digest")
	}

	log.Info(fmt.Sprintf("Signature %s is valid and digest of normalised cd matches calculated digest", o.SignatureName))
	return nil

}

func CheckCdDigests(cd *cdv2.ComponentDescriptor, repoContext cdv2.OCIRegistryRepository, ociClient ociclient.Client, ctx context.Context, skipAccessTypes []string) error {
	skipAccessTypesMap := map[string]bool{}
	for _, v := range skipAccessTypes {
		skipAccessTypesMap[v] = true
	}
	for _, reference := range cd.ComponentReferences {
		ociRef, err := cdoci.OCIRef(repoContext, reference.Name, reference.Version)
		if err != nil {
			return fmt.Errorf("invalid component reference: %w", err)
		}

		cdresolver := cdoci.NewResolver(ociClient)
		childCd, err := cdresolver.Resolve(ctx, &repoContext, reference.ComponentName, reference.Version)
		if err != nil {
			return fmt.Errorf("unable to to fetch component descriptor %s: %w", ociRef, err)
		}

		if reference.Digest == nil || reference.Digest.HashAlgorithm == "" || reference.Digest.NormalisationAlgorithm == "" || reference.Digest.Value == "" {
			return fmt.Errorf("component reference is missing digest %s:%s", reference.ComponentName, reference.Version)
		}

		hasherForCdReference, err := cdv2Sign.HasherForName(reference.Digest.HashAlgorithm)
		if err != nil {
			return fmt.Errorf("failed creating hasher for algorithm %s for referenceCd %s %s: %w", reference.Digest.HashAlgorithm, reference.Name, reference.Version, err)
		}

		digest, err := recursivelyCheckCdsDigests(childCd, repoContext, ociClient, ctx, hasherForCdReference, skipAccessTypes)
		if err != nil {
			return fmt.Errorf("checking of component reference %s:%s failed: %w", reference.ComponentName, reference.Version, err)
		}

		if !reflect.DeepEqual(reference.Digest, digest) {
			return fmt.Errorf("component reference digest for  %s:%s is different to stored one", reference.ComponentName, reference.Version)
		}

	}
	for _, resource := range cd.Resources {
		log := logger.Log.WithValues("componentDescriptor", cd, "resource.name", resource.Name, "resource.version", resource.Version, "resource.extraIdentity", resource.ExtraIdentity)

		//skip ignored access type
		if _, ok := skipAccessTypesMap[resource.Access.Type]; ok {
			log.Info("skipping resource as defined in --skip-access-types")
			continue
		}
		if resource.Digest == nil || resource.Digest.HashAlgorithm == "" || resource.Digest.NormalisationAlgorithm == "" || resource.Digest.Value == "" {
			return fmt.Errorf("resource is missing digest %s:%s", resource.Name, resource.Version)
		}

		hasher, err := cdv2Sign.HasherForName(resource.Digest.HashAlgorithm)
		if err != nil {
			return fmt.Errorf("failed creating hasher for algorithm %s for resource %s %s: %w", resource.Digest.HashAlgorithm, resource.Name, resource.Version, err)
		}
		digester := signatures.NewDigester(ociClient, *hasher, skipAccessTypes)

		digest, err := digester.DigestForResource(ctx, *cd, resource)
		if err != nil {
			return fmt.Errorf("failed creating digest for resource %s: %w", resource.Name, err)
		}

		if !reflect.DeepEqual(resource.Digest, digest) {
			return fmt.Errorf("resource digest is different to stored one %s:%s", resource.Name, resource.Version)
		}

	}
	return nil
}

func recursivelyCheckCdsDigests(cd *cdv2.ComponentDescriptor, repoContext cdv2.OCIRegistryRepository, ociClient ociclient.Client, ctx context.Context, hasherForCd *cdv2Sign.Hasher, skipAccessTypes []string) (*cdv2.DigestSpec, error) {
	skipAccessTypesMap := map[string]bool{}
	for _, v := range skipAccessTypes {
		skipAccessTypesMap[v] = true
	}

	for referenceIndex, reference := range cd.ComponentReferences {
		reference := reference

		ociRef, err := cdoci.OCIRef(repoContext, reference.Name, reference.Version)
		if err != nil {
			return nil, fmt.Errorf("invalid component reference: %w", err)
		}

		cdresolver := cdoci.NewResolver(ociClient)
		childCd, err := cdresolver.Resolve(ctx, &repoContext, reference.ComponentName, reference.Version)
		if err != nil {
			return nil, fmt.Errorf("unable to to fetch component descriptor %s: %w", ociRef, err)
		}

		hasher, err := cdv2Sign.HasherForName(cdv2Sign.SHA256)
		if err != nil {
			return nil, fmt.Errorf("failed creating hasher: %w", err)
		}

		digest, err := recursivelyCheckCdsDigests(childCd, repoContext, ociClient, ctx, hasher, skipAccessTypes)
		if err != nil {
			return nil, fmt.Errorf("unable to resolve component reference to %s:%s: %w", reference.ComponentName, reference.Version, err)
		}
		reference.Digest = digest
		cd.ComponentReferences[referenceIndex] = reference
	}
	for resourceIndex, resource := range cd.Resources {
		resource := resource
		log := logger.Log.WithValues("componentDescriptor", cd, "resource.name", resource.Name, "resource.version", resource.Version, "resource.extraIdentity", resource.ExtraIdentity)

		//skip ignored access type
		if _, ok := skipAccessTypesMap[resource.Access.Type]; ok {
			log.Info("skipping resource as defined in --skip-access-types")
			continue
		}

		hasher, err := cdv2Sign.HasherForName(cdv2Sign.SHA256)
		if err != nil {
			return nil, fmt.Errorf("failed creating hasher: %w", err)
		}

		digester := signatures.NewDigester(ociClient, *hasher, skipAccessTypes)

		digest, err := digester.DigestForResource(ctx, *cd, resource)
		if err != nil {
			return nil, fmt.Errorf("failed creating digest for resource %s: %w", resource.Name, err)
		}

		// For better user information, log resource with mismatching digest.
		// Since we do not trust the digest data in this cd, it is only for information purpose.
		// The mismatch will be noted in the propagated cd reference digest in the root cd.
		if resource.Digest != nil && !reflect.DeepEqual(resource.Digest, digest) {
			log.Info(fmt.Sprintf("digest in (untrusted) cd %+v mismatches with calculated digest %+v ", resource.Digest, digest))
		}

		resource.Digest = digest
		cd.Resources[resourceIndex] = resource
	}

	hashCd, err := cdv2Sign.HashForComponentDescriptor(*cd, *hasherForCd)
	if err != nil {
		return nil, fmt.Errorf("failed hashing cd %s:%s: %w", cd.Name, cd.Version, err)
	}
	return hashCd, nil
}
