package main

import (
	"github.com/OpenVPN/terraform-provider-cloudconnexa/cloudconnexa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// main is the entry point for the Terraform provider plugin.
// It initializes and serves the CloudConnexa provider using the Terraform plugin SDK.
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return cloudconnexa.Provider()
		},
	})
}
