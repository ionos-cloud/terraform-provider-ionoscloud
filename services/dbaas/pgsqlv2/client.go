package pgsqlv2

import (
	"fmt"
	"sort"
	"strings"

	pgsqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// Client is a wrapper around the PgSQL v2 SDK client.
type Client struct {
	sdkClient *pgsqlv2.APIClient
}

// NewClientFromConfig creates a *Client from an existing shared.Configuration.
func NewClientFromConfig(config *shared.Configuration) *Client {
	return &Client{
		sdkClient: pgsqlv2.NewAPIClient(config),
	}
}

var LocationToURL = map[string]string{
	"de/txl":   "https://postgresql.de-txl.ionos.com",
	"de/fra":   "https://postgresql.de-fra.ionos.com",
	"de/fra/2": "https://postgresql.de-fra.ionos.com",
	"fr/par":   "https://postgresql.fr-par.ionos.com",
	"es/vit":   "https://postgresql.es-vit.ionos.com",
	"gb/lhr":   "https://postgresql.gb-lhr.ionos.com",
	"gb/bhx":   "https://postgresql.gb-bhx.ionos.com",
	"us/las":   "https://postgresql.us-las.ionos.com",
	"us/mci":   "https://postgresql.us-mci.ionos.com",
	"us/ewr":   "https://postgresql.us-ewr.ionos.com",
}

// AvailableLocations returns a sorted list of available PgSQL v2 locations.
func AvailableLocations() []string {
	locations := make([]string, 0, len(LocationToURL))
	for loc := range LocationToURL {
		locations = append(locations, loc)
	}
	sort.Strings(locations)
	return locations
}

// AvailableLocationsString returns a comma-separated string of available PgSQL v2 locations,
// each location enclosed in backticks for display in Terraform schema descriptions.
func AvailableLocationsString() string {
	locations := AvailableLocations()
	quoted := make([]string, len(locations))
	for i, loc := range locations {
		quoted[i] = fmt.Sprintf("`%s`", loc)
	}
	return strings.Join(quoted, ", ")
}
