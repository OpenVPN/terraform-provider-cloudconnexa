package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceVpnRegions returns a Terraform data source resource for CloudConnexa VPN regions.
// This resource allows users to read information about available VPN regions.
//
// Returns:
//   - *schema.Resource: A Terraform resource definition for VPN regions data source
func dataSourceVpnRegions() *schema.Resource {
	return &schema.Resource{
		Description: "Use a `cloudconnexa_vpn_regions` data source to read all CloudConnexa VPN regions.",
		ReadContext: dataSourceVpnRegionsRead,
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of regions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the region.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the region.",
						},
						"continent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The continent of the region.",
						},
						"country": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The country of the region.",
						},
						"country_iso": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ISO code of the country of the region.",
						},
					},
				},
			},
		},
	}
}

// dataSourceVpnRegionsRead handles the read operation for the VPN regions data source.
// It retrieves VPN regions information from CloudConnexa
// and updates the Terraform state with the retrieved data.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceVpnRegionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	regions, err := c.VPNRegions.List()
	if err != nil {
		return diag.Errorf("Failed to get vpn regions: %s", err)
	}
	configRegions := make([]map[string]interface{}, len(regions))
	for i, r := range regions {
		region := make(map[string]interface{})
		region["id"] = r.ID
		region["name"] = r.RegionName
		region["country"] = r.Country
		region["country_iso"] = r.CountryISO
		region["continent"] = r.Continent
		configRegions[i] = region
	}
	d.SetId("regions")
	d.Set("regions", configRegions)
	return diags
}
