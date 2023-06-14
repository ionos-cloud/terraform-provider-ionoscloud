
[![Gitter](https://img.shields.io/gitter/room/ionos-cloud/sdk-general)](https://gitter.im/ionos-cloud/sdk-general)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=alert_status)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=bugs)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=security_rating)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=terraform-provider&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=terraform-provider)
[![Release](https://img.shields.io/github/v/release/ionos-cloud/terraform-provider-ionoscloud.svg)](https://github.com/ionos-cloud/terraform-provider-ionoscloud/releases/latest)
[![Release Date](https://img.shields.io/github/release-date/ionos-cloud/terraform-provider-ionoscloud.svg)](https://github.com/ionos-cloud/terraform-provider-ionoscloud/releases/latest)
[![Compute-engine test run](https://github.com/ionos-cloud/terraform-provider-ionoscloud/actions/workflows/compute-test-run.yml/badge.svg)](https://github.com/ionos-cloud/terraform-provider-ionoscloud/actions/workflows/compute-test-run.yml)
[![Go](https://img.shields.io/github/go-mod/go-version/ionos-cloud/terraform-provider-ionoscloud.svg)](https://github.com/ionos-cloud/terraform-provider-ionoscloud)

![Alt text](.github/IONOS.CLOUD.BLU.svg?raw=true "Title")

# IONOS Cloud Terraform Provider

The IonosCloud provider gives the ability to deploy and configure resources using the IonosCloud APIs.

## Migrating from the ProfitBricks provider

Please see the [Documentation](docs/index.md#migrating-from-the-profitbricks-provider) on how to migrate from the ProfitBricks provider.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x+
- [Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)

**NOTE:** In order to use a specific version of this provider, please include the following block at the beginning of your terraform config files [details](https://www.terraform.io/docs/configuration/terraform.html#specifying-a-required-terraform-version):

```terraform
provider "ionoscloud" {
  version = "~> 6.3.0"
}
```

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/ionos-cloud/terraform-provider-ionoscloud`

```sh
$ mkdir -p $GOPATH/src/github.com/ionos-cloud; cd $GOPATH/src/github.com/ionos-cloud
$ git clone https://github.com/ionos-cloud/terraform-provider-ionoscloud.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/ionos-cloud/terraform-provider-ionoscloud
$ make build
```

## Using the provider

See the [IonosCloud Provider documentation](https://registry.terraform.io/providers/ionos-cloud/ionoscloud/latest/docs) to get started using the IonosCloud provider.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.18+ is _required_). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-ionoscloud
...
```

## Testing the Provider

### What Are We Testing?

The purpose of our acceptance tests is to **provision** resources containing all the available arguments, followed by **updates** on all arguments that allow this action. Beside the provisioning part, **data-sources** with all possible arguments and **imports** are also tested.

All tests are integrated into [github actions](https://github.com/ionos-cloud/terraform-provider-ionoscloud/actions) that run daily and are also run manually before any release.

### How to Run Tests Locally 

⚠️ **Warning:** Acceptance tests provision resources in the IONOS Cloud, and often may involve extra billing charges on your account.

In order to test the provider, you can simply run:

``` sh
$ make test
```

In order to run the full suite of Acceptance tests, run:

``` sh
$ make testacc TAGS=all
```

#### Test Tags

Tests can also be run for a batch of resources or for a single resource, using tags.

_Example of running server and lan tests:_
``` sh
$ make testacc TAGS=server,lan
```

<details> <summary title="Click to toggle">See more details about <b>test tags</b></summary>

**Build tags** are named as follows:

- `compute` - all **compute engine** tests (datacenter, firewall rule, image, IP block, IP failover, lan, location, nic, private cross connect, server, snapshot, template, volume)
- `nlb` - **network load balancer** and **network load balancer forwarding rule** tests
- `natgateway` - **NAT gateway** and **NAT gateway rule** tests
- `k8s` - **k8s cluster** and **k8s node pool** tests
- `dbaas` - **DBaaS postgres cluster** tests
- `alb` - **Application Load Balancer** tests

``` sh
$ make testacc TAGS=dbaas
```

You can also test one single resource, using one of the tags: `backup`, `datacenter`, `dbaas`, `firewall`, `group`, `image`, `ipblock`, `ipfailover`, `k8s`, `lan`, `location`, `natgateway`,
`nlb`, `nic`, `pcc`, `resource`, `s3key`, `server`, `share`, `snapshot`, `template`, `user`, `volume`

</details>

## Certificate pinning:

You can enable certificate pinning if you want to bypass the normal certificate checking procedure, 
by doing the following: 

Set env variable IONOS_PINNED_CERT=<insert_sha256_public_fingerprint_here>

You can get the sha256 fingerprint most easily from the browser by inspecting the certificate.

## Debugging

In the default mode, the Terraform provider returns only HTTP client errors. These usually consist only of the HTTP status code. There is no clear description of the problem. But if you want to see the API call error messages as well, you need to set the SDK and Terraform provider environment variables.

You can enable logging now using the `IONOS_LOG_LEVEL` env variable. Allowed values: `off`, `debug` and `trace`. Defaults to `off`.

⚠️ **Note:** We recommend you only use `trace` level for debugging purposes. Disable it in your production environments because it can log sensitive data. It logs the full request and response without encryption, even for an HTTPS call.
Verbose request and response logging can also significantly impact your application’s performance.

```bash
$ export IONOS_LOG_LEVEL=debug
```

⚠️ **Note:** `IONOS_DEBUG` is now deprecated and will be removed in a future release.

⚠️ **Note:** We recommend you only use `IONOS_DEBUG` for debugging purposes. Disable it in your production environments because it can log sensitive data. It logs the full request and response without encryption, even for an HTTPS call.
Verbose request and response logging can also significantly impact your application’s performance.

```bash
$ export TF_LOG=debug
$ export IONOS_DEBUG=true
$ terraform apply
```
now you can see the response body incl. api error message:
```json
{
  "httpStatus" : 422,
  "messages" : [ {
    "errorCode" : "200",
    "message" : "[VDC-yy-xxxx] Operation cannot be executed since this Kubernetes Nodepool is already marked for deletion. Current state of the resource is FAILED_DESTROYING."
  }]
}
```