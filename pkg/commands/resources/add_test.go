// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/codec"
	"github.com/gardener/component-spec/bindings-go/ctf"
	testlog "github.com/go-logr/logr/testing"
	"github.com/mandelsoft/vfs/pkg/layerfs"
	"github.com/mandelsoft/vfs/pkg/memoryfs"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/projectionfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/gardener/component-cli/pkg/commands/resources"
	"github.com/gardener/component-cli/pkg/utils"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resources Test Suite")
}

var _ = Describe("Add", func() {

	var testdataFs vfs.FileSystem

	BeforeEach(func() {
		fs, err := projectionfs.New(osfs.New(), "./testdata")
		Expect(err).ToNot(HaveOccurred())
		testdataFs = layerfs.New(memoryfs.New(), fs)
	})

	It("should add a resource defined by a file", func() {
		opts := &resources.Options{
			ComponentArchivePath: "./00-component",
			ResourceObjectPath:   "./resources/00-res.yaml",
		}

		Expect(opts.Run(context.TODO(), testlog.NullLogger{}, testdataFs)).To(Succeed())

		data, err := vfs.ReadFile(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.ComponentDescriptorFileName))
		Expect(err).ToNot(HaveOccurred())

		cd := &cdv2.ComponentDescriptor{}
		Expect(codec.Decode(data, cd)).To(Succeed())

		Expect(cd.Resources).To(HaveLen(1))
		Expect(cd.Resources[0].IdentityObjectMeta).To(MatchFields(IgnoreExtras, Fields{
			"Name":    Equal("ubuntu"),
			"Version": Equal("v0.0.1"),
			"Type":    Equal("ociImage"),
		}))
		Expect(cd.Resources[0]).To(MatchFields(IgnoreExtras, Fields{
			"Relation": Equal(cdv2.ResourceRelation("external")),
		}))
		Expect(cd.Resources[0].Access.Object).To(HaveKeyWithValue("type", "ociRegistry"))
		Expect(cd.Resources[0].Access.Object).To(HaveKeyWithValue("imageReference", "ubuntu:18.0"))
	})

	It("should add a resource defined by stdin", func() {
		oldstdin := os.Stdin
		defer func() {
			os.Stdin = oldstdin
		}()
		r, w, err := os.Pipe()
		Expect(err).ToNot(HaveOccurred())
		_, err = w.WriteString(`
name: 'ubuntu'
version: 'v0.0.1'
type: 'ociImage'
relation: 'external'
access:
  type: 'ociRegistry'
  imageReference: 'ubuntu:18.0'
`)
		Expect(err).ToNot(HaveOccurred())
		Expect(w.Close()).To(Succeed())
		os.Stdin = r

		opts := &resources.Options{
			ComponentArchivePath: "./00-component",
		}

		Expect(opts.Run(context.TODO(), testlog.NullLogger{}, testdataFs)).To(Succeed())

		data, err := vfs.ReadFile(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.ComponentDescriptorFileName))
		Expect(err).ToNot(HaveOccurred())

		cd := &cdv2.ComponentDescriptor{}
		Expect(codec.Decode(data, cd)).To(Succeed())

		Expect(cd.Resources).To(HaveLen(1))
		Expect(cd.Resources[0].IdentityObjectMeta).To(MatchFields(IgnoreExtras, Fields{
			"Name":    Equal("ubuntu"),
			"Version": Equal("v0.0.1"),
			"Type":    Equal("ociImage"),
		}))
		Expect(cd.Resources[0]).To(MatchFields(IgnoreExtras, Fields{
			"Relation": Equal(cdv2.ResourceRelation("external")),
		}))
		Expect(cd.Resources[0].Access.Object).To(HaveKeyWithValue("type", "ociRegistry"))
		Expect(cd.Resources[0].Access.Object).To(HaveKeyWithValue("imageReference", "ubuntu:18.0"))
	})

	It("should automatically set the version for a local resource", func() {
		opts := &resources.Options{
			ComponentArchivePath: "./00-component",
			ResourceObjectPath:   "./resources/01-local.yaml",
		}

		Expect(opts.Run(context.TODO(), testlog.NullLogger{}, testdataFs)).To(Succeed())

		data, err := vfs.ReadFile(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.ComponentDescriptorFileName))
		Expect(err).ToNot(HaveOccurred())

		cd := &cdv2.ComponentDescriptor{}
		Expect(codec.Decode(data, cd)).To(Succeed())

		Expect(cd.Resources).To(HaveLen(1))
		Expect(cd.Resources[0].IdentityObjectMeta).To(MatchFields(IgnoreExtras, Fields{
			"Name":    Equal("testres"),
			"Version": Equal("v0.0.0"),
			"Type":    Equal("mytype"),
		}))
	})

	It("should add multiple resources via multi yaml docs", func() {
		opts := &resources.Options{
			ComponentArchivePath: "./00-component",
			ResourceObjectPath:   "./resources/02-multidoc.yaml",
		}

		Expect(opts.Run(context.TODO(), testlog.NullLogger{}, testdataFs)).To(Succeed())

		data, err := vfs.ReadFile(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.ComponentDescriptorFileName))
		Expect(err).ToNot(HaveOccurred())

		cd := &cdv2.ComponentDescriptor{}
		Expect(codec.Decode(data, cd)).To(Succeed())

		Expect(cd.Resources).To(HaveLen(2))
		Expect(cd.Resources[0].IdentityObjectMeta).To(MatchFields(IgnoreExtras, Fields{
			"Name":    Equal("ubuntu"),
			"Version": Equal("v0.0.1"),
			"Type":    Equal("ociImage"),
		}))
		Expect(cd.Resources[1].IdentityObjectMeta).To(MatchFields(IgnoreExtras, Fields{
			"Name":    Equal("testres"),
			"Version": Equal("v0.0.0"),
			"Type":    Equal("mytype"),
		}))
	})

	It("should throw an error if an invalid resource is defined", func() {
		opts := &resources.Options{
			ComponentArchivePath: "./00-component",
			ResourceObjectPath:   "./resources/10-res-invalid.yaml",
		}

		Expect(opts.Run(context.TODO(), testlog.NullLogger{}, testdataFs)).To(HaveOccurred())

		data, err := vfs.ReadFile(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.ComponentDescriptorFileName))
		Expect(err).ToNot(HaveOccurred())
		cd := &cdv2.ComponentDescriptor{}
		Expect(codec.Decode(data, cd)).To(Succeed())
		Expect(cd.Resources).To(HaveLen(0))
	})

	Context("With Input", func() {
		It("should add a resource defined by a file with a jsonfile input", func() {
			opts := &resources.Options{
				ComponentArchivePath: "./00-component",
				// jsonschema example copied from https://json-schema.org/learn/miscellaneous-examples.html
				ResourceObjectPath: "./resources/20-res-json.yaml",
			}

			Expect(opts.Run(context.TODO(), testlog.NullLogger{}, testdataFs)).To(Succeed())

			data, err := vfs.ReadFile(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.ComponentDescriptorFileName))
			Expect(err).ToNot(HaveOccurred())
			cd := &cdv2.ComponentDescriptor{}
			Expect(codec.Decode(data, cd)).To(Succeed())

			Expect(cd.Resources).To(HaveLen(1))
			Expect(cd.Resources[0].IdentityObjectMeta).To(MatchFields(IgnoreExtras, Fields{
				"Name":    Equal("myconfig"),
				"Version": Equal("v0.0.1"),
				"Type":    Equal("jsonschema"),
			}))
			Expect(cd.Resources[0]).To(MatchFields(IgnoreExtras, Fields{
				"Relation": Equal(cdv2.ResourceRelation("external")),
			}))
			Expect(cd.Resources[0].Access.Object).To(HaveKeyWithValue("type", cdv2.LocalFilesystemBlobType))
			Expect(cd.Resources[0].Access.Object).To(HaveKeyWithValue("filename", BeAssignableToTypeOf("")))

			blobs, err := vfs.ReadDir(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.BlobsDirectoryName))
			Expect(err).ToNot(HaveOccurred())
			Expect(blobs).To(HaveLen(1))
		})

		It("should automatically tar a directory input and add it as resource", func() {
			opts := &resources.Options{
				ComponentArchivePath: "./00-component",
				ResourceObjectPath:   "./resources/20-res-json.yaml",
			}

			Expect(opts.Run(context.TODO(), testlog.NullLogger{}, testdataFs)).To(Succeed())

			data, err := vfs.ReadFile(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.ComponentDescriptorFileName))
			Expect(err).ToNot(HaveOccurred())
			cd := &cdv2.ComponentDescriptor{}
			Expect(codec.Decode(data, cd)).To(Succeed())

			Expect(cd.Resources).To(HaveLen(1))
			Expect(cd.Resources[0].IdentityObjectMeta).To(MatchFields(IgnoreExtras, Fields{
				"Name":    Equal("myconfig"),
				"Version": Equal("v0.0.1"),
				"Type":    Equal("jsonschema"),
			}))
			Expect(cd.Resources[0]).To(MatchFields(IgnoreExtras, Fields{
				"Relation": Equal(cdv2.ResourceRelation("external")),
			}))
			Expect(cd.Resources[0].Access.Object).To(HaveKeyWithValue("type", cdv2.LocalFilesystemBlobType))
			Expect(cd.Resources[0].Access.Object).To(HaveKeyWithValue("filename", BeAssignableToTypeOf("")))

			blobs, err := vfs.ReadDir(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.BlobsDirectoryName))
			Expect(err).ToNot(HaveOccurred())
			Expect(blobs).To(HaveLen(1))
		})

		It("should gzip a input blob and add it as resource if the gzip flag is provided", func() {
			opts := &resources.Options{
				ComponentArchivePath: "./00-component",
				ResourceObjectPath:   "./resources/21-res-dir-zip.yaml",
			}

			Expect(opts.Run(context.TODO(), testlog.NullLogger{}, testdataFs)).To(Succeed())

			data, err := vfs.ReadFile(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.ComponentDescriptorFileName))
			Expect(err).ToNot(HaveOccurred())
			cd := &cdv2.ComponentDescriptor{}
			Expect(codec.Decode(data, cd)).To(Succeed())

			Expect(cd.Resources).To(HaveLen(1))
			Expect(cd.Resources[0].IdentityObjectMeta).To(MatchFields(IgnoreExtras, Fields{
				"Name":    Equal("myconfig"),
				"Version": Equal("v0.0.1"),
				"Type":    Equal("jsonschema"),
			}))
			Expect(cd.Resources[0]).To(MatchFields(IgnoreExtras, Fields{
				"Relation": Equal(cdv2.ResourceRelation("external")),
			}))
			Expect(cd.Resources[0].Access.Object).To(HaveKeyWithValue("type", cdv2.LocalFilesystemBlobType))
			Expect(cd.Resources[0].Access.Object).To(HaveKeyWithValue("filename", BeAssignableToTypeOf("")))

			blobs, err := vfs.ReadDir(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.BlobsDirectoryName))
			Expect(err).ToNot(HaveOccurred())
			Expect(blobs).To(HaveLen(1))

			mimetype, err := utils.GetFileType(testdataFs, filepath.Join(opts.ComponentArchivePath, ctf.BlobsDirectoryName, blobs[0].Name()))
			Expect(err).ToNot(HaveOccurred())
			Expect(mimetype).To(Equal("application/x-gzip"))
		})

	})

})
