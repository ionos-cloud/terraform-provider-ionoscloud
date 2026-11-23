# CLAUDE.md - terraform-provider-ionoscloud

## Project Overview

Hybrid Terraform provider: legacy SDK v2 in `ionoscloud/`, modern terraform-plugin-framework in `internal/framework/`, muxed via `terraform-plugin-mux` in `main.go`. Resource/data-source names must be unique across both providers. **Use the Framework for all new resources, data sources, and ephemerals.** The codebase quality is uneven — before implementing something new or imitating an existing pattern, evaluate the reference critically, don't treat it as a template just because it exists.

## Build & Test Commands

```bash
make build          # Build and install provider
make test           # Unit tests
make testacc        # Acceptance tests (requires TF_ACC=1 and IONOS credentials)
make lint           # golangci-lint on changed files
make fmt            # Format code
make fmtcheck       # Check formatting
```

### Running specific acceptance tests

```bash
# By build tag (compute, k8s, dbaas, server, dns, alb, nlb, etc.)
TF_ACC=1 go test ./... -v -failfast -timeout 240m -tags compute

# By test name
TF_ACC=1 go test ./... -v -run TestAccResourceServer -timeout 120m
```

## Linting

Run `make lint` after writing code. Config in `.golangci.yml`.

## Writing Framework resources

When implementing new resources / data-sources / ephemerals or modifying anything under `internal/framework/`, see the `framework-development` skill — it covers file layout, schema design rules, CRUD patterns, async polling, response mapping, write-only fields, import state, data source and ephemeral patterns, the service layer, and acceptance tests.
