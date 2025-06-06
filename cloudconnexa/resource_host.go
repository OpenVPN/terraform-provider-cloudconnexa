package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceHost returns a Terraform resource schema for managing CloudConnexa hosts
func resourceHost() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_host` to create an CloudConnexa host.",
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		UpdateContext: resourceHostUpdate,
		DeleteContext: resourceHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the host.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "The description for the UI. Defaults to `Managed by Terraform`.",
			},
			"domain": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 253),
				Description:  "The domain of the host.",
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
				Description: "The IPV4 and IPV6 subnets automatically assigned to this host.",
			},
		},
	}
}

// resourceHostCreate creates a new CloudConnexa host
func resourceHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	h := cloudconnexa.Host{
		Name:           d.Get("name").(string),
		Domain:         d.Get("domain").(string),
		Description:    d.Get("description").(string),
		InternetAccess: d.Get("internet_access").(string),
	}
	host, err := c.Hosts.Create(h)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(host.ID)

	return diags
}

// resourceHostRead retrieves information about an existing CloudConnexa host
func resourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	host, err := c.Hosts.Get(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if host == nil {
		d.SetId("")
		return diags
	}
	d.Set("name", host.Name)
	d.Set("description", host.Description)
	d.Set("domain", host.Domain)
	d.Set("internet_access", host.InternetAccess)
	d.Set("system_subnets", host.SystemSubnets)

	return diags
}

// resourceHostUpdate updates an existing CloudConnexa host
func resourceHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	_, newName := d.GetChange("name")
	_, newDescription := d.GetChange("description")
	_, newDomain := d.GetChange("domain")
	_, newAccess := d.GetChange("internet_access")
	err := c.Hosts.Update(cloudconnexa.Host{
		ID:             d.Id(),
		Name:           newName.(string),
		Description:    newDescription.(string),
		Domain:         newDomain.(string),
		InternetAccess: newAccess.(string),
	})
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return append(diags, resourceHostRead(ctx, d, m)...)
}

// resourceHostDelete removes an existing CloudConnexa host
func resourceHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	hostId := d.Id()
	err := c.Hosts.Delete(hostId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
