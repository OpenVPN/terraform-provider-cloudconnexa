package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// dataSourceNetworkApplication returns a Terraform data source resource for network applications.
// This resource allows users to read information about a specific network application by its ID.
func dataSourceNetworkApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkApplicationRead,
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
				Elem:     resourceNetworkApplicationRoute(),
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceNetworkApplicationConfig(),
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceNetworkApplicationRead handles the read operation for the network application data source.
// It retrieves application details using the provided ID and sets the data in the Terraform state.
// Parameters:
//   - ctx: Context for the operation
//   - data: The Terraform resource data
//   - i: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors or warnings
func dataSourceNetworkApplicationRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	meta := i.(*providerMeta)
	c := meta.Client
	var diags diag.Diagnostics
	var id = data.Get("id").(string)
	var application *cloudconnexa.NetworkApplicationResponse
	var err error
	if id != "" {
		application, err = withRetry(ctx, meta.RetryConfig, func() (*cloudconnexa.NetworkApplicationResponse, error) {
			return c.NetworkApplications.Get(id)
		})
	} else {
		var name = data.Get("name").(string)
		application, err = withRetry(ctx, meta.RetryConfig, func() (*cloudconnexa.NetworkApplicationResponse, error) {
			return c.NetworkApplications.GetByName(name)
		})
	}
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if application == nil {
		return append(diags, diag.Errorf("Network application was not found")...)
	}
	setNetworkApplicationData(data, application)
	return nil
}
