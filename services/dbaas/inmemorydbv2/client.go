package inmemorydbv2

import (
	"fmt"
	"sort"
	"strings"

	inmemorydbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// Client wraps the InMemoryDB v2 API SDK client.
type Client struct {
	sdkClient *inmemorydbv3.APIClient
}

// NewClientFromConfig creates a *Client from an existing shared.Configuration.
func NewClientFromConfig(config *shared.Configuration) *Client {
	return &Client{
		sdkClient: inmemorydbv3.NewAPIClient(config),
	}
}

var LocationToURL = map[string]string{
	"de/fra": "https://in-memory-db.de-fra.ionos.com/v2",
	"de/txl": "https://in-memory-db.de-txl.ionos.com/v2",
	"es/vit": "https://in-memory-db.es-vit.ionos.com/v2",
	"gb/bhx": "https://in-memory-db.gb-bhx.ionos.com/v2",
	"gb/lhr": "https://in-memory-db.gb-lhr.ionos.com/v2",
	"us/ewr": "https://in-memory-db.us-ewr.ionos.com/v2",
	"us/las": "https://in-memory-db.us-las.ionos.com/v2",
	"us/mci": "https://in-memory-db.us-mci.ionos.com/v2",
	"fr/par": "https://in-memory-db.fr-par.ionos.com/v2",
}

// AvailableLocations returns a sorted list of available InMemoryDB v2 locations.
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
	quoted := make([]string, len(locations))
	for i, loc := range locations {
		quoted[i] = fmt.Sprintf("`%s`", loc)
	}
	return strings.Join(quoted, ", ")
}
