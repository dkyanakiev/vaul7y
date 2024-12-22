# Usage 

## Setup
### Installation

### Brew

```shell
brew install dkyanakiev/tap/vaul7y
```
to upgrade
```shell
brew update && brew upgrade vaul7y
```

### Download from GitHub

Download the relevant binary for your operating system (macOS = Darwin) from
the [latest Github release](https://github.com/dkyanakiev/vaul7y/releases). Unpack it, then move the binary to
somewhere accessible in your `PATH`, e.g. `mv ./vaul7y /usr/local/bin`.

### > Using [go installed on your machine](https://go.dev/doc/install)

```shell
go install github.com/dkyanakiev/vaul7y@latest
```

### Building from source and Run Vaul7y

Make sure you have your go environment setup:

1. Clone the project
1. Run `$ make build` to build the binary
1. Run `$ make run` to run the binary
1. You can use `$ make install-osx` on a Mac to cp the binary to `/usr/local/bin/vaul7y`

or

```
$ go install ./cmd/vaul7y
```

### How to use it

Once `Vaul7y` is installed and avialable in your path, simply run:

```
$ vaul7y
```

![image](../images/screen1.png)


### Environment variables

In order to use the tool you must expose the needed env variables, that would generally be used by the vault cli to auth to a given cluster. 

Required:  
`VAULT_ADDR`  
`VAULT_TOKEN`

For the full list see the [official docs](https://developer.hashicorp.com/vault/docs/commands#environment-variables)

Another option is to store your configs in yaml file named `.vaul7y.yaml` stored in your home directory.  
Example: [`~/myuser/.vaul7y.yaml`](./examples/vaul7y.yaml)

Or alternatively pass a config file as an argument using `-c <path/file.yaml>`  
Example: `vaul7y -c ./new-env.yml`

### Features

Currently the capabilities are limited. 

* Support for navigation between KV mounts
    * Currently only KV2
* Looking up secret objects
    * Show/hide secrets and coping data
    * Update/patch secrets
    * Create new secrets
    * Filter paths/secrets 
* Support for exploring and filtering ACL Policies
* Namespace support for Enteprise versions
