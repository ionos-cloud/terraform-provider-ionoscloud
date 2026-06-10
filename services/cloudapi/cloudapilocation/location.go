package cloudapilocation

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// findLocationById retrieves a single location by its "<region>/<id>" identifier.
func findLocationById(ctx context.Context, client *ionoscloud.APIClient, locationID string) (*ionoscloud.Location, error) {
	parts := strings.SplitN(locationID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid location id %q, expected format <region>/<id>", locationID)
	}
	location, apiResponse, err := client.LocationsApi.LocationsFindByRegionIdAndId(ctx, parts[0], parts[1]).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, err
	}
	return &location, nil
}

// ResolveParentLocation returns the location ids whose images and image aliases are
// usable from the requested location: the location itself, plus the parent location it
// inherits them from (LocationProperties.MetroRegion) when it has one. The result always
// contains at least the requested location id; a fetch failure is logged and degrades to
// that. The parent location is deliberately never fetched — restricted contracts may not
// be able to read it, and image/alias matching works off the image listing instead.
func ResolveParentLocation(ctx context.Context, client *ionoscloud.APIClient, locationID string) []string {
	locationIDs := []string{locationID}

	location, err := findLocationById(ctx, client, locationID)
	if err != nil {
		tflog.Warn(ctx, "could not resolve parent location, using the requested location only",
			map[string]any{"location": locationID, "error": err.Error()})
		return locationIDs
	}
	if location == nil || location.Properties == nil ||
		location.Properties.MetroRegion == nil || *location.Properties.MetroRegion == "" {
		return locationIDs
	}
	// Classic locations carry a self-referential metroRegion (e.g. de/fra -> de/fra):
	// they are their own parent, not children.
	if strings.EqualFold(*location.Properties.MetroRegion, locationID) {
		return locationIDs
	}

	return append(locationIDs, *location.Properties.MetroRegion)
}

// LocationInSet reports whether loc is one of the given location ids, case-insensitively.
func LocationInSet(locations []string, loc string) bool {
	for _, l := range locations {
		if strings.EqualFold(l, loc) {
			return true
		}
	}
	return false
}
