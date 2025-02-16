package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_network` to create an CloudConnexa Network.",
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the network.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "The display description for this resource. Defaults to `Managed by Terraform`.",
			},
			"egress": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Boolean to control whether this network provides an egress or not.",
			},
			"internet_access": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "SPLIT_TUNNEL_ON",
				ValidateFunc: validation.StringInSlice([]string{"SPLIT_TUNNEL_ON", "SPLIT_TUNNEL_OFF", "RESTRICTED_INTERNET"}, false),
				Description:  "The type of internet access provided. Valid values are `SPLIT_TUNNEL_ON`, `SPLIT_TUNNEL_OFF`, or `RESTRICTED_INTERNET`. Defaults to `SPLIT_TUNNEL_ON`.",
			},
			"system_subnets": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IPV4 and IPV6 subnets automatically assigned to this network.",
			},
			"default_route": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The default route of this network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "IP_V4",
							ValidateFunc: validation.StringInSlice([]string{"IP_V4", "IP_V6"}, false),
							Description:  "The type of route. Valid values are `IP_V4`, `IP_V6`, and `DOMAIN`.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "Managed by Terraform.",
							Description: "The default route description.",
						},
						"subnet": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The target value of the default route.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of this resource.",
						},
					},
				},
			},
			"default_connector": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The default connector of this network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of this resource.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "Managed by Terraform.",
							Description: "The default connection description.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the connector automatically created and attached to this network.",
						},
						"vpn_region_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the region where the default connector will be deployed.",
						},
						"network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parent network id.",
						},
						"ip_v4_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPV4 address of the default connector.",
						},
						"ip_v6_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPV6 address of the default connector.",
						},
						"profile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OpenVPN profile of the connector.",
						},
					},
				},
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	configConnector := d.Get("default_connector").([]interface{})[0].(map[string]interface{})
	connectors := []cloudconnexa.NetworkConnector{
		{
			Name:        configConnector["name"].(string),
			VpnRegionID: configConnector["vpn_region_id"].(string),
			Description: configConnector["description"].(string),
		},
	}
	n := cloudconnexa.Network{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Egress:         d.Get("egress").(bool),
		InternetAccess: d.Get("internet_access").(string),
		Connectors:     connectors,
	}
	network, err := c.Networks.Create(n)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(network.ID)
	configRoute := d.Get("default_route").([]interface{})[0].(map[string]interface{})
	defaultRoute, err := c.Routes.Create(network.ID, cloudconnexa.Route{
		Type:        configRoute["type"].(string),
		Description: configRoute["description"].(string),
		Subnet:      configRoute["subnet"].(string),
	})
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	defaultRouteWithIdSlice := make([]map[string]interface{}, 1)
	defaultRouteWithIdSlice[0] = map[string]interface{}{
		"id":          defaultRoute.ID,
		"description": defaultRoute.Description,
		"type":        defaultRoute.Type,
		"subnet":      defaultRoute.Subnet,
	}
	d.Set("default_route", defaultRouteWithIdSlice)
	connectorsList := make([]interface{}, 1)
	connector := make(map[string]interface{})
	connector["id"] = network.Connectors[0].ID
	connector["name"] = network.Connectors[0].Name
	connector["network_id"] = network.Connectors[0].NetworkItemID
	connector["vpn_region_id"] = network.Connectors[0].VpnRegionID
	connector["ip_v4_address"] = network.Connectors[0].IPv4Address
	connector["ip_v6_address"] = network.Connectors[0].IPv6Address
	client := m.(*cloudconnexa.Client)
	profile, err := client.NetworkConnectors.GetProfile(network.Connectors[0].ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	connector["profile"] = profile
	connectorsList[0] = connector
	err = d.Set("default_connector", connectorsList)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "The default connector for this network needs to be set up manually",
		Detail:   "Terraform only creates the CloudConnexa default connector object for this network, but additional manual steps are required to associate a host in your infrastructure with this connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
	})
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	network, err := c.Networks.Get(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if network == nil {
		d.SetId("")
		return diags
	}
	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("egress", network.Egress)
	d.Set("internet_access", network.InternetAccess)
	d.Set("system_subnets", network.SystemSubnets)
	if len(d.Get("default_connector").([]interface{})) > 0 {
		configConnector := d.Get("default_connector").([]interface{})[0].(map[string]interface{})
		connectorName := configConnector["name"].(string)
		networkConnectors, err := c.NetworkConnectors.ListByNetworkID(network.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		retrievedConnector, err := getConnectorSlice(networkConnectors, network.ID, connectorName, m)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		err = d.Set("default_connector", retrievedConnector)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if len(d.Get("default_route").([]interface{})) > 0 {
		configRoute := d.Get("default_route").([]interface{})[0].(map[string]interface{})
		route, err := c.Routes.GetNetworkRoute(d.Id(), configRoute["id"].(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		if route == nil {
			d.Set("default_route", []map[string]interface{}{})
		} else {
			defaultRoute := []map[string]interface{}{
				{
					"id":          configRoute["id"].(string),
					"type":        route.Type,
					"description": route.Description,
				},
			}
			if route.Type == "IP_V4" || route.Type == "IP_V6" {
				defaultRoute[0]["subnet"] = route.Subnet
			}
			d.Set("default_route", defaultRoute)
		}
	}
	return diags
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	if d.HasChange("default_connector") {
		old, new := d.GetChange("default_connector")
		oldSlice := old.([]interface{})
		newSlice := new.([]interface{})
		if len(oldSlice) == 0 && len(newSlice) == 1 {
			// This happens when importing the resource
			newConnector := cloudconnexa.NetworkConnector{
				Name:            newSlice[0].(map[string]interface{})["name"].(string),
				VpnRegionID:     newSlice[0].(map[string]interface{})["vpn_region_id"].(string),
				Description:     newSlice[0].(map[string]interface{})["description"].(string),
				NetworkItemType: "NETWORK",
			}
			_, err := c.NetworkConnectors.Create(newConnector, d.Id())
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
		} else {
			oldMap := oldSlice[0].(map[string]interface{})
			newMap := newSlice[0].(map[string]interface{})
			if oldMap["name"].(string) != newMap["name"].(string) ||
				oldMap["description"].(string) != newMap["description"].(string) ||
				oldMap["vpn_region_id"].(string) != newMap["vpn_region_id"].(string) {

				newConnector := cloudconnexa.NetworkConnector{
					ID:              oldMap["id"].(string),
					Name:            newMap["name"].(string),
					VpnRegionID:     newMap["vpn_region_id"].(string),
					Description:     newMap["description"].(string),
					NetworkItemType: "NETWORK",
				}
				_, err := c.NetworkConnectors.Update(newConnector)
				if err != nil {
					return append(diags, diag.FromErr(err)...)
				}
			}
		}
	}
	if d.HasChange("default_route") {
		old, new := d.GetChange("default_route")
		oldSlice := old.([]interface{})
		newSlice := new.([]interface{})
		if len(oldSlice) == 0 && len(newSlice) == 1 {
			// This happens when importing the resource
			newMap := newSlice[0].(map[string]interface{})
			routeType := newMap["type"]
			routeDesc := newMap["description"]
			routeSubnet := newMap["subnet"]
			route := cloudconnexa.Route{
				Type:        routeType.(string),
				Description: routeDesc.(string),
				Subnet:      routeSubnet.(string),
			}
			defaultRoute, err := c.Routes.Create(d.Id(), route)
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
			defaultRouteWithIdSlice := make([]map[string]interface{}, 1)
			defaultRouteWithIdSlice[0] = map[string]interface{}{
				"id":          defaultRoute.ID,
				"description": defaultRoute.Description,
			}
			err = d.Set("default_route", defaultRouteWithIdSlice)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		} else {
			newMap := newSlice[0].(map[string]interface{})
			routeId := newMap["id"]
			routeType := newMap["type"]
			routeDesc := newMap["description"]
			routeSubnet := newMap["subnet"]
			route := cloudconnexa.Route{
				ID:          routeId.(string),
				Type:        routeType.(string),
				Description: routeDesc.(string),
				Subnet:      routeSubnet.(string),
			}
			err := c.Routes.Update(route)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	}
	if d.HasChanges("name", "description", "internet_access", "egress") {
		_, newName := d.GetChange("name")
		_, newDescription := d.GetChange("description")
		_, newEgress := d.GetChange("egress")
		_, newAccess := d.GetChange("internet_access")
		err := c.Networks.Update(cloudconnexa.Network{
			ID:             d.Id(),
			Name:           newName.(string),
			Description:    newDescription.(string),
			Egress:         newEgress.(bool),
			InternetAccess: newAccess.(string),
		})
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	return append(diags, resourceNetworkRead(ctx, d, m)...)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	networkId := d.Id()
	err := c.Networks.Delete(networkId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
