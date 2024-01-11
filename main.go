package main

import (
	"github.com/OpenVPN/terraform-provider-openvpn-cloud/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return cloudconnexa.Provider()
		},
	})
}
