# envii

```
envii streamlines the management and synchronization of .env files among different environments or team members, facilitating consistent configuration across development workflows.

Usage:
  envii [command] [flags]
  envii [command]

Available Commands:
  auth        Display commands for authenticating envii with an account
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  pull        Pulls the latest .env file from the remote server and replaces the local .env file with it.
  push        Pushes the latest .env file to the remote server.
  undo        undoes the last `envii pull` command

Flags:
  -h, --help                         help for envii
  -k, --k8s-values-path string       Path to the k8s values file
  -p, --provider string              Provider to use to pull the .env file: zipper | k8s (default "zipper")
  -s, --secrets-declaration string   Path or identifier for the secrets declaration
  -v, --version                      version for envii

Use "envii [command] --help" for more information about a command.
```
