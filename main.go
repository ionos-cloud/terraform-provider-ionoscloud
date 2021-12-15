package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/ionoscloud"
	"log"
)

func main() {
	var debugMode bool

	//set this to true to enable cli debugging your provider, by running headless and setting TF_REATTACH_PROVIDERS on the terraform terminal
	//this will enable you to debug when running plans from cli.
	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve/goland")
	flag.Parse()

	if debugMode {
		err := plugin.Debug(context.Background(), "ionoscloud",
			&plugin.ServeOpts{
				ProviderFunc: ionoscloud.Provider,
			})
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: ionoscloud.Provider,
		})
	}
}
