package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/jameswoolfenden/terraform-provider-bridgecrew/bridgecrew"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	opts := &plugin.ServeOpts{
		ProviderFunc: bridgecrew.Provider,
	}

	if debugMode {
		err := plugin.Debug(context.Background(), "jameswoolfenden/dev/bridgecrew", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
