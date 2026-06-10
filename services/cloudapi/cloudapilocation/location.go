package cloudapilocation

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// ParentLocationInfo describes a location together with the parent location it inherits
// images and image aliases from (LocationProperties.MetroRegion), if any.
type ParentLocationInfo struct {
	// LocationIDs holds the location ids whose images are usable from the requested
	// location: itself, plus its parent location when set.
	LocationIDs []string
	// ParentID is the parent location id, or "" when the location has none.
	ParentID string
}

// FindLocationById retrieves a single location by its "<region>/<id>" identifier.
func FindLocationById(ctx context.Context, client *ionoscloud.APIClient, locationID string) (*ionoscloud.Location, error) {
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

// ResolveParentLocation expands a location id into the set of locations whose images are
// usable from it, following the parent (metroRegion) pointer when present. The result
// always contains at least the requested location id; a fetch failure is logged and
// degrades to that. The parent location is deliberately never fetched — restricted
// contracts may not be able to read it, and image/alias matching works off the image
// listing instead.
func ResolveParentLocation(ctx context.Context, client *ionoscloud.APIClient, locationID string) ParentLocationInfo {
	info := ParentLocationInfo{LocationIDs: []string{locationID}}

	location, err := FindLocationById(ctx, client, locationID)
	if err != nil {
		LogParentResolveWarn(ctx, locationID, err)
		return info
	}
	if location == nil || location.Properties == nil ||
		location.Properties.MetroRegion == nil || *location.Properties.MetroRegion == "" {
		return info
	}
	// Classic locations carry a self-referential metroRegion (e.g. de/fra -> de/fra):
	// they are their own parent, not children.
	if strings.EqualFold(*location.Properties.MetroRegion, locationID) {
		return info
	}

	info.ParentID = *location.Properties.MetroRegion
	info.LocationIDs = append(info.LocationIDs, info.ParentID)
	return info
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

// LogParentResolveWarn emits a uniform warning when parent-location resolution fails and
// the caller falls back to the requested location only.
func LogParentResolveWarn(ctx context.Context, location string, err error) {
	if err == nil {
		return
	}
	tflog.Warn(ctx, "could not resolve parent location, using the requested location only",
		map[string]any{"location": location, "error": err.Error()})
}
