package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceVpnRegion returns a Terraform data source resource for CloudConnexa VPN regions.
// This resource allows users to read information about a specific VPN region by its ID.
//
// Returns:
//   - *schema.Resource: A Terraform resource definition for VPN region data source
func dataSourceVpnRegion() *schema.Resource {
	return &schema.Resource{
		Description: "Use a `cloudconnexa_vpn_region` data source to read an CloudConnexa VPN region.",
		ReadContext: dataSourceVpnRegionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the region.",
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
			"region_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the region.",
			},
		},
	}
}

// dataSourceVpnRegionRead handles the read operation for the VPN region data source.
// It retrieves VPN region information from CloudConnexa using the provided ID
// and updates the Terraform state with the retrieved data.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceVpnRegionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := d.Get("id").(string)

	vpnRegion, err := c.VPNRegions.GetByID(id)
	if err != nil {
		return diag.Errorf("Failed to get vpn region with ID: %s, %s", id, err)
	}
	if vpnRegion == nil {
		return append(diags, diag.Errorf("VPN region with id %s was not found", id)...)
	}

	d.SetId(vpnRegion.ID)
	d.Set("continent", vpnRegion.Continent)
	d.Set("country", vpnRegion.Country)
	d.Set("country_iso", vpnRegion.CountryISO)
	d.Set("region_name", vpnRegion.RegionName)

	return diags
}
