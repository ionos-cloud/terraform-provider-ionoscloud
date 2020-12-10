# Terraform Provider

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Migrating from the ProfitBricks provider

### Provider Name

The provider name changed from `profitbricks` to `ionoscloud`.
This reflects in the following change in your terraform file:
`provider "profitbricks"` becomes `provider "ionoscloud"`

### Resources and Datasources

The migration affects resource names and datasource names.
Every resource and datasource changed its prefix from `profitbricks_` to `ionoscloud_`.

In order to accommodate that, the terraform state must be updated.

The local state, in json format, can be updated with a simple find and replace procedure.
For example, on Linux, sed can be used:
```
$ sed -i 's/profitbricks_/ionoscloud_/g' ./main.tf
```

On OSX the same command becomes:
```
$ sed -i bak 's/profitbricks_/ionoscloud_/g' ./main.tf
```

If you manage your state using remote backends you need to take the appropriate steps specific to your backend.

### Environment Variables

The following env variables have changed:

| Old Variable Name     | New Variable Name |
|-----------------------|-------------------|
| PROFITBRICKS_USERNAME | IONOS_USERNAME    |
| PROFITBRICKS_PASSWORD | IONOS_PASSWORD    |
| PROFITBRICKS_TOKEN    | IONOS_TOKEN       |
| PROFITBRICKS_API_URL  | IONOS_API_URL     |

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

**NOTE:** In order to use a speciffic version of this provider, please include the following block at the beginning of your terraform config files [details](https://www.terraform.io/docs/configuration/terraform.html#specifying-a-required-terraform-version):

```terraform
provider "ionoscloud" {
  version = "~> 5.0.0"
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

See the [IonosCloud Provider documentation](https://www.terraform.io/docs/providers/ionoscloud/index.html) to get started using the IonosCloud provider.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is _required_). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
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
