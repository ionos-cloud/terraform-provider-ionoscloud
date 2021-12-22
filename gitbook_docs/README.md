# Terraform provider

**NOTE:**
_Please consider using Terraform IonosCloud Provider v6 releases that are using IONOS Cloud API Version 6, the latest stable API version. New features and improvements will be integrated in IONOS Cloud API Version 6._

## Overview

The Terraform IonosCloud Provider provides you with access to the IONOS Cloud. The provider supports both simple and complex requests. It is designed for devops engineers and developers who are building their infrastructure in the IONOS Cloud . The provider wraps the IONOS Cloud GO SDK. All operations are performed over SSL and authenticated using your IONOS Cloud portal credentials. The provider can be used within an instance running in IONOS Cloud or directly over the Internet from any application that can send an HTTPS request and receive an HTTPS response.

## Getting Started

An IONOS account is required for access to the Ionos Cloud via the ionoscloud terraform provider; credentials from your registration are used to authenticate against the IONOS Cloud API.

### Installation

Terraform is needed to operate the ionoscloud provider:

* [Terraform](https://www.terraform.io/downloads.html) 0.12.x

**NOTE:** In order to use a speciffic version of this provider, please include the following block at the beginning of your terraform config files [details](https://www.terraform.io/docs/configuration/terraform.html#specifying-a-required-terraform-version):

```text
provider "ionoscloud" {
  version = "~> 5.0.4"
}
```

### Authentication

The username, password and optionally the api endpoint can be manually specified when initializing the provider

```text
provider "ionoscloud" {
  username          = var.ionos_username
  password          = var.ionos_password
  endpoint          = var.ionos_api_endpoint
}
```

Environment variables can also be used; the provider uses the following variables:

* IONOS\_USERNAME - to specify the username used to login
* IONOS\_PASSWORD - to specify the password
* IONOS\_API\_URL - to specify the Ionos Cloud API endpoint \(to be used for development/testing purposes only\)

**Warning**: Make sure to follow the Information Security Best Practices when using credentials within your terraform configuration files.

## FAQ

1. How can I open a bug report/feature request? 

Bug reports and feature requests can be opened in the Issues repository: [https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/new/choose](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/new/choose)

1. Can I contribute to the provider?

Pull requests can be open in the [https://github.com/ionos-cloud/terraform-provider-ionoscloud](https://github.com/ionos-cloud/terraform-provider-ionoscloud) repository.

