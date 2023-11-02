# Introduction

## IONOS Cloud Terraform Provider

The IonosCloud provider gives the ability to deploy and configure resources using the IonosCloud APIs.

### Requirements

* [Terraform](https://www.terraform.io/downloads.html) 0.12.x

**NOTE:** In order to use a specific version of this provider, please include the following block at the beginning of your terraform config files [details](https://www.terraform.io/docs/configuration/terraform.html#specifying-a-required-terraform-version):

```
provider "ionoscloud" {
  version = ">= 6.4.10"
}
```

### Using the provider

See the [IonosCloud Provider documentation](https://registry.terraform.io/providers/ionos-cloud/ionoscloud/latest/docs) to get started using the IonosCloud provider.

### Environment Variables

| Environment Variable    | Description                                                                                                                                                                |
|-------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `IONOS_USERNAME`        | Specify the username used to login, to authenticate against the IONOS Cloud API, unless a token is provided                                                                |
| `IONOS_PASSWORD`        | Specify the password used to login, to authenticate against the IONOS Cloud API, unless a token is provided                                                                |
| `IONOS_TOKEN`           | Specify the token used to login, if a token is being used instead of username and password                                                                                 |
| `IONOS_API_URL`         | Specify th`e API URL. It will overwrite the API endpoint default value `api.ionos.com`.  It is not necessary to override this value unless you have special routing config |
| `IONOS_LOG_LEVEL`       | Specify the Log Level used to log messages. Possible values: Off, Debug, Trace                                                                                             |
| `IONOS_PINNED_CERT`     | Specify the SHA-256 public fingerprint here, enables certificate pinning                                                                                                   |
| `IONOS_CONTRACT_NUMBER` | Specify the contract number on which you wish to provision. Only valid for reseller accounts, for other types of accounts the header will be ignored                       |

### Certificate pinning:

You can enable certificate pinning if you want to bypass the normal certificate checking procedure,
by doing the following:

Set env variable IONOS_PINNED_CERT=`<insert_sha256_public_fingerprint_here>`

You can get the sha256 fingerprint most easily form the browser by inspecting the certificate.

### Debugging

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