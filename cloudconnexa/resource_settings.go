package cloudconnexa

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

func resourceSettings() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_settings` to define settings",
		CreateContext: resourceSettingsUpdate,
		ReadContext:   resourceSettingsRead,
		DeleteContext: resourceSettingsDelete,
		UpdateContext: resourceSettingsUpdate,
		Schema: map[string]*schema.Schema{
			"allow_trusted_devices": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"two_factor_auth": {
				Type:     schema.TypeBool,
				Required: true,
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
				Required: true,
			},
			"dns_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     dnsZoneSchema(),
			},
			"connect_auth": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"NO_AUTH", "ON_PRIOR_AUTH", "EVERY_TIME"}, false),
			},
			"device_allowance_per_user": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"device_allowance_force_update": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"device_enforcement": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OFF", "LEARN_AND_ENFORCE", "ENFORCE"}, false),
			},
			"profile_distribution": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AUTOMATIC", "MANUAL"}, false),
			},
			"connection_timeout": {
				Type:     schema.TypeInt,
				Required: true,
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
				Required: true,
			},
			"domain_routing_subnet": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem:     domainRoutingSubnet(),
			},
			"snat": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"subnet": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem:     subnetSchema(),
			},
			"topology": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"FULL_MESH", "CUSTOM"}, false),
			},
		},
	}
}

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

func resourceSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	if d.HasChange("allow_trusted_devices") {
		_, err := c.Settings.SetTrustedDevicesAllowed(d.Get("allow_trusted_devices").(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.HasChange("two_factor_auth") {
		_, err := c.Settings.SetTwoFactorAuthEnabled(d.Get("allow_trusted_devices").(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.HasChange("dns_servers") {
		value := &cloudconnexa.DnsServers{}
		servers := d.Get("dns_servers").([]interface{})
		if len(servers) > 0 && servers[0] != nil {
			value.PrimaryIpV4 = servers[0].(map[string]interface{})["primary_ip_v4"].(string)
			value.SecondaryIpV4 = servers[0].(map[string]interface{})["secondary_ip_v4"].(string)
		}
		_, err := c.Settings.SetDnsServers(value)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.HasChange("default_dns_suffix") {
		_, err := c.Settings.SetDefaultDnsSuffix(d.Get("default_dns_suffix").(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.HasChange("dns_proxy_enabled") {
		_, err := c.Settings.SetDnsProxyAuthEnabled(d.Get("dns_proxy_enabled").(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.HasChange("dns_zones") {
		value := []cloudconnexa.DnsZone{}
		zones := d.Get("dns_zones").([]interface{})
		for _, zone := range zones {
			var addresses []string
			for _, addr := range zone.(map[string]interface{})["addresses"].([]interface{}) {
				addresses = append(addresses, addr.(string))
			}

			value = append(value, cloudconnexa.DnsZone{
				zone.(map[string]interface{})["name"].(string),
				addresses,
			})
		}

		_, err := c.Settings.SetDnsZones(value)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("connect_auth") {
		_, err := c.Settings.SetDefaultConnectAuth(d.Get("connect_auth").(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("device_allowance_per_user") {
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

	if d.HasChange("device_enforcement") {
		_, err := c.Settings.SetDeviceEnforcement(d.Get("device_enforcement").(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("profile_distribution") {
		_, err := c.Settings.SetProfileDistribution(d.Get("profile_distribution").(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("connection_timeout") {
		_, err := c.Settings.SetConnectionTimeout(d.Get("connection_timeout").(int))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("client_options") {
		_, err := c.Settings.SetClientOptions(d.Get("client_options").([]string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("default_region") {
		_, err := c.Settings.SetDefaultRegion(d.Get("default_region").(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if d.HasChange("domain_routing_subnet") {
		value := cloudconnexa.DomainRoutingSubnet{}
		subnet := d.Get("domain_routing_subnet").([]interface{})
		if len(subnet) > 0 && subnet[0] != nil {
			value.IpV4Address = subnet[0].(map[string]interface{})["ip_v4_address"].(string)
			value.IpV6Address = subnet[0].(map[string]interface{})["ip_v6_address"].(string)
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

	if d.HasChange("subnet") {
		value := cloudconnexa.Subnet{}
		subnet := d.Get("subnet").([]interface{})
		if len(subnet) > 0 && subnet[0] != nil {
			for _, addr := range subnet[0].(map[string]interface{})["ip_v4_address"].([]interface{}) {
				value.IpV4Address = append(value.IpV4Address, addr.(string))
			}
			for _, addr := range subnet[0].(map[string]interface{})["ip_v6_address"].([]interface{}) {
				value.IpV6Address = append(value.IpV6Address, addr.(string))
			}
		}
		_, err := c.Settings.SetSubnet(value)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}

		if d.HasChange("topology") {
			_, err := c.Settings.SetTopology(d.Get("topology").(string))
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
		}
	}

	d.SetId("settings")
	return diags
}

func resourceSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	allowTrustedDevices, err := c.Settings.GetTrustedDevicesAllowed()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("allow_trusted_devices", allowTrustedDevices)

	dnsServers, err := c.Settings.GetDnsServers()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if dnsServers != nil && (dnsServers.PrimaryIpV4 != "" || dnsServers.SecondaryIpV4 != "") {
		value := make(map[string]interface{})
		value["primary_ip_v4"] = dnsServers.PrimaryIpV4
		value["secondary_ip_v4"] = dnsServers.SecondaryIpV4
		d.Set("dns_servers", []interface{}{value})
	}

	defaultDnsSuffix, err := c.Settings.GetDefaultDnsSuffix()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("default_dns_suffix", defaultDnsSuffix)

	dnsProxyEnabled, err := c.Settings.GetDnsProxyEnabled()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("dns_proxy_enabled", dnsProxyEnabled)

	dnsZones, err := c.Settings.GetDnsZones()
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

	connectAuth, err := c.Settings.GetDefaultConnectAuth()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("connect_auth", connectAuth)

	deviceAllowance, err := c.Settings.GetDefaultDeviceAllowancePerUser()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("device_allowance_per_user", deviceAllowance)

	deviceAllowanceForceUpdate, err := c.Settings.GetForceUpdateDeviceAllowanceEnabled()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("device_allowance_force_update", deviceAllowanceForceUpdate)

	deviceEnforcement, err := c.Settings.GetDeviceEnforcement()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("device_enforcement", deviceEnforcement)

	profileDistribution, err := c.Settings.GetProfileDistribution()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("profile_distribution", profileDistribution)

	connectionTimeout, err := c.Settings.GetConnectionTimeout()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("connection_timeout", connectionTimeout)

	clientOptions, err := c.Settings.GetClientOptions()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("client_options", clientOptions)

	defaultRegion, err := c.Settings.GetDefaultRegion()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("default_region", defaultRegion)

	domainRoutingSubnet, err := c.Settings.GetDomainRoutingSubnet()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if domainRoutingSubnet.IpV4Address != "" || domainRoutingSubnet.IpV6Address != "" {
		value := make(map[string]interface{})
		value["ip_v4_address"] = domainRoutingSubnet.IpV4Address
		value["ip_v6_address"] = domainRoutingSubnet.IpV6Address
		d.Set("domain_routing_subnet", []interface{}{value})
	}

	snat, err := c.Settings.GetSnatEnabled()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("snat", snat)

	subnet, err := c.Settings.GetSubnet()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	subnetValue := make(map[string]interface{})
	subnetValue["ip_v4_address"] = subnet.IpV4Address
	subnetValue["ip_v6_address"] = subnet.IpV6Address
	d.Set("subnet", subnetValue)

	topology, err := c.Settings.GetTopology()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("topology", topology)

	d.SetId("settings")
	return diags
}

func resourceSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
