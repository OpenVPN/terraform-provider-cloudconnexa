package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceNetworkConnector returns a Terraform resource schema for managing network connectors
func resourceNetworkConnector() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_connector` to create an CloudConnexa connector.\n\n~> NOTE: This only creates the CloudConnexa connector object. Additional manual steps are required to associate a host in your infrastructure with the connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
		CreateContext: resourceNetworkConnectorCreate,
		ReadContext:   resourceNetworkConnectorRead,
		DeleteContext: resourceNetworkConnectorDelete,
		UpdateContext: resourceNetworkConnectorUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The connector display name.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "The description for the UI. Defaults to `Managed by Terraform`.",
			},
			"vpn_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the region where the connector will be deployed.",
			},
			"network_id": {
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
				Sensitive:   true,
				Description: "OpenVPN profile of the connector.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Connector token.",
			},
			"ipsec_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable IPsec tunnel for this connector.",
			},
			"ipsec_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current IPsec tunnel status.",
			},
		},
	}
}

// resourceNetworkConnectorUpdate updates an existing network connector
func resourceNetworkConnectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	connector := cloudconnexa.NetworkConnector{
		ID:          d.Id(),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		VpnRegionID: d.Get("vpn_region_id").(string),
	}
	_, err := c.NetworkConnectors.Update(connector)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// Handle IPsec configuration changes
	if d.HasChange("ipsec_enabled") {
		ipsecEnabled := d.Get("ipsec_enabled").(bool)
		if ipsecEnabled {
			_, err := c.NetworkConnectors.StartIPsec(d.Id())
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
			// Set ipsec_enabled to true in state after successful start
			d.Set("ipsec_enabled", true)
		} else {
			_, err := c.NetworkConnectors.StopIPsec(d.Id())
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
			// Set ipsec_enabled to false in state after successful stop
			d.Set("ipsec_enabled", false)
		}

		// Refresh state to update ipsec_status after starting/stopping IPsec
		return resourceNetworkConnectorRead(ctx, d, m)
	}

	return diags
}

// resourceNetworkConnectorCreate creates a new network connector
func resourceNetworkConnectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	networkItemId := d.Get("network_id").(string)
	vpnRegionId := d.Get("vpn_region_id").(string)
	connector := cloudconnexa.NetworkConnector{
		Name:            name,
		NetworkItemID:   networkItemId,
		NetworkItemType: "NETWORK",
		VpnRegionID:     vpnRegionId,
		Description:     description,
	}
	conn, err := c.NetworkConnectors.Create(connector, networkItemId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(conn.ID)
	profile, err := c.NetworkConnectors.GetProfile(conn.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("profile", profile)
	token, err := c.NetworkConnectors.GetToken(conn.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("token", token)

	// Handle IPsec configuration if enabled
	ipsecEnabled := d.Get("ipsec_enabled").(bool)
	if ipsecEnabled {
		_, err := c.NetworkConnectors.StartIPsec(conn.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}

		// Set ipsec_enabled to true in state after successful start
		d.Set("ipsec_enabled", true)

		// Refresh state to update ipsec_status after starting IPsec
		readDiags := resourceNetworkConnectorRead(ctx, d, m)
		diags = append(diags, readDiags...)
	} else {
		// Set ipsec_enabled to false in state when not enabled
		d.Set("ipsec_enabled", false)
	}

	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Connector needs to be set up manually",
		Detail:   "Terraform only creates the CloudConnexa connector object, but additional manual steps are required to associate a host in your infrastructure with this connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
	})
}

// resourceNetworkConnectorRead reads the state of a network connector
func resourceNetworkConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	connector, err := c.NetworkConnectors.GetByID(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	token, err := c.NetworkConnectors.GetToken(connector.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if connector == nil {
		d.SetId("")
	} else {
		d.SetId(connector.ID)
		d.Set("name", connector.Name)
		d.Set("vpn_region_id", connector.VpnRegionID)
		d.Set("network_id", connector.NetworkItemID)
		d.Set("ip_v4_address", connector.IPv4Address)
		d.Set("ip_v6_address", connector.IPv6Address)
		d.Set("token", token)
		profile, err := c.NetworkConnectors.GetProfile(connector.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("profile", profile)

		// Update IPsec status based on current configuration
		// Note: Since the API doesn't provide IPsec status directly in NetworkConnector,
		// we set a descriptive status based on the ipsec_enabled configuration
		ipsecEnabled := d.Get("ipsec_enabled").(bool)
		if ipsecEnabled {
			d.Set("ipsec_status", "enabled")
		} else {
			d.Set("ipsec_status", "disabled")
		}
	}
	return diags
}

// resourceNetworkConnectorDelete deletes a network connector
func resourceNetworkConnectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	err := c.NetworkConnectors.Delete(d.Id(), d.Get("network_id").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
