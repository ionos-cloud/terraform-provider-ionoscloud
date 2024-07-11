package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/provider"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/ionoscloud"
)

func main() {

	var debugMode bool

	// set this to true to enable cli debugging your provider, by running headless and setting TF_REATTACH_PROVIDERS on the terraform terminal
	//this will enable you to debug when running plans from cli.
	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve/goland")
	flag.Parse()
	// log levels need to be shown correctly in terraform when enabling TF_LOG
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	providers := []func() tfprotov5.ProviderServer{
		providerserver.NewProtocol5(provider.New()), // terraform-plugin-framework provider
		ionoscloud.Provider().GRPCProvider,          // terraform-plugin-sdk provider
	}
	ctx := context.Background()
	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt
	if debugMode {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"registry.terraform.io/ionos-cloud/ionoscloud",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
