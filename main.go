package main

import (
	"context"
	"flag"
	"log"

	"github.com/SimCubeLtd/terraform-provider-nexus/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Generate docs for website
//go:generate go run github.com/datadrivers/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/datadrivers/nexus",
			&plugin.ServeOpts{
				ProviderFunc: provider.Provider,
			})
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: provider.Provider})
	}
}
