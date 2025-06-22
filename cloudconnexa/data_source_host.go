package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// dataSourceHost returns a Terraform data source resource for CloudConnexa hosts.
// This resource allows users to read information about existing CloudConnexa connectors
// by providing a host ID.
func dataSourceHost() *schema.Resource {
	return &schema.Resource{
		Description: "Use an `cloudconnexa_host` data source to read an existing CloudConnexa connector.",
		ReadContext: dataSourceHostRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The host ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the host.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the host.",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host domain.",
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
		},
	}
}

// dataSourceHostRead implements the read operation for the CloudConnexa host data source.
// It retrieves host information using the provided ID and updates the Terraform state
// with the host's attributes.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The meta interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	host, err := c.Hosts.Get(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get host with ID: %s, %s", id, err)...)
	}
	if host == nil {
		return append(diags, diag.Errorf("Host with id %s was not found", id)...)
	}

	d.SetId(host.ID)
	d.Set("name", host.Name)
	d.Set("description", host.Description)
	d.Set("domain", host.Domain)
	d.Set("internet_access", host.InternetAccess)
	d.Set("system_subnets", host.SystemSubnets)
	return diags
}
