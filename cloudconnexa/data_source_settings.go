package cloudconnexa

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceSettings returns a Terraform data source schema for CloudConnexa settings.
// It defines the schema for retrieving various configuration settings including DNS,
// authentication, device management, and network topology settings.
//
// Returns:
//   - *schema.Resource: A configured Terraform data source for CloudConnexa settings
func dataSourceSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceSettingsRead,
		Schema: map[string]*schema.Schema{
			"allow_trusted_devices": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"two_factor_auth": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dnsServersSchema(),
			},
			"default_dns_suffix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_proxy_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dns_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dnsZoneSchema(),
			},
			"connect_auth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_allowance_per_user": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"device_allowance_force_update": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"device_enforcement": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_distribution": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"client_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_routing_subnet": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     domainRoutingSubnetSchema(),
			},
			"snat": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"subnet": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     subnetSchema(),
			},
			"topology": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_log_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether DNS logging is enabled.",
			},
			"access_visibility_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether access visibility is enabled.",
			},
		},
	}
}
