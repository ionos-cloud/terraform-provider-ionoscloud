# Introduction

## IONOS Cloud Terraform Provider

The IonosCloud provider gives the ability to deploy and configure resources using the IonosCloud APIs.

### Migrating from the ProfitBricks provider

Please see the [Documentation](../docs/index.md#migrating-from-the-profitbricks-provider) on how to migrate from the ProfitBricks provider.

### Requirements

* [Terraform](https://www.terraform.io/downloads.html) 0.12.x
* [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

**NOTE:** In order to use a specific version of this provider, please include the following block at the beginning of your terraform config files [details](https://www.terraform.io/docs/configuration/terraform.html#specifying-a-required-terraform-version):

```
provider "ionoscloud" {
  version = "~> 6.2.0"
}
```

### Building The Provider

Clone repository to: `$GOPATH/src/github.com/ionos-cloud/terraform-provider-ionoscloud`

```
$ mkdir -p $GOPATH/src/github.com/ionos-cloud; cd $GOPATH/src/github.com/ionos-cloud
$ git clone git@github.com:ionos-cloud/terraform-provider-ionoscloud
```

Enter the provider directory and build the provider

```
$ cd $GOPATH/src/github.com/ionos-cloud/terraform-provider-ionoscloud
$ make build
```

### Using the provider

See the [IonosCloud Provider documentation](https://registry.terraform.io/providers/ionos-cloud/ionoscloud/latest/docs) to get started using the IonosCloud provider.

### Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is _required_). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```
$ make build
...
$ $GOPATH/bin/terraform-provider-ionoscloud
...
```

In order to test the provider, you can simply run `make test`.

```
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```
$ make testacc
```

## Certificate pinning:

You can enable certificate pinning if you want to bypass the normal certificate checking procedure,
by doing the following:

Set env variable IONOS_PINNED_CERT=<insert_sha256_public_fingerprint_here>

You can get the sha256 fingerprint most easily form the browser by inspecting the certificate.