package mongov2

import (
	"fmt"
	"sort"
	"strings"

	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// Client wraps the MongoDB v2 API SDK client.
type Client struct {
	sdkClient *mongov3.APIClient
}

// NewClientFromConfig creates a *Client from an existing shared.Configuration.
func NewClientFromConfig(config *shared.Configuration) *Client {
	return &Client{
		sdkClient: mongov3.NewAPIClient(config),
	}
}

// LocationToURL maps a location to the MongoDB v2 regional API endpoint.
// The endpoints are the servers advertised by the MongoDB v2 API spec.
var LocationToURL = map[string]string{
	"gb/bhx": "https://mongodb.gb-bhx.ionos.com/v2",
	"us/ewr": "https://mongodb.us-ewr.ionos.com/v2",
}

// AvailableLocations returns a sorted list of available MongoDB v2 locations.
func AvailableLocations() []string {
	locations := make([]string, 0, len(LocationToURL))
	for loc := range LocationToURL {
		locations = append(locations, loc)
	}
	sort.Strings(locations)
	return locations
}

// AvailableLocationsString returns a comma-separated list of available locations,
// each enclosed in backticks for use in schema descriptions.
func AvailableLocationsString() string {
	locations := AvailableLocations()
	quoted := make([]string, 0, len(locations))
	for _, loc := range locations {
		quoted = append(quoted, fmt.Sprintf("`%s`", loc))
	}
	return strings.Join(quoted, ", ")
}
