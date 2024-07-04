package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceConnector() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_connector` to create an CloudConnexa connector.\n\n~> NOTE: This only creates the CloudConnexa connector object. Additional manual steps are required to associate a host in your infrastructure with the connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
		CreateContext: resourceConnectorCreate,
		ReadContext:   resourceConnectorRead,
		DeleteContext: resourceConnectorDelete,
		UpdateContext: resourceConnectorUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The connector display name.",
			},
			"vpn_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the region where the connector will be deployed.",
			},
			"network_item_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HOST", "NETWORK"}, false),
				Description:  "The type of network item of the connector. Supported values are `HOST` and `NETWORK`.",
			},
			"network_item_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the network with which this connector is associated.",
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
				Description: "OpenVPN profile of the connector.",
			},
		},
	}
}

func resourceConnectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	connector := cloudconnexa.Connector{
		Id:              d.Id(),
		Name:            d.Get("name").(string),
		NetworkItemId:   d.Get("network_item_id").(string),
		NetworkItemType: d.Get("network_item_type").(string),
		VpnRegionId:     d.Get("vpn_region_id").(string),
	}
	_, err := c.Connectors.Update(connector)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func resourceConnectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	name := d.Get("name").(string)
	networkItemId := d.Get("network_item_id").(string)
	networkItemType := d.Get("network_item_type").(string)
	vpnRegionId := d.Get("vpn_region_id").(string)
	connector := cloudconnexa.Connector{
		Name:            name,
		NetworkItemId:   networkItemId,
		NetworkItemType: networkItemType,
		VpnRegionId:     vpnRegionId,
	}
	conn, err := c.Connectors.Create(connector, networkItemId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(conn.Id)
	profile, err := c.Connectors.GetProfile(conn.Id)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("profile", profile)
	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Connector needs to be set up manually",
		Detail:   "Terraform only creates the CloudConnexa connector object, but additional manual steps are required to associate a host in your infrastructure with this connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
	})
}

func resourceConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	connector, err := c.Connectors.GetByID(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if connector == nil {
		d.SetId("")
	} else {
		d.SetId(connector.Id)
		d.Set("name", connector.Name)
		d.Set("vpn_region_id", connector.VpnRegionId)
		d.Set("network_item_type", connector.NetworkItemType)
		d.Set("network_item_id", connector.NetworkItemId)
		d.Set("ip_v4_address", connector.IPv4Address)
		d.Set("ip_v6_address", connector.IPv6Address)
		profile, err := c.Connectors.GetProfile(connector.Id)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("profile", profile)
	}
	return diags
}

func resourceConnectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	err := c.Connectors.Delete(d.Id(), d.Get("network_item_id").(string), d.Get("network_item_type").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func getConnectorSlice(connectors []cloudconnexa.Connector, networkItemId string, connectorName string, m interface{}) ([]interface{}, error) {
	if len(connectors) == 0 {
		return nil, nil
	}
	connectorsList := make([]interface{}, 1)
	for _, c := range connectors {
		if c.NetworkItemId == networkItemId && c.Name == connectorName {
			connector := make(map[string]interface{})
			connector["id"] = c.Id
			connector["name"] = c.Name
			connector["network_item_id"] = c.NetworkItemId
			connector["network_item_type"] = c.NetworkItemType
			connector["vpn_region_id"] = c.VpnRegionId
			connector["ip_v4_address"] = c.IPv4Address
			connector["ip_v6_address"] = c.IPv6Address
			client := m.(*cloudconnexa.Client)
			profile, err := client.Connectors.GetProfile(c.Id)
			if err != nil {
				return nil, err
			}
			connector["profile"] = profile
			connectorsList[0] = connector
			break
		}
	}
	return connectorsList, nil
}
