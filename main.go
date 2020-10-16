package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/utilitywarehouse/terraform-provider-megaport/megaport"
)

var (
	buildCommit  = "unknown"
	buildVersion = "dev"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: megaport.Provider})
}
