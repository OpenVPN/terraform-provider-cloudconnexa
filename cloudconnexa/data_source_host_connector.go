package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

func dataSourceHostConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Use an `cloudconnexa_connector` data source to read an existing CloudConnexa connector.",
		ReadContext: dataSourceHostConnectorRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the connector.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the connector.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the connector.",
			},
			"host_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the host with which the connector is associated.",
			},
			"vpn_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region where the connector is deployed.",
			},
			"ip_v4_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPV4 address of the connector.",
			},
			"ip_v6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPV6 address of the connector.",
			},
			"profile": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OpenVPN profile",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Connector token",
			},
		},
	}
}

func dataSourceHostConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var connector *cloudconnexa.HostConnector
	var err error
	var token string
	id := d.Get("id").(string)
	connector, err = c.HostConnectors.GetByID(id)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if connector == nil {
		return append(diags, diag.Errorf("Connector with id %s was not found", id)...)
	}
	token, err = c.HostConnectors.GetToken(connector.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId(connector.ID)
	d.Set("name", connector.Name)
	d.Set("description", connector.Description)
	d.Set("host_id", connector.NetworkItemID)
	d.Set("vpn_region_id", connector.VpnRegionID)
	d.Set("ip_v4_address", connector.IPv4Address)
	d.Set("ip_v6_address", connector.IPv6Address)
	d.Set("token", token)

	profile, err := c.HostConnectors.GetProfile(connector.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("profile", profile)
	return diags
}
