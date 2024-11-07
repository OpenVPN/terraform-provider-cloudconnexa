package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

func dataSourceApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "Application ID",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "Application name",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceApplicationRoute(),
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceApplicationConfig(),
			},
			"network_item_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_item_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApplicationRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var application *cloudconnexa.ApplicationResponse
	var err error
	applicationId := data.Get("id").(string)
	applicationName := data.Get("name").(string)
	if applicationId != "" {
		application, err = c.Applications.Get(applicationId)
		if err != nil {
			return diag.FromErr(err)
		}
		if application == nil {
			return append(diags, diag.Errorf("Application with id %s was not found", applicationId)...)
		}
	} else if applicationName != "" {
		applicationsAll, err := c.Applications.List()
		var applicationCount int
		if err != nil {
			return diag.FromErr(err)
		}

		for _, app := range applicationsAll {
			if app.Name == applicationName {
				applicationCount++
			}
		}

		if applicationCount == 0 {
			return append(diags, diag.Errorf("Application with name %s was not found", applicationName)...)
		} else if applicationCount > 1 {
			return append(diags, diag.Errorf("More than 1 application with name %s was found. Please use id instead", applicationName)...)
		} else {
			application, err = c.Applications.GetByName(applicationName)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		return append(diags, diag.Errorf("Application name or id is missing")...)
	}
	setApplicationData(data, application)
	return nil
}
