package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// dataSourceLocationContext returns a Terraform data source resource for CloudConnexa location contexts.
// This resource allows users to read existing location context configurations.
func dataSourceLocationContext() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocationContextRead,
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
			"user_groups_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"ip_check": {
				Type:     schema.TypeList,
				Elem:     ipCheckConfig(),
				Computed: true,
			},
			"country_check": {
				Type:     schema.TypeList,
				Elem:     countryCheckConfig(),
				Computed: true,
			},
			"default_check": {
				Type:     schema.TypeList,
				Elem:     defaultCheckConfig(),
				Computed: true,
			},
		},
	}
}

// dataSourceLocationContextRead handles the read operation for the location context data source.
// It retrieves the location context configuration from CloudConnexa using the provided ID
// and updates the Terraform state with the retrieved data.
//
// Parameters:
//   - ctx: The context for the operation
//   - data: The Terraform resource data
//   - i: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceLocationContextRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	meta := i.(*providerMeta)
	c := meta.Client
	var diags diag.Diagnostics
	var id = data.Get("id").(string)
	var lc *cloudconnexa.LocationContext
	var err error
	if id != "" {
		lc, err = withRetry(ctx, meta.RetryConfig, func() (*cloudconnexa.LocationContext, error) {
			return c.LocationContexts.Get(id)
		})
	} else {
		var name = data.Get("name").(string)
		lc, err = withRetry(ctx, meta.RetryConfig, func() (*cloudconnexa.LocationContext, error) {
			return c.LocationContexts.GetByName(name)
		})
	}
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if lc == nil {
		return append(diags, diag.Errorf("Location context was not found")...)
	}
	setLocationContextData(data, lc)
	return nil
}
