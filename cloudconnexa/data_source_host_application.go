package cloudconnexa

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

func dataSourceHostApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostApplicationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Application ID",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Application name",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceHostApplicationRoute(),
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceHostApplicationConfig(),
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceHostApplicationRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var application *cloudconnexa.ApplicationResponse
	var err error
	id := data.Get("id").(string)
	application, err = c.HostApplications.Get(id)
	if err != nil {
		if strings.Contains(err.Error(), "status code: 404") {
			return append(diags, diag.Errorf("Application with id %s was not found", id)...)
		} else {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if application == nil {
		return append(diags, diag.Errorf("Application with id %s was not found", id)...)
	}
	setApplicationData(data, application)
	return nil
}
