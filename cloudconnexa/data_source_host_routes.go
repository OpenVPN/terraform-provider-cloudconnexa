package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceHostRoutes returns a Terraform data source resource for CloudConnexa host routes.
// This resource allows users to read information about routes for a specific host.
func dataSourceHostRoutes() *schema.Resource {
	return &schema.Resource{
		Description: "Use `cloudconnexa_host_routes` data source to retrieve routes for a specific host.",
		ReadContext: dataSourceHostRoutesRead,
		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the host to get routes for.",
			},
			"routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of routes for the host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The route ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The route type.",
						},
						"subnet": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The route subnet.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The route domain (if applicable).",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The route description.",
						},
					},
				},
			},
		},
	}
}

// dataSourceHostRoutesRead handles the read operation for the host routes data source.
func dataSourceHostRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	hostID := d.Get("host_id").(string)
	routes, err := c.HostRoutes.List(hostID)
	if err != nil {
		return diag.Errorf("Failed to get routes for host %s: %s", hostID, err)
	}

	d.SetId(hostID)
	d.Set("routes", flattenHostRoutes(routes))

	return diags
}

// flattenHostRoutes converts a slice of CloudConnexa host routes into a slice of interface{}
func flattenHostRoutes(routes []cloudconnexa.HostRoute) []interface{} {
	result := make([]interface{}, len(routes))
	for i, r := range routes {
		route := map[string]interface{}{
			"id":          r.ID,
			"type":        r.Type,
			"subnet":      r.Subnet,
			"domain":      r.Domain,
			"description": r.Description,
		}
		result[i] = route
	}
	return result
}
