package cloudconnexa

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// dataSourceHostIPService returns a Terraform data source resource for CloudConnexa host IP services.
// This resource allows users to read existing host IP service configurations.
func dataSourceHostIPService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostIPServiceRead,
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceHostIpServiceConfig(),
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceHostIPServiceRead handles the read operation for the host IP service data source.
// It retrieves the host IP service configuration from CloudConnexa using the provided ID
// and updates the Terraform state with the retrieved data.
//
// Parameters:
//   - ctx: The context for the operation
//   - data: The Terraform resource data
//   - i: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceHostIPServiceRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	service, err := c.HostIPServices.Get(data.Get("id").(string))

	if err != nil {
		return diag.FromErr(err)
	}
	setHostIpServiceResourceData(data, service)
	return nil
}
