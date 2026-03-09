package shared

import (
	"context"
	"fmt"
	"regexp"
)

var uuidRegex = regexp.MustCompile(
	`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
)

// IsUUID returns true if s matches the standard UUID format (any version).
func IsUUID(s string) bool {
	return uuidRegex.MatchString(s)
}

// Resolve resolves a human-readable name to a resource ID.
// Works directly with value-type slices returned by SDK list calls.
//
// If nameOrID is already a valid UUID, it is returned as-is without
// making any API calls. Otherwise, listByName is called to find
// resources matching the given name. Exactly one match is expected;
// zero matches returns an error, and multiple matches returns an
// ambiguity error.
//
// Example usage with the DNS SDK:
//
//	id, err := shared.Resolve(ctx, "example.com", func(ctx context.Context, name string) ([]dns.ZoneRead, error) {
//	    list, _, err := client.ZonesApi.ZonesGet(ctx).FilterZoneName(name).Limit(2).Execute()
//	    if err != nil {
//	        return nil, err
//	    }
//	    return list.GetItems(), nil
//	})
func Resolve[T any, PT interface {
	*T
	Identifiable
}](
	ctx context.Context,
	nameOrID string,
	listByName func(ctx context.Context, name string) ([]T, error),
) (string, error) {
	if IsUUID(nameOrID) {
		return nameOrID, nil
	}

	items, err := listByName(ctx, nameOrID)
	if err != nil {
		return "", fmt.Errorf("resolving %q: %w", nameOrID, err)
	}

	switch len(items) {
	case 0:
		return "", fmt.Errorf("no resource found matching name %q", nameOrID)
	case 1:
		return PT(&items[0]).GetId(), nil
	default:
		return "", fmt.Errorf("ambiguous: %d resources match name %q", len(items), nameOrID)
	}
}
