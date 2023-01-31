# ionoscloud_mongo_template

The **DbaaS Mongo Template data source** can be used to search for and return an existing DbaaS MongoDB Template.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_mongo_template" "example" {
  id = <template_id>
}
```
### By name

```hcl
data "ionoscloud_mongo_template" "example" {
  name = <name>
}
```

### By name, using partial_match

```hcl
data "ionoscloud_mongo_template" "example" {
  name = <name>
  partial_match = true
}
```

* `name` - (Optional) The name of the template you want to search for.
* `id` - (Optional) ID of the template you want to search for.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

Either `name` or `id` must be provided. If none or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The ID of the template.
* `name` - The name of the template.
* `edition` - The edition of the template (e.g. enterprise).
* `cores` - The number of CPU cores.
* `ram` - The amount of memory in GB.
* `storage_size` - The amount of storage size in GB.