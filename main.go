package main

import (
	"context"
	"flag"
	"log"

	"github.com/OpenVPN/terraform-provider-cloudconnexa/cloudconnexa"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()

	// The provider is served by two implementations behind one mux: the SDK v2 provider owns all
	// existing resources and data sources, the framework provider owns the ephemeral resources.
	// Both speak protocol 5, and both must expose an identical provider block schema.
	providers := []func() tfprotov5.ProviderServer{
		cloudconnexa.Provider().GRPCProvider,
		providerserver.NewProtocol5(cloudconnexa.NewFrameworkProvider()),
	}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt
	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}
	if err := tf5server.Serve("registry.terraform.io/OpenVPN/cloudconnexa", muxServer.ProviderServer, serveOpts...); err != nil {
		log.Fatal(err)
	}
}
