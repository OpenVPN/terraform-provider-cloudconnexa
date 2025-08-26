package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// dataSourceNetworkConnector returns a Terraform data source resource for CloudConnexa network connectors.
// This resource allows users to read existing network connector configurations.
func dataSourceNetworkConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Use an `cloudconnexa_connector` data source to read an existing CloudConnexa connector.",
		ReadContext: dataSourceNetworkConnectorRead,
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
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the network with which the connector is associated.",
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
			"ipsec_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ipSecConfigSchema(),
			},
		},
	}
}

// dataSourceNetworkConnectorRead handles the read operation for the network connector data source.
// It retrieves the network connector configuration from CloudConnexa using the provided ID
// and updates the Terraform state with the retrieved data.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceNetworkConnectorRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var id = data.Get("id").(string)
	var connector *cloudconnexa.NetworkConnector
	var err error
	if id != "" {
		connector, err = c.NetworkConnectors.GetByID(id)
		if err != nil {
			return append(diags, diag.Errorf("Failed to get network connector with ID: %s, %s", id, err)...)
		}
		if connector == nil {
			return append(diags, diag.Errorf("Network connector with id %s was not found", id)...)
		}
	} else {
		var name = data.Get("name").(string)
		connector, err = c.NetworkConnectors.GetByName(name)
		if err != nil {
			return append(diags, diag.Errorf("Failed to get network connector with name: %s, %s", name, err)...)
		}
		if connector == nil {
			return append(diags, diag.Errorf("Network connector with name %s was not found", name)...)
		}
	}

	setNetworkConnectorData(data, connector)
	if connector.TunnelingProtocol == "OPENVPN" {
		token, err := c.NetworkConnectors.GetToken(connector.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		data.Set("token", token)
		profile, err := c.NetworkConnectors.GetProfile(connector.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		data.Set("profile", profile)
	}
	return diags
}
