package cloudconnexa

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

func dataSourceAccessGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccessGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Access group name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Access group description.",
			},
			"source": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     resourceSource(),
			},
			"destination": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     resourceDestination(),
			},
		},
	}
}

func dataSourceAccessGroupRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var id = data.Get("id").(string)
	group, err := c.AccessGroups.Get(id)

	if err != nil {
		return diag.FromErr(err)
	}
	if group == nil {
		return append(diags, diag.Errorf("Access Group with id %s was not found", id)...)
	}
	setAccessGroupData(data, group)
	return nil
}
