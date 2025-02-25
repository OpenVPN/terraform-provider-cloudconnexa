package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func dataSourceNetworkRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	networkId := d.Get("id").(string)
	routes, err := c.Routes.List(networkId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	configRoutes := make([]map[string]interface{}, len(routes))
	for i, r := range routes {
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
	d.SetId(networkId)

	return diags
}
