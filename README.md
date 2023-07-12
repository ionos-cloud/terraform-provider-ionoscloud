
[![Gitter](https://img.shields.io/gitter/room/ionos-cloud/sdk-general)](https://gitter.im/ionos-cloud/sdk-general)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=alert_status)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=bugs)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=security_rating)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Release](https://img.shields.io/github/v/release/ionos-cloud/terraform-provider-ionoscloud.svg)](https://github.com/ionos-cloud/terraform-provider-ionoscloud/releases/latest)
[![Release Date](https://img.shields.io/github/release-date/ionos-cloud/terraform-provider-ionoscloud.svg)](https://github.com/ionos-cloud/terraform-provider-ionoscloud/releases/latest)
[![Go](https://img.shields.io/github/go-mod/go-version/ionos-cloud/terraform-provider-ionoscloud.svg)](https://github.com/ionos-cloud/terraform-provider-ionoscloud)

![Alt text](.github/IONOS.CLOUD.BLU.svg?raw=true "Title")


**IMPORTANT NOTE**: 

Terraform IONOS Cloud Provider v5 is deprecated and no longer maintained. Please upgrade to v6, which uses the latest stable API version. 

Terraform IONOS Cloud Provider **v5 will reach End of Life by September 30, 2023**. After this date, the v5 API will not be accessible. If you require any assistance, please contact our support team.


# IONOS Cloud Terraform Provider

The IonosCloud provider gives the ability to deploy and configure resources using the IonosCloud APIs.

## Migrating from the ProfitBricks provider

Please see the [documentation](docs/index.md#migrating-from-the-profitbricks-provider) on how to migrate from the ProfitBricks provider.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.15.x
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

**NOTE:** In order to use a specific version of this provider, please include the following block at the beginning of your terraform config files [details](https://www.terraform.io/docs/configuration/terraform.html#specifying-a-required-terraform-version):

```terraform
terraform {
  required_providers {
    ionoscloud = {
      source = "ionos-cloud/ionoscloud"
      version = "~> 6.2.0"
    }
  }
}

provider "ionoscloud" {
  username = "ionoscloud_username"
  password = "ionoscloud_password"
}

resource "ionoscloud_datacenter" "main" {
  # ...
}
```

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/ionos-cloud/terraform-provider-ionoscloud`

```sh
$ mkdir -p $GOPATH/src/github.com/ionos-cloud; cd $GOPATH/src/github.com/ionos-cloud
$ git clone git@github.com:ionos-cloud/terraform-provider-ionoscloud
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/ionos-cloud/terraform-provider-ionoscloud
$ make build
```

## Using the provider

See the [IonosCloud Provider documentation](https://registry.terraform.io/providers/ionos-cloud/ionoscloud/latest/docs) to get started using the IonosCloud provider.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is _required_). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-ionoscloud
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
