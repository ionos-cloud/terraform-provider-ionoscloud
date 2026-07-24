# s3-key-secret-exposure Specification

## Purpose
TBD - created by archiving change find-1-vulnerability-on-the-terraform-prov-ucolmlg. Update Purpose after archive.
## Requirements
### Requirement: secret_key field marked sensitive
The `secret_key` attribute of `ionoscloud_s3_key` SHALL be declared with `Sensitive: true` in its schema definition so that Terraform masks the value in CLI output, plan diffs, and structured logs.

#### Scenario: Plan output masks secret_key
- **WHEN** a user runs `terraform plan` on a configuration that creates or updates an `ionoscloud_s3_key` resource
- **THEN** the `secret_key` value SHALL appear as `(sensitive value)` in the plan output, never as the raw string

#### Scenario: State file marks secret_key sensitive
- **WHEN** Terraform writes resource state after a successful apply
- **THEN** the state entry for `secret_key` SHALL carry the Terraform sensitive marker, preventing accidental exposure via `terraform output` without explicit opt-in

#### Scenario: terraform output blocks raw value
- **WHEN** a user runs `terraform output secret_key` without flags
- **THEN** Terraform SHALL print `(sensitive value)` and exit 0, not the raw secret

### Requirement: Error logs do not expose secret_key
The provider SHALL NOT log or format the full S3 key API response struct using `%+v` or equivalent reflective formatting in any error, warning, or info log path.

#### Scenario: Read error does not leak secret
- **WHEN** `resourceS3KeyRead` encounters an API error fetching a key
- **THEN** the resulting error message SHALL contain the error text and the key ID but SHALL NOT contain the `SecretKey` string value

#### Scenario: Other CRUD error paths do not leak secret
- **WHEN** any CRUD operation on `ionoscloud_s3_key` returns an error
- **THEN** log messages and returned `diag.Diagnostics` detail strings SHALL NOT include the raw `SecretKey` value

