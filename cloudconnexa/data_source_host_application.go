package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// dataSourceHostApplication returns a Terraform data source resource for CloudConnexa host applications.
// This resource allows users to read existing host application configurations.
func dataSourceHostApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostApplicationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application ID",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Application name",
				ExactlyOneOf: []string{"id", "name"},
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

// dataSourceHostApplicationRead handles the read operation for the host application data source.
// It retrieves the host application configuration from CloudConnexa using the provided ID
// and updates the Terraform state with the retrieved data.
//
// Parameters:
//   - ctx: The context for the operation
//   - data: The Terraform resource data
//   - i: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceHostApplicationRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var id = data.Get("id").(string)
	var application *cloudconnexa.ApplicationResponse
	var err error
	if id != "" {
		application, err = c.HostApplications.Get(id)
		if err != nil {
			return append(diags, diag.Errorf("Failed to get application with ID: %s, %s", id, err)...)
		}
		if application == nil {
			return append(diags, diag.Errorf("Application with id %s was not found", id)...)
		}
	} else {
		var name = data.Get("name").(string)
		application, err = c.HostApplications.GetByName(name)
		if err != nil {
			return append(diags, diag.Errorf("Failed to get application with name: %s, %s", name, err)...)
		}
		if application == nil {
			return append(diags, diag.Errorf("Application with name %s was not found", name)...)
		}
	}
	setApplicationData(data, application)
	return nil
}
