## component-cli component-archive signature sign rsa

fetch the component descriptor from an oci registry and sign it using RSASSA-PKCS1-V1_5-SIGN

```
component-cli component-archive signature sign rsa BASE_URL COMPONENT_NAME VERSION [flags]
```

### Options

```
      --allow-plain-http            allows the fallback to http if the oci registry does not support https
      --cc-config string            path to the local concourse config file
      --force                       force overwrite of already existing component descriptors
  -h, --help                        help for rsa
      --insecure-skip-tls-verify    If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --private-key string          path to private key file used for signing
      --recursive                   recursively sign and upload all referenced component descriptors
      --registry-config string      path to the dockerconfig.json with the oci registry authentication information
      --signature-name string       name of the signature
      --skip-access-types strings   comma separated list of access types that will not be digested and signed
      --upload-base-url string      target repository context to upload the signed cd
```

### Options inherited from parent commands

```
      --cli                  logger runs as cli logger. enables cli logging
      --dev                  enable development logging which result in console encoding, enabled stacktrace and enabled caller
      --disable-caller       disable the caller of logs (default true)
      --disable-stacktrace   disable the stacktrace of error logs (default true)
      --disable-timestamp    disable timestamp output (default true)
  -v, --verbosity int        number for the log level verbosity (default 1)
```

### SEE ALSO

* [component-cli component-archive signature sign](component-cli_component-archive_signature_sign.md)	 - command to sign component descriptors

