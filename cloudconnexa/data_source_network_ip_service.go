package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// dataSourceNetworkIPService returns a Terraform data source resource for CloudConnexa network IP services.
// This resource allows users to read existing network IP service configurations.
func dataSourceNetworkIPService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkIPServiceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
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
				Elem:     resourceNetworkIpServiceConfig(),
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceNetworkIPServiceRead handles the read operation for the network IP service data source.
// It retrieves the network IP service configuration from CloudConnexa using the provided ID
// and updates the Terraform state with the retrieved data.
//
// Parameters:
//   - ctx: The context for the operation
//   - data: The Terraform resource data
//   - i: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceNetworkIPServiceRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var id = data.Get("id").(string)
	var service *cloudconnexa.NetworkIPServiceResponse
	var err error
	if id != "" {
		service, err = c.NetworkIPServices.Get(id)
		if err != nil {
			return append(diags, diag.Errorf("Failed to get network IP service with ID: %s, %s", id, err)...)
		}
		if service == nil {
			return append(diags, diag.Errorf("Network IP service with id %s was not found", id)...)
		}
	} else {
		var name = data.Get("name").(string)
		service, err = c.NetworkIPServices.GetByName(name)
		if err != nil {
			return append(diags, diag.Errorf("Failed to get network IP service with name: %s, %s", name, err)...)
		}
		if service == nil {
			return append(diags, diag.Errorf("Network IP service with name %s was not found", name)...)
		}
	}
	setNetworkIpServiceResourceData(data, service)
	return nil
}
