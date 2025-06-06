package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceNetwork returns a Terraform data source resource for CloudConnexa networks.
// This resource allows users to read existing network configurations.
//
// Returns:
//   - *schema.Resource: A Terraform resource definition for network data source
func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Use a `cloudconnexa_network` data source to read an CloudConnexa network.",
		ReadContext: dataSourceNetworkRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the network.",
			},
			"egress": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Boolean to indicate whether this network provides an egress or not.",
			},
			"internet_access": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of internet access provided. Valid values are `SPLIT_TUNNEL_ON`, `SPLIT_TUNNEL_OFF`, or `RESTRICTED_INTERNET`. Defaults to `SPLIT_TUNNEL_ON`.",
			},
			"system_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IPV4 and IPV6 subnets automatically assigned to this network.",
			},
		},
	}
}

// dataSourceNetworkRead handles the read operation for the network data source.
// It retrieves the network configuration from CloudConnexa using the provided ID
// and updates the Terraform state with the retrieved data.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var network *cloudconnexa.Network
	var err error
	id := d.Get("id").(string)
	network, err = c.Networks.Get(id)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if network == nil {
		return append(diags, diag.Errorf("Network with id %s was not found", id)...)
	}
	d.SetId(network.ID)
	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("egress", network.Egress)
	d.Set("internet_access", network.InternetAccess)
	d.Set("system_subnets", network.SystemSubnets)
	return diags
}
