package mariadbv2

import (
	"fmt"
	"sort"
	"strings"

	mariadbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// Client wraps the MariaDB v2 API SDK client.
type Client struct {
	sdkClient *mariadbv3.APIClient
}

// NewClientFromConfig creates a *Client from an existing shared.Configuration.
func NewClientFromConfig(config *shared.Configuration) *Client {
	return &Client{
		sdkClient: mariadbv3.NewAPIClient(config),
	}
}

var LocationToURL = map[string]string{
	"de/fra":   "https://mariadb.de-fra.ionos.com/v2",
	"de/fra/1": "https://mariadb.de-fra.ionos.com/v2",
	"de/fra/2": "https://mariadb.de-fra.ionos.com/v2",
	"de/txl":   "https://mariadb.de-txl.ionos.com/v2",
	"es/vit":   "https://mariadb.es-vit.ionos.com/v2",
	"fr/par":   "https://mariadb.fr-par.ionos.com/v2",
	"gb/bhx":   "https://mariadb.gb-bhx.ionos.com/v2",
	"gb/lhr":   "https://mariadb.gb-lhr.ionos.com/v2",
	"us/ewr":   "https://mariadb.us-ewr.ionos.com/v2",
	"us/las":   "https://mariadb.us-las.ionos.com/v2",
	"us/mci":   "https://mariadb.us-mci.ionos.com/v2",
}

// AvailableLocations returns a sorted list of available MariaDB v2 locations.
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
