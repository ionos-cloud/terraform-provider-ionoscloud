---
name: framework-development
description: Implement Terraform Provider resources, data sources and ephemerals using the Plugin Framework. Use when writing, reading, modifying code for a terraform-plugin-framework entity (resource, data-source or ephemeral), developing CRUD operations, schema design, state management, acceptance testing for provider resources.
---

# Framework Development

Guidance for writing or modifying resources, data sources, ephemerals, and tests using `terraform-plugin-framework` (the modern provider half) and the matching service layer in `services/<product>/`.

Always look at existing similar resources for reference, but **question the implementation** — the codebase is not very clean or efficient. Don't copy patterns blindly.

## File & Package Layout

New resources go in `internal/framework/services/<product>/`:
```
internal/framework/services/<product>/
  resource_<name>.go
  resource_<name>_test.go
  data_source_<name>.go
  data_source_<name>_test.go
  ephemeral_<name>.go
  ephemeral_<name>_test.go
  resources.go              # Factory: Resources() []func() resource.Resource
  data_sources.go           # Factory: DataSources() []func() datasource.DataSource
  ephemeral_resources.go    # Factory: EphemeralResources() []func() ephemeral.EphemeralResource
```

## Key Conventions

- Resource type name: `req.ProviderTypeName + "_<product>_<name>"`.
- Client access: cast `req.ProviderData` to `*bundleclient.SdkBundle`, then use the product-specific client.

## Registering New Resources/Data Sources

1. Add the factory function to the product's `resources.go` / `data_sources.go` / `ephemeral_resources.go`.
2. If new product package: import it in `internal/framework/provider/provider.go` and add to `Resources()` / `DataSources()` / `EphemeralResources()` slices.
3. Add the resource name constant to `utils/constant/constants.go`.

## Code Style

- Use `resp.Diagnostics.AddError()`/`AddWarning()` for error reporting (never return raw errors).
- Use the `backoff` package for retry/polling on async API operations.
- Timeouts: use the `terraform-plugin-framework-timeouts` block pattern.
- Plan modifiers: use `stringplanmodifier.RequiresReplace()` etc. for `ForceNew` behavior.
- Imports: stdlib, then third-party, then local — enforced by goimports with local prefix `github.com/ionos-cloud/terraform-provider-ionoscloud`.

## Schema Design Rules

### Defaults belong to the API, not the provider
Never hardcode default values in the schema (e.g., `Default: booldefault.StaticBool(false)`). Use `Optional + Computed` instead and let the API assign defaults.

### Use the narrowest type that matches the SDK
If the underlying SDK uses `int32`, use `schema.Int32Attribute` and `types.Int32` — not `Int64Attribute` with casts. Same principle applies to other types.

### Plan modifiers must be justified by the swagger
- Derive resource behavior from the swagger, don't make assumptions.
- When writing schema attributes, check how attributes are defined in the swagger, see create-only vs updatable, `writeOnly`, required vs optional, and add plan modifiers (`RequiresReplace()`, `UseStateForUnknown()`, etc.) accordingly.
- Do not add `UseStateForUnknown()` on `Required` fields.

### Sort schema attributes alphabetically
Within a schema block, sort attributes by name.

## Model Definitions

- Define a separate Go struct for each level of nesting, with `tfsdk` struct tags matching schema attribute names exactly.
- All fields must use Terraform framework types (`types.String`, `types.Int32`, `types.Bool`, etc.) — never raw Go types.
- Nested objects use pointers to their model struct (e.g., `*instancesModel`) to allow null handling.
- Timeouts use `timeouts.Value` from `terraform-plugin-framework-timeouts`.
- Data source models may reuse resource model structures or define a subset.

## Interface Compliance

Always declare compile-time interface checks at the top of the file:
```go
var (
    _ resource.ResourceWithImportState = (*exampleResource)(nil)
    _ resource.ResourceWithConfigure   = (*exampleResource)(nil)
)
```

## CRUD Implementation

All CRUD operations follow the same pattern: **extract from Terraform request → build API request → call API → poll if async → fetch fresh state → map to model → set state**. Polling is needed only if the swagger contains states such as `AVAILABLE` / `PROVISIONING` that imply a wait-until-available mechanism.

**Create:**
1. Extract plan: `req.Plan.Get(ctx, &plan)`.
2. Get timeout: `plan.Timeouts.Create(ctx, utils.DefaultTimeout)`.
3. Build API request from plan via a `buildXxxCreateProperties()` function.
4. Call the API.
5. Poll for completion using `backoff.Retry` with `backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout))`.
6. Fetch the full resource (the create response may not have all fields).
7. Map API response to model via `mapXxxResponseToModel()`.
8. Set state: `resp.State.Set(ctx, &plan)`.

**Read:**
1. Extract state: `req.State.Get(ctx, &state)`.
2. Call the API to fetch the resource.
3. If 404: `resp.State.RemoveResource(ctx)` and return (do not error).
4. Map API response to model.
5. Set state.

