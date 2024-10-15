package cloudconnexa

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

func dataSourceLocationContext() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocationContextRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_groups_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"ip_policy": {
				Type:     schema.TypeList,
				Elem:     ipPolicyConfig(),
				Computed: true,
			},
			"country_policy": {
				Type:     schema.TypeList,
				Elem:     countryPolicyConfig(),
				Computed: true,
			},
			"default_policy": {
				Type:     schema.TypeList,
				Elem:     defaultPolicyConfig(),
				Computed: true,
			},
		},
	}
}

func dataSourceLocationContextRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var id = data.Get("id").(string)
	context, err := c.LocationContexts.Get(id)

	if err != nil {
		return diag.FromErr(err)
	}
	if context == nil {
		return append(diags, diag.Errorf("Location Context with ID %s was not found", id)...)
	}
	setLocationContextData(data, context)
	return nil
}
