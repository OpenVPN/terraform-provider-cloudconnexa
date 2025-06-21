package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceHostConnector returns a Terraform resource schema for managing CloudConnexa host connectors.
// It defines the CRUD operations and schema for connector configuration.
func resourceHostConnector() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_connector` to create an CloudConnexa connector.\n\n~> NOTE: This only creates the CloudConnexa connector object. Additional manual steps are required to associate a host in your infrastructure with the connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
		CreateContext: resourceHostConnectorCreate,
		ReadContext:   resourceHostConnectorRead,
		DeleteContext: resourceHostConnectorDelete,
		UpdateContext: resourceHostConnectorUpdate,
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
			"host_id": {
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
		},
	}
}

// resourceHostConnectorUpdate updates an existing CloudConnexa host connector with new configuration.
// It handles updating the connector's name, description, and VPN region.
func resourceHostConnectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	connector := cloudconnexa.HostConnector{
		ID:          d.Id(),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		VpnRegionID: d.Get("vpn_region_id").(string),
	}
	_, err := c.HostConnectors.Update(connector)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

// resourceHostConnectorCreate creates a new CloudConnexa host connector.
// It initializes the connector with the specified configuration and returns a warning about manual setup requirements.
func resourceHostConnectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	networkItemId := d.Get("host_id").(string)
	vpnRegionId := d.Get("vpn_region_id").(string)
	connector := cloudconnexa.HostConnector{
		Name:            name,
		NetworkItemID:   networkItemId,
		NetworkItemType: "HOST",
		VpnRegionID:     vpnRegionId,
		Description:     description,
	}
	conn, err := c.HostConnectors.Create(connector, networkItemId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(conn.ID)
	profile, err := c.HostConnectors.GetProfile(conn.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("profile", profile)
	token, err := c.HostConnectors.GetToken(conn.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("token", token)
	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Connector needs to be set up manually",
		Detail:   "Terraform only creates the CloudConnexa connector object, but additional manual steps are required to associate a host in your infrastructure with this connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
	})
}

// resourceHostConnectorRead retrieves the current state of a CloudConnexa host connector.
// It fetches the connector's configuration, profile, and token information.
func resourceHostConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := d.Id()
	connector, err := c.HostConnectors.GetByID(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get host connector with ID: %s, %s", id, err)...)
	}
	token, err := c.HostConnectors.GetToken(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if connector == nil {
		d.SetId("")
	} else {
		d.SetId(connector.ID)
		d.Set("name", connector.Name)
		d.Set("vpn_region_id", connector.VpnRegionID)
		d.Set("host_id", connector.NetworkItemID)
		d.Set("ip_v4_address", connector.IPv4Address)
		d.Set("ip_v6_address", connector.IPv6Address)
		d.Set("token", token)
		profile, err := c.HostConnectors.GetProfile(connector.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("profile", profile)
	}
	return diags
}

// resourceHostConnectorDelete removes a CloudConnexa host connector.
// It deletes the connector and its associated host configuration.
func resourceHostConnectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	err := c.HostConnectors.Delete(d.Id(), d.Get("host_id").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
