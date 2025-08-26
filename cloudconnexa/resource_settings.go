package cloudconnexa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// resourceSettings returns a Terraform resource for managing CloudConnexa settings
func resourceSettings() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_settings` to define settings",
		CreateContext: resourceSettingsUpdate,
		ReadContext:   resourceSettingsRead,
		DeleteContext: resourceSettingsDelete,
		UpdateContext: resourceSettingsUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSettingsImport,
		},
		Schema: map[string]*schema.Schema{
			"allow_trusted_devices": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"two_factor_auth": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     dnsServersSchema(),
			},
			"default_dns_suffix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_proxy_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"dns_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     dnsZoneSchema(),
			},
			"connect_auth": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NO_AUTH", "ON_PRIOR_AUTH", "EVERY_TIME"}, false),
			},
			"device_allowance_per_user": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"device_allowance_force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"device_enforcement": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"OFF", "LEARN_AND_ENFORCE", "ENFORCE"}, false),
			},
			"profile_distribution": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"AUTOMATIC", "MANUAL"}, false),
			},
			"connection_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"client_options": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_routing_subnet": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     domainRoutingSubnet(),
			},
			"snat": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"subnet": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     subnetSchema(),
			},
			"topology": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"FULL_MESH", "CUSTOM"}, false),
			},
			"dns_log_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"access_visibility_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

// dnsServersSchema returns a Terraform schema for DNS server configuration
func dnsServersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"primary_ip_v4": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secondary_ip_v4": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// dnsZoneSchema returns a Terraform schema for DNS zone configuration
func dnsZoneSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// domainRoutingSubnet returns a Terraform schema for domain routing subnet configuration
func domainRoutingSubnet() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_v4_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_v6_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// subnetSchema returns a Terraform schema for subnet configuration
func subnetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_v4_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_v6_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// resourceSettingsUpdate handles the creation and update of CloudConnexa settings
func resourceSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	if d.HasChange("allow_trusted_devices") {
		_, err := c.Settings.SetTrustedDevicesAllowed(d.Get("allow_trusted_devices").(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.HasChange("two_factor_auth") && d.Get("two_factor_auth") != "" {
		_, err := c.Settings.SetTwoFactorAuthEnabled(d.Get("allow_trusted_devices").(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.HasChange("dns_servers") && len(d.Get("dns_servers").([]interface{})) > 0 {
		value := &cloudconnexa.DNSServers{}
		servers := d.Get("dns_servers").([]interface{})
		if len(servers) > 0 && servers[0] != nil {
			value.PrimaryIPV4 = servers[0].(map[string]interface{})["primary_ip_v4"].(string)
			value.SecondaryIPV4 = servers[0].(map[string]interface{})["secondary_ip_v4"].(string)
		}
		_, err := c.Settings.SetDNSServers(value)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.HasChange("default_dns_suffix") && d.Get("default_dns_suffix") != "" {
		_, err := c.Settings.SetDefaultDNSSuffix(d.Get("default_dns_suffix").(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.HasChange("dns_proxy_enabled") {
		_, err := c.Settings.SetDNSProxyEnabled(d.Get("dns_proxy_enabled").(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.HasChange("dns_zones") && len(d.Get("dns_zones").([]interface{})) > 0 {
		value := []cloudconnexa.DNSZone{}
		zones := d.Get("dns_zones").([]interface{})
		for _, zone := range zones {
			var addresses []string
			for _, addr := range zone.(map[string]interface{})["addresses"].([]interface{}) {
				addresses = append(addresses, addr.(string))
			}

			value = append(value, cloudconnexa.DNSZone{
				Name:      zone.(map[string]interface{})["name"].(string),
				Addresses: addresses,
			})
		}

		_, err := c.Settings.SetDNSZones(value)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("connect_auth") && d.Get("connect_auth") != "" {
		_, err := c.Settings.SetDefaultConnectAuth(d.Get("connect_auth").(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("device_allowance_per_user") && d.Get("device_allowance_per_user") != 0 {
		_, err := c.Settings.SetDefaultDeviceAllowancePerUser(d.Get("device_allowance_per_user").(int))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("device_allowance_force_update") {
		_, err := c.Settings.SetForceUpdateDeviceAllowanceEnabled(d.Get("device_allowance_force_update").(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("device_enforcement") && d.Get("device_enforcement") != "" {
		_, err := c.Settings.SetDeviceEnforcement(d.Get("device_enforcement").(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("profile_distribution") && d.Get("profile_distribution") != "" {
		_, err := c.Settings.SetProfileDistribution(d.Get("profile_distribution").(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("connection_timeout") && d.Get("connection_timeout") != 0 {
		_, err := c.Settings.SetConnectionTimeout(d.Get("connection_timeout").(int))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("client_options") && len(d.Get("client_options").([]interface{})) > 0 {
		_, err := c.Settings.SetClientOptions(toStrings(d.Get("client_options").([]interface{})))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("default_region") && d.Get("default_region") != "" {
		_, err := c.Settings.SetDefaultRegion(d.Get("default_region").(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("domain_routing_subnet") && len(d.Get("domain_routing_subnet").([]interface{})) > 0 {
		value := cloudconnexa.DomainRoutingSubnet{}
		subnet := d.Get("domain_routing_subnet").([]interface{})
		if len(subnet) > 0 && subnet[0] != nil {
			value.IPV4Address = subnet[0].(map[string]interface{})["ip_v4_address"].(string)
			value.IPV6Address = subnet[0].(map[string]interface{})["ip_v6_address"].(string)
		}
		_, err := c.Settings.SetDomainRoutingSubnet(value)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("snat") {
		_, err := c.Settings.SetSnatEnabled(d.Get("snat").(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("subnet") && len(d.Get("subnet").([]interface{})) > 0 {
		value := cloudconnexa.Subnet{}
		subnet := d.Get("subnet").([]interface{})
		if len(subnet) > 0 && subnet[0] != nil {
			for _, addr := range subnet[0].(map[string]interface{})["ip_v4_address"].([]interface{}) {
				value.IPV4Address = append(value.IPV4Address, addr.(string))
			}
			for _, addr := range subnet[0].(map[string]interface{})["ip_v6_address"].([]interface{}) {
				value.IPV6Address = append(value.IPV6Address, addr.(string))
			}
		}
		_, err := c.Settings.SetSubnet(value)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}

		if d.HasChange("topology") && d.Get("topology") != "" {
			_, err := c.Settings.SetTopology(d.Get("topology").(string))
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
		}
	}

	if d.HasChange("dns_log_enabled") {
		err := c.Settings.SetDNSLogEnabled(d.Get("dns_log_enabled").(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("access_visibility_enabled") {
		err := c.Settings.SetAccessVisibilityEnabled(d.Get("access_visibility_enabled").(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	d.SetId("settings")
	return diags
}

// resourceSettingsRead reads the settings resource from the CloudConnexa API
// and updates the Terraform state with the current values.
func resourceSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	allowTrustedDevices, err := c.Settings.GetTrustedDevicesAllowed()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("allow_trusted_devices", allowTrustedDevices)

	if len(d.Get("dns_servers").([]interface{})) > 0 {
		dnsServers, err := c.Settings.GetDNSServers()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		if dnsServers != nil && (dnsServers.PrimaryIPV4 != "" || dnsServers.SecondaryIPV4 != "") {
			value := make(map[string]interface{})
			value["primary_ip_v4"] = dnsServers.PrimaryIPV4
			value["secondary_ip_v4"] = dnsServers.SecondaryIPV4
			d.Set("dns_servers", []interface{}{value})
		}
	}

	if d.Get("default_dns_suffix") != "" {
		defaultDNSSuffix, err := c.Settings.GetDefaultDNSSuffix()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("default_dns_suffix", defaultDNSSuffix)
	}

	dnsProxyEnabled, err := c.Settings.GetDNSProxyEnabled()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("dns_proxy_enabled", dnsProxyEnabled)

	if len(d.Get("dns_zones").([]interface{})) > 0 {
		dnsZones, err := c.Settings.GetDNSZones()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		var zoneValues []interface{}
		for _, zone := range dnsZones {
			value := make(map[string]interface{})
			value["name"] = zone.Name
			value["addresses"] = zone.Addresses
			zoneValues = append(zoneValues, value)
		}
		d.Set("dns_zones", zoneValues)
	}

	if d.Get("connect_auth") != "" {
		connectAuth, err := c.Settings.GetDefaultConnectAuth()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("connect_auth", connectAuth)
	}

	if d.Get("device_allowance_per_user") != 0 {
		deviceAllowance, err := c.Settings.GetDefaultDeviceAllowancePerUser()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("device_allowance_per_user", deviceAllowance)
	}

	deviceAllowanceForceUpdate, err := c.Settings.GetForceUpdateDeviceAllowanceEnabled()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("device_allowance_force_update", deviceAllowanceForceUpdate)

	if d.Get("device_enforcement") != "" {
		deviceEnforcement, err := c.Settings.GetDeviceEnforcement()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("device_enforcement", deviceEnforcement)
	}

	if d.Get("profile_distribution") != "" {
		profileDistribution, err := c.Settings.GetProfileDistribution()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("profile_distribution", profileDistribution)
	}

	if d.Get("connection_timeout") != 0 {
		connectionTimeout, err := c.Settings.GetConnectionTimeout()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("connection_timeout", connectionTimeout)
	}

	// Get and set client options if they exist
	if len(d.Get("client_options").([]interface{})) > 0 {
		clientOptions, err := c.Settings.GetClientOptions()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("client_options", clientOptions)
	}

	// Get and set default region if specified
	if d.Get("default_region") != "" {
		defaultRegion, err := c.Settings.GetDefaultRegion()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("default_region", defaultRegion)
	}

	// Get and set domain routing subnet if specified
	if len(d.Get("domain_routing_subnet").([]interface{})) > 0 {
		domainRoutingSubnet, err := c.Settings.GetDomainRoutingSubnet()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		if domainRoutingSubnet.IPV4Address != "" || domainRoutingSubnet.IPV6Address != "" {
			value := make(map[string]interface{})
			value["ip_v4_address"] = domainRoutingSubnet.IPV4Address
			value["ip_v6_address"] = domainRoutingSubnet.IPV6Address
			d.Set("domain_routing_subnet", []interface{}{value})
		}
	}

	// Get and set SNAT enabled status
	snat, err := c.Settings.GetSnatEnabled()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("snat", snat)

	// Get and set subnet configuration if specified
	if len(d.Get("subnet").([]interface{})) > 0 {
		subnet, err := c.Settings.GetSubnet()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		subnetValue := make(map[string]interface{})
		subnetValue["ip_v4_address"] = subnet.IPV4Address
		subnetValue["ip_v6_address"] = subnet.IPV6Address
		d.Set("subnet", subnetValue)
	}

	// Get and set topology if specified
	if d.Get("topology") != "" {
		topology, err := c.Settings.GetTopology()
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("topology", topology)
	}

	// Get and set DNS Log enabled status
	dnsLog, err := c.Settings.GetDNSLogEnabled()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("dns_log_enabled", dnsLog)

	// Get and set Access Visibility enabled status
	accessVisibility, err := c.Settings.GetAccessVisibilityEnabled()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("access_visibility_enabled", accessVisibility)

	d.SetId("settings")
	return diags
}

// resourceSettingsDelete handles the deletion of CloudConnexa settings.
// Since settings are managed as a single resource, deletion is a no-op operation.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The provider meta interface
//
// Returns:
//   - diag.Diagnostics: Empty diagnostics since deletion is a no-op
func resourceSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// resourceSettingsImport handles the import of CloudConnexa settings.
// Since settings is a singleton resource, it always uses the fixed ID "settings"
// and ignores the provided import ID.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The provider meta interface
//
// Returns:
//   - []*schema.ResourceData: Resource data slice for imported settings
//   - error: Any error that occurred during import
func resourceSettingsImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Settings is a singleton resource with a fixed ID
	d.SetId("settings")

	// Read the current settings from the API to populate the state
	diags := resourceSettingsRead(ctx, d, m)
	if diags.HasError() {
		// Convert diagnostics to error for import function
		var err error
		for _, diagnostic := range diags {
			if diagnostic.Severity == diag.Error {
				if err == nil {
					err = fmt.Errorf("summary: %s; detail: %s", diagnostic.Summary, diagnostic.Detail)
				} else {
					err = fmt.Errorf("%s; summary: %s; detail: %s", err.Error(), diagnostic.Summary, diagnostic.Detail)
				}
			}
		}
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