**Update:**
1. Extract both plan and state.
2. Use state for IDs and plan for new values.
3. Build update request via `buildXxxUpdateProperties()`.
4. Call the API, poll, fetch fresh state, map, set state (same as Create steps 4–8).

**Delete:**
1. Extract state.
2. Call the delete API.
3. Poll until the resource returns 404 (use an `IsXxxDeleted()` helper).

## Error Handling

- Always check `resp.Diagnostics.HasError()` after appending diagnostics and return early.
- On Read, handle 404 by removing the resource — never error on a missing resource:
  ```go
  if apiResponse != nil && apiResponse.HttpNotFound() {
      resp.State.RemoveResource(ctx)
      return
  }
  ```
- Use `resp.Diagnostics.AddError(summary, detail)` for all errors.
- In polling helpers: `backoff.Permanent(err)` for non-retryable (e.g., API errors); `fmt.Errorf(...)` for retryable (e.g., not yet ready).

## Timeout & Polling

- Extract timeouts early in each CRUD method: `plan.Timeouts.Create(ctx, utils.DefaultTimeout)`.
- Wrap long-running operations with exponential backoff:
  ```go
  err = backoff.Retry(func() error {
      return client.IsXxxReady(ctx, id)
  }, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)))
  ```
- Polling helpers (`IsXxxReady`, `IsXxxDeleted`) belong in the service layer, not the resource.

## Response-to-State Mapping

- Always create dedicated mapping functions — do not inline mapping logic in CRUD methods.
- Use `types.StringValue()` for required strings, `types.StringPointerValue()` for optional/nullable strings, same principle for other types.

## Preserving WriteOnly / Sensitive Fields

When the API does not return a field (e.g., passwords), preserve the value from the existing model during mapping:
```go
var existingPassword types.String
if model.Credentials != nil {
    existingPassword = model.Credentials.Password
}
model.Credentials = &credentialsModel{
    Username: types.StringValue(props.Credentials.Username),
    Password: existingPassword, // Preserved — API never returns this
}
```

## Model-to-Request Building

- Create dedicated `buildXxxCreateProperties()` / `buildXxxUpdateProperties()` functions returning `(ApiType, diag.Diagnostics)`.
- Required fields: use `.ValueString()`, `.ValueInt32()`, etc. directly depending on the type.
- Optional fields: check `.IsNull()` before converting. Use `.ValueStringPointer()`, `.ValueBoolPointer()`, etc. for API nullable fields.

## Import State

Use colon-separated compound IDs when the API needs more than just the resource ID (e.g., location-scoped resources):
```go
func (r *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    parts := strings.Split(req.ID, ":")
    if len(parts) != 2 {
        resp.Diagnostics.AddError("Unexpected Import Identifier",
            fmt.Sprintf("Expected format: '<location>:<cluster_id>'. Got: %q", req.ID))
        return
    }
    resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("location"), parts[0])...)
    resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
```

## Data Source Patterns

**Single resource (by ID or name):** use `ConfigValidators` with `datasourcevalidator.ExactlyOneOf()` to enforce mutual exclusivity:
```go
func (d *clusterDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
    return []datasource.ConfigValidator{
        datasourcevalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
    }
}
```

## Ephemeral Resource Patterns

See `internal/framework/services/kafka/ephemeral_user_credentials.go` as reference.

**Testing ephemeral resources:**
- Use the `echoprovider` from `terraform-plugin-testing` to verify ephemeral values (since they're not in state, you can't use `TestCheckResourceAttr`).
- Register `echoprovider.NewProviderServer()` alongside the regular provider factories.
- Use `ConfigStateChecks` with `statecheck.ExpectKnownValue` against the echo resource to assert ephemeral values.

## Service Layer

- Wrap the generated SDK client in a thin service struct with typed methods.
- All API methods return `(result, *shared.APIResponse, error)`.
- Log API responses with `apiResponse.LogInfo()`.
- For location-scoped APIs, maintain a `LocationToURL` map and create the client with the correct base URL.
- Place polling helpers (`IsXxxReady`, `IsXxxDeleted`) in the service layer — return `nil` when done, `backoff.Permanent(err)` for fatal, `fmt.Errorf(...)` for retryable states.

## Testing

### Structure
```go
//go:build all || dbaas || psqlv2

package pgsqlv2_test
```
- Build tags: `all || <category> || <product>`.
- Use `acctest.TestAccProtoV6ProviderFactories` for Framework resources.
- Run full lifecycle in a single test: Create → Read assertions → Update → Read assertions → Import → Data source checks → Destroy.

### Helpers
- **Exists check:** Fetch the resource from the API and return an error if it doesn't exist.
- **Destroy check:** Iterate all resources of the type in state. For each, call the API and assert 404. Return error if any still exist.
- **Import state ID:** Build the compound import ID (e.g., `location:id`) from the Terraform state.

