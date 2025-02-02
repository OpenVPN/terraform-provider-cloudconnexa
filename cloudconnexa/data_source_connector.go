package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

func dataSourceConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Use an `cloudconnexa_connector` data source to read an existing CloudConnexa connector.",
		ReadContext: dataSourceConnectorRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The ID of the connector.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The name of the connector.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the connector.",
			},
			"network_item_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the network or host with which the connector is associated.",
			},
			"network_item_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network object type of the connector. This typically will be set to either `NETWORK` or `HOST`.",
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

func dataSourceConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var connector *cloudconnexa.Connector
	var err error
	var token string
	connectorName := d.Get("name").(string)
	connectorId := d.Get("id").(string)
	networkItemType := d.Get("network_item_type").(string)
	if connectorId != "" {
		connector, err = c.Connectors.GetByID(connectorId, networkItemType)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		if connector == nil {
			return append(diags, diag.Errorf("Connector with id %s was not found", connectorId)...)
		}
		token, err = c.Connectors.GetToken(connector.Id, networkItemType)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	} else if connectorName != "" {
		connectorsAll, err := c.Connectors.List(networkItemType)
		var connectorCount int
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}

		for _, con := range connectorsAll {
			if con.Name == connectorName {
				connectorCount++
			}
		}

		if connectorCount == 0 {
			return append(diags, diag.Errorf("Connector with name %s was not found", connectorName)...)
		} else if connectorCount > 1 {
			return append(diags, diag.Errorf("More than 1 connector with name %s was found. Please use id instead", connectorName)...)
		} else {
			connector, err = c.Connectors.GetByName(connectorName, networkItemType)
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
		}

		token, err = c.Connectors.GetToken(connector.Id, networkItemType)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}

	} else {
		return append(diags, diag.Errorf("Connector name or id is missing")...)
	}

	d.SetId(connector.Id)
	d.Set("name", connector.Name)
	d.Set("description", connector.Description)
	d.Set("network_item_id", connector.NetworkItemId)
	d.Set("network_item_type", connector.NetworkItemType)
	d.Set("vpn_region_id", connector.VpnRegionId)
	d.Set("ip_v4_address", connector.IPv4Address)
	d.Set("ip_v6_address", connector.IPv6Address)
	d.Set("token", token)

	profile, err := c.Connectors.GetProfile(connector.Id, networkItemType)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("profile", profile)
	return diags
}
