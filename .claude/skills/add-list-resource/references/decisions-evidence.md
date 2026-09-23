# Evidence for step 0.5 — the four decisions

`scripts/probe.sh` probes gather this evidence and print verdicts; this is what they mean.

## 1. Pagination

Two questions. `limit-cloud`/`limit-bundle <Op>` answers the first, *does the client send a
default?* Only the request type answers the second, *does the endpoint page at all?*: `Limit` and
`Offset` on it with no default means the server's page size governs, so say that rather than invent
a number — and where that operation's own doc comment names the default ("Default limit is the first
100 items"), state it as the API's figure, with nothing attached to disown it. Neither: it does not
page, so drop the truncation warning. A `datasource-limit <resource>` hit: a bare `.Depth(1)` fetch would cap the
list resource below its own data source, so match it and name that number in the docs.

## 2. Identity attributes

`splitImportID` (`ionoscloud/utils.go`) shows the tuple: `location` where the resource has
one, then one string attribute per parsed part, its own being `id`.

## 3. Filter fields

`name` plus the regional key under the name *the resource's own schema* uses (`location` vs
`region`), checked against that schema: `FilterAttribute` allow-lists a free-form string, so a
`field_name` matching no attribute validates fine and then matches nothing.

- **One field is not enough**: the AND subtest and the two-filter docs example need two that match
  different items. Take the second from the remaining **string** attributes — `MatchesFilters`
  compares strings, so bools and lists cannot be filtered. Nothing narrows → drop both, and say why.
- **Drop a field that cannot narrow**: a one-value validator always matches
  (`ionoscloud_target_group`'s `protocol`: only `"HTTP"`).
- `MatchesFilters` is case-sensitive (`internal/framework/identity/filter.go`), SDKv2
  `StringInSlice(…, true)` is not: a bad `field_name` errors at plan time, a mis-cased value is
  silent. Say so in the docs.

## 4. Which client, regional or global

Global follows the collection's shape, not the constructor: `ionoscloud/resource_datacenter.go`
uses a per-location client in CRUD, yet its list resource makes one
`NewCloudAPIClientWithFailover(ctx)` call — `location` only picks an endpoint override, not a
partition. Never that constructor for a bundle product: its `*ionoscloud.APIClient` cannot reach
that API. Settle partitioning from THIS resource; a product-wide grep corroborates, never decides.

1. Client built per location (`.New<X>Client(ctx, location)`)? A direct yes.
2. Else, does its own schema carry `location`/`region`? If not, one call is the whole collection.
3. Only if partitioned, `locations <product>`. Enumerable → fan out over it
   (`services/vpn`: `AvailableLocations`). **Partitioned but nothing enumerable → stop and ask**: one
   call returns one location's worth of objects and the listing is silently incomplete, with no
   error and nothing in the output to show it. Enumerating the map's keys means exporting a location
   list from that service package — a provider change this skill does not get to make on its own.
