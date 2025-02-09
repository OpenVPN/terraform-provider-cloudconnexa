package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Use a `cloudconnexa_network` data source to read an CloudConnexa network.",
		ReadContext: dataSourceNetworkRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The network ID.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The network name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the network.",
			},
			"egress": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Boolean to indicate whether this network provides an egress or not.",
			},
			"internet_access": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of internet access provided. Valid values are `SPLIT_TUNNEL_ON`, `SPLIT_TUNNEL_OFF`, or `RESTRICTED_INTERNET`. Defaults to `SPLIT_TUNNEL_ON`.",
			},
			"system_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IPV4 and IPV6 subnets automatically assigned to this network.",
			},
			"routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The routes associated with this network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The route id.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of route. Valid values are `IP_V4`, `IP_V6`, and `DOMAIN`.",
						},
						"subnet": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of the route, either an IPV4 address, an IPV6 address, or a DNS hostname.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the route.",
						},
					},
				},
			},
			"connectors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of connectors associated with this network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connector id.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connector name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default connection description.",
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
					},
				},
			},
		},
	}
}

func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var network *cloudconnexa.Network
	var err error
	networkId := d.Get("id").(string)
	networkName := d.Get("name").(string)
	if networkId != "" {
		network, err = c.Networks.Get(networkId)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		if network == nil {
			return append(diags, diag.Errorf("Network with id %s was not found", networkId)...)
		}
	} else if networkName != "" {
		network, err = c.Networks.GetByName(networkName)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		if network == nil {
			return append(diags, diag.Errorf("Network with name %s was not found", networkName)...)
		}
	} else {
		return append(diags, diag.Errorf("Network name or id is missing")...)
	}
	d.SetId(network.Id)
	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("egress", network.Egress)
	d.Set("internet_access", network.InternetAccess)
	d.Set("system_subnets", network.SystemSubnets)
	d.Set("routes", getRoutesSlice(&network.Routes))
	d.Set("connectors", getConnectorsSliceByNetworkConnectors(&network.Connectors))
	return diags
}

func getRoutesSlice(networkRoutes *[]cloudconnexa.Route) []interface{} {
	routes := make([]interface{}, len(*networkRoutes))
	for i, r := range *networkRoutes {
		route := make(map[string]interface{})
		route["id"] = r.Id
		route["subnet"] = r.Subnet
		route["type"] = r.Type
		route["description"] = r.Description
		routes[i] = route
	}
	return routes
}

func getConnectorsSliceByNetworkConnectors(connectors *[]cloudconnexa.NetworkConnector) []interface{} {
	conns := make([]interface{}, len(*connectors))
	for i, c := range *connectors {
		connector := make(map[string]interface{})
		connector["id"] = c.Id
		connector["name"] = c.Name
		connector["network_id"] = c.NetworkItemId
		connector["vpn_region_id"] = c.VpnRegionId
		connector["ip_v4_address"] = c.IPv4Address
		connector["ip_v6_address"] = c.IPv6Address
		connector["description"] = c.Description
		conns[i] = connector
	}
	return conns
}
