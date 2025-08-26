package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// dataSourceHostConnector returns a Terraform data source resource for CloudConnexa host connectors.
// This resource allows users to read existing host connector configurations.
func dataSourceHostConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Use an `cloudconnexa_connector` data source to read an existing CloudConnexa connector.",
		ReadContext: dataSourceHostConnectorRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the connector.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The name of the connector.",
				ExactlyOneOf: []string{"id", "name"},
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
				Sensitive:   true,
				Description: "OpenVPN profile",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Connector token",
			},
		},
	}
}

// dataSourceHostConnectorRead handles the read operation for the host connector data source.
// It retrieves the host connector configuration from CloudConnexa using the provided ID
// and updates the Terraform state with the retrieved data.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceHostConnectorRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var id = data.Get("id").(string)
	var connector *cloudconnexa.HostConnector
	var token string
	var err error
	if id != "" {
		connector, err = c.HostConnectors.GetByID(id)
		if err != nil {
			return append(diags, diag.Errorf("Failed to get host connector with ID: %s, %s", id, err)...)
		}
		if connector == nil {
			return append(diags, diag.Errorf("Host connector with id %s was not found", id)...)
		}
	} else {
		var name = data.Get("name").(string)
		connector, err = c.HostConnectors.GetByName(name)
		if err != nil {
			return append(diags, diag.Errorf("Failed to get host connector with name: %s, %s", name, err)...)
		}
		if connector == nil {
			return append(diags, diag.Errorf("Host connector with name %s was not found", name)...)
		}
	}
	token, err = c.HostConnectors.GetToken(connector.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	data.SetId(connector.ID)
	data.Set("name", connector.Name)
	data.Set("description", connector.Description)
	data.Set("host_id", connector.NetworkItemID)
	data.Set("vpn_region_id", connector.VpnRegionID)
	data.Set("ip_v4_address", connector.IPv4Address)
	data.Set("ip_v6_address", connector.IPv6Address)
	data.Set("token", token)

	profile, err := c.HostConnectors.GetProfile(connector.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	data.Set("profile", profile)
	return diags
}