### Configurations
- If shared infrastructure (datacenter, LAN, etc.) is needed, define a `var infraConfig` string and reuse it.
- Build test configs by concatenating: `infraConfig + resourceConfig + dataSourceConfig`.
- Use `fmt.Sprintf` with positional verbs (`%[1]s`) for repeated values like location.
- Use `ImportStateVerifyIgnore` for fields the API doesn't return (passwords, timeouts, restore blocks).

## Reference Tables and Examples

### Attribute Types

| Terraform Type | Framework Type | Use Case |
|----------------|----------------|----------|
| `string` | `schema.StringAttribute` | Names, IDs |
| `number` | `schema.Int32Attribute`, `schema.Int64Attribute`, `schema.Float64Attribute` | Counts, sizes |
| `bool` | `schema.BoolAttribute` | Feature flags |
| `list` | `schema.ListAttribute`, `schema.ListNestedAttribute` | Ordered collections |
| `set` | `schema.SetAttribute`, `schema.SetNestedAttribute` | Unordered unique items |
| `map` | `schema.MapAttribute` | Key-value pairs |
| `object` | `schema.SingleNestedAttribute` | Complex nested config |

Pick the narrowest type that matches the SDK (see Schema Design Rules above).

### Plan Modifier Examples

```go
// Force replacement when value changes
stringplanmodifier.RequiresReplace()

// Preserve unknown value during plan
stringplanmodifier.UseStateForUnknown()

// Conditional replacement
stringplanmodifier.RequiresReplaceIf(
    func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
        // Custom logic
    },
    "short description",
    "markdown description",
)
```

### Validator Examples

```go
// String validators
stringvalidator.LengthBetween(1, 255)
stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z0-9-]+$`), "must be lowercase alphanumeric with hyphens")
stringvalidator.OneOf("option1", "option2", "option3")

// Int32/Int64 validators
int32validator.Between(1, 100)
int32validator.AtLeast(1)
int64validator.AtMost(1000)

// List validators
listvalidator.SizeAtLeast(1)
listvalidator.SizeAtMost(10)
```

### Sensitive Attributes

```go
"password": schema.StringAttribute{
    Required:  true,
    Sensitive: true,
    Validators: []validator.String{
        stringvalidator.LengthAtLeast(8),
    },
}
```

For write-only fields the API never returns, also follow the "Preserving WriteOnly / Sensitive Fields" pattern above.

### SDKv2 Resource Pattern (legacy `ionoscloud/` package)

> New work goes through the Framework. This pattern is provided for reference when modifying existing SDKv2 resources under `ionoscloud/`.

```go
func ResourceExample() *schema.Resource {
    return &schema.Resource{
        CreateWithoutTimeout: resourceExampleCreate,
        ReadWithoutTimeout:   resourceExampleRead,
        UpdateWithoutTimeout: resourceExampleUpdate,
        DeleteWithoutTimeout: resourceExampleDelete,

        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },

        Schema: map[string]*schema.Schema{
            "name": {
                Type:         schema.TypeString,
                Required:     true,
                ForceNew:     true,
                ValidateFunc: validation.StringLenBetween(1, 255),
            },
        },
    }
}
```

### Diagnostic Helper Examples

```go
// Error
resp.Diagnostics.AddError(
    "Error creating resource",
    fmt.Sprintf("Could not create resource: %v", err),
)

// Warning
resp.Diagnostics.AddWarning(
    "Resource modified outside Terraform",
    "Resource was modified outside of Terraform, state may be inconsistent",
)

// Attribute-scoped error
resp.Diagnostics.AddAttributeError(
    path.Root("name"),
    "Invalid name",
    "Name must be lowercase alphanumeric",
)
```

## Pre-Submission Checklist

- [ ] `make build` passes
- [ ] `make lint` is clean
- [ ] Acceptance tests pass for the affected build tag
- [ ] All CRUD operations are implemented (resources)
- [ ] Import is implemented and verified by an `ImportState` test step
- [ ] Documentation under `docs/resources/`, `docs/data-sources/`, or `docs/ephemerals/` is added/updated with examples
- [ ] Resource name constant added to `utils/constant/constants.go`
- [ ] Factory registered in the product's `resources.go` / `data_sources.go` / `ephemeral_resources.go`
- [ ] Compile-time interface checks declared at the top of the file
- [ ] Schema attributes sorted alphabetically; plan modifiers justified by the swagger
- [ ] Sensitive / write-only attributes are marked and preserved across reads
- [ ] Error messages use `resp.Diagnostics.AddError`/`AddWarning` with clear summaries

## References

- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [Resource Development](https://developer.hashicorp.com/terraform/plugin/framework/resources)
- [Data Source Development](https://developer.hashicorp.com/terraform/plugin/framework/data-sources)
- [Ephemeral Resources](https://developer.hashicorp.com/terraform/plugin/framework/ephemeral-resources)
- [Acceptance Testing](https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests)
- [terraform-plugin-framework GitHub](https://github.com/hashicorp/terraform-plugin-framework)