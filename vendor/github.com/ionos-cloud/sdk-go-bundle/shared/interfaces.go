package shared

// Identifiable is satisfied by any SDK model type that exposes a string ID.
// All *Read types across SDK products (dns.ZoneRead, compute.Datacenter,
// vpn.WireguardGatewayRead, etc.) implement this interface via their
// generated GetId() method.
type Identifiable interface {
	GetId() string
}

// HasHref is satisfied by any SDK model type that exposes an href link.
type HasHref interface {
	GetHref() string
}

// Resource combines Identifiable and HasHref — the two universally
// consistent accessor methods across all SDK product model types.
type Resource interface {
	Identifiable
	HasHref
}

// Listable is satisfied by SDK list response types that expose their
// items via GetItems(). Most generated list types (dns.ZoneReadList,
// vpn.IPSecGatewayReadList, compute.Datacenters, etc.) satisfy this.
type Listable[T any] interface {
	GetItems() []T
}

// HasProperties is satisfied by SDK model types that wrap a properties
// sub-object (e.g., dns.ZoneRead exposes GetProperties() dns.Zone).
type HasProperties[P any] interface {
	GetProperties() P
}

// ExtractIDs returns the IDs from a slice of value-type items whose pointer
// type satisfies Identifiable. Works directly with SDK list results:
//
//	ids := shared.ExtractIDs(list.GetItems())
func ExtractIDs[T any, PT interface {
	*T
	Identifiable
}](items []T) []string {
	ids := make([]string, len(items))
	for i := range items {
		ids[i] = PT(&items[i]).GetId()
	}
	return ids
}

// FindByID returns a pointer to the first item matching the given ID, and
// true if found. Works directly with SDK list results:
//
//	zone, ok := shared.FindByID(list.GetItems(), id)
func FindByID[T any, PT interface {
	*T
	Identifiable
}](items []T, id string) (*T, bool) {
	for i := range items {
		if PT(&items[i]).GetId() == id {
			return &items[i], true
		}
	}
	return nil, false
}

