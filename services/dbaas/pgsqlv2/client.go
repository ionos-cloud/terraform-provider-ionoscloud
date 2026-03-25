package pgsqlv2

import (
	pgsqlv2 "github.com/ionos-cloud/pgsqlv2"
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
	"":         "https://postgresql.de-txl.ionos.com",
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
