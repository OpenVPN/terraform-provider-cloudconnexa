package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceNetwork returns a Terraform resource schema for managing networks
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
			"tunneling_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "OPENVPN",
				ValidateFunc: validation.StringInSlice([]string{"OPENVPN", "IPSEC"}, false),
				Description:  "The tunneling protocol used for this network.",
			},
		},
	}
}

// resourceNetworkCreate creates a new network
func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	n := cloudconnexa.Network{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Egress:            d.Get("egress").(bool),
		InternetAccess:    d.Get("internet_access").(string),
		TunnelingProtocol: d.Get("tunneling_protocol").(string),
	}
	network, err := c.Networks.Create(n)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(network.ID)
	return diags
}

// resourceNetworkRead reads the state of a network
func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := d.Id()
	network, err := c.Networks.Get(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get network with ID: %s, %s", id, err)...)
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
	d.Set("tunneling_protocol", network.TunnelingProtocol)
	return diags
}

// resourceNetworkUpdate updates an existing network
func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	_, newName := d.GetChange("name")
	_, newDescription := d.GetChange("description")
	_, newEgress := d.GetChange("egress")
	_, newAccess := d.GetChange("internet_access")
	_, tunnelingProtocol := d.GetChange("tunneling_protocol")
	err := c.Networks.Update(cloudconnexa.Network{
		ID:                d.Id(),
		Name:              newName.(string),
		Description:       newDescription.(string),
		Egress:            newEgress.(bool),
		InternetAccess:    newAccess.(string),
		TunnelingProtocol: tunnelingProtocol.(string),
	})
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return append(diags, resourceNetworkRead(ctx, d, m)...)
}

// resourceNetworkDelete deletes a network
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
