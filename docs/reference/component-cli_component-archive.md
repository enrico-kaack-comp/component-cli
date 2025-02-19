## component-cli component-archive



```
component-cli component-archive [component-archive-path] [ctf-path] [flags]
```

### Options

```
  -a, --archive string                  path to the component archive directory
      --component-name string           name of the component
      --component-name-mapping string   [OPTIONAL] repository context name mapping (default "urlPath")
  -c, --component-ref stringArray       path to component references definition
      --component-version string        version of the component
      --format CAOutputFormat           archive format of the component archive. Can be "tar" or "tgz" (default tar)
  -h, --help                            help for component-archive
      --repo-ctx string                 [OPTIONAL] repository context url for component to upload. The repository url will be automatically added to the repository contexts.
  -r, --resources stringArray           path to resources definition
  -s, --sources stringArray             path to sources definition
      --temp-dir string                 temporary directory where the component archive is build. Defaults to a os-specific temp dir
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

* [component-cli](component-cli.md)	 - component cli
* [component-cli component-archive component-references](component-cli_component-archive_component-references.md)	 - command to modify component references of a component descriptor
* [component-cli component-archive create](component-cli_component-archive_create.md)	 - Creates a component archive with a component descriptor
* [component-cli component-archive export](component-cli_component-archive_export.md)	 - Exports a component archive as defined by CTF
* [component-cli component-archive remote](component-cli_component-archive_remote.md)	 - command to interact with component descriptors stored in an oci registry
* [component-cli component-archive resources](component-cli_component-archive_resources.md)	 - command to modify resources of a component descriptor
* [component-cli component-archive signature](component-cli_component-archive_signature.md)	 - [EXPERIMENTAL] command to work with signatures and digests in component descriptors
* [component-cli component-archive sources](component-cli_component-archive_sources.md)	 - command to modify sources of a component descriptor

