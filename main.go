package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/ionoscloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ionoscloud.Provider})
}
