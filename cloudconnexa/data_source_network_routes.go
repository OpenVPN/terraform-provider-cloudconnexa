package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceNetworkRoutes returns a Terraform data source resource for CloudConnexa network routes.
// This resource allows users to read all routes associated with a specific CloudConnexa network.
//
// Returns:
//   - *schema.Resource: A Terraform resource definition for network routes
func dataSourceNetworkRoutes() *schema.Resource {
	return &schema.Resource{
		Description: "Use an `cloudconnexa_network_routes` data source to read all the routes associated with an CloudConnexa network.",
		ReadContext: dataSourceNetworkRoutesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the CloudConnexa network of the routes to be discovered.",
			},
			"routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of routes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the route.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of route. Valid values are `IP_V4`, `IP_V6`.",
						},
						"subnet": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subnet of the route.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A description of the route.",
						},
					},
				},
			},
		},
	}
}

// dataSourceNetworkRoutesRead handles the read operation for the network routes data source.
// It retrieves all routes associated with a specific network ID from CloudConnexa
// and updates the Terraform state with the retrieved data.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceNetworkRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	id := d.Get("id").(string)
	if id == "" {
		return append(diags, diag.Errorf("ID cannot be empty")...)
	}
	network, err := c.Networks.Get(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get network with ID: %s, %s", id, err)...)
	}

	configRoutes := make([]map[string]interface{}, len(network.Routes))
	for i, r := range network.Routes {
		route := make(map[string]interface{})
		route["id"] = r.ID
		route["type"] = r.Type
		route["subnet"] = r.Subnet
		route["description"] = r.Description
		configRoutes[i] = route
	}

	if err := d.Set("routes", configRoutes); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(id)

	return diags
}
