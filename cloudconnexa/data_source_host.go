package cloudconnexa

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

func dataSourceHost() *schema.Resource {
	return &schema.Resource{
		Description: "Use an `cloudconnexa_host` data source to read an existing CloudConnexa connector.",
		ReadContext: dataSourceHostRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the host.",
			},
			"internet_access": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of internet access provided.",
			},
			"system_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IPV4 and IPV6 subnets automatically assigned to this host.",
			},
			"connectors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of connectors to be associated with this host.",
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
						"network_item_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the host with which the connector is associated.",
						},
						"network_item_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network object type of the connector. This typically will be set to `HOST`.",
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

func dataSourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	host, err := c.Hosts.GetByName(d.Get("name").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("name", host.Name)
	d.Set("internet_access", host.InternetAccess)
	d.Set("system_subnets", host.SystemSubnets)
	d.Set("connectors", getConnectorsSliceByConnectors(&host.Connectors))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func getConnectorsSliceByConnectors(connectors *[]cloudconnexa.Connector) []interface{} {
	conns := make([]interface{}, len(*connectors))
	for i, c := range *connectors {
		connector := make(map[string]interface{})
		connector["id"] = c.Id
		connector["name"] = c.Name
		connector["network_item_id"] = c.NetworkItemId
		connector["network_item_type"] = c.NetworkItemType
		connector["vpn_region_id"] = c.VpnRegionId
		connector["ip_v4_address"] = c.IPv4Address
		connector["ip_v6_address"] = c.IPv6Address
		connector["description"] = c.Description
		conns[i] = connector
	}
	return conns
}
