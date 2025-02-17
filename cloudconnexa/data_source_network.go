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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network name.",
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
		},
	}
}

func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var network *cloudconnexa.Network
	var err error
	id := d.Get("id").(string)
	network, err = c.Networks.Get(id)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if network == nil {
		return append(diags, diag.Errorf("Network with id %s was not found", id)...)
	}
	d.SetId(network.ID)
	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("egress", network.Egress)
	d.Set("internet_access", network.InternetAccess)
	d.Set("system_subnets", network.SystemSubnets)
	return diags
}

func getRoutesSlice(networkRoutes *[]cloudconnexa.Route) []interface{} {
	routes := make([]interface{}, len(*networkRoutes))
	for i, r := range *networkRoutes {
		route := make(map[string]interface{})
		route["id"] = r.ID
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
		connector["id"] = c.ID
		connector["name"] = c.Name
		connector["network_id"] = c.NetworkItemID
		connector["vpn_region_id"] = c.VpnRegionID
		connector["ip_v4_address"] = c.IPv4Address
		connector["ip_v6_address"] = c.IPv6Address
		connector["description"] = c.Description
		conns[i] = connector
	}
	return conns
}
