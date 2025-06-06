package cloudconnexa

import (
	"context"
	"slices"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// validValues contains the list of allowed service types for network IP services
var (
	validValues = []string{"ANY", "BGP", "DHCP", "DNS", "FTP", "HTTP", "HTTPS", "IMAP", "IMAPS", "NTP", "POP3", "POP3S", "SMTP", "SMTPS", "SNMP", "SSH", "TELNET", "TFTP"}
)

// resourceNetworkIPService returns a Terraform resource schema for managing network IP services
func resourceNetworkIPService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkIpServiceCreate,
		ReadContext:   resourceNetworkIpServiceRead,
		DeleteContext: resourceNetworkIpServiceDelete,
		UpdateContext: resourceNetworkIpServiceUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 255),
				Optional:     true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"IP_SOURCE", "SERVICE_DESTINATION"}, false),
			},
			"routes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     resourceNetworkIpServiceConfig(),
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// resourceNetworkIpServiceUpdate updates an existing network IP service
func resourceNetworkIpServiceUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)

	s, err := c.NetworkIPServices.Update(data.Id(), resourceDataToNetworkIpService(data))
	if err != nil {
		return diag.FromErr(err)
	}
	setNetworkIpServiceResourceData(data, s)
	return nil
}

// resourceNetworkIpServiceConfig returns the schema for network IP service configuration
func resourceNetworkIpServiceConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"custom_service_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: customNetworkIpServiceTypesConfig(),
				},
			},
			"service_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						val := i.(string)
						for _, validValue := range validValues {
							if val == validValue {
								return nil
							}
						}
						return diag.Errorf("service type must be one of %s", validValues)
					},
				},
			},
		},
	}
}

// customNetworkIpServiceTypesConfig returns the schema for custom network IP service types
func customNetworkIpServiceTypesConfig() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false),
		},
		"from_port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"to_port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}

// resourceNetworkIpServiceRead retrieves a network IP service by ID
func resourceNetworkIpServiceRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	service, err := c.NetworkIPServices.Get(data.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if service == nil {
		data.SetId("")
		return diags
	}
	setNetworkIpServiceResourceData(data, service)
	return diags
}

// setNetworkIpServiceResourceData sets the resource data from a network IP service response
func setNetworkIpServiceResourceData(data *schema.ResourceData, service *cloudconnexa.IPServiceResponse) {
	data.SetId(service.ID)
	_ = data.Set("name", service.Name)
	_ = data.Set("description", service.Description)
	_ = data.Set("type", service.Type)
	_ = data.Set("routes", flattenNetworkIpServiceRoutes(service.Routes))
	_ = data.Set("config", flattenNetworkIpServiceConfig(service.Config))
	_ = data.Set("network_id", service.NetworkItemID)
}

// resourceNetworkIpServiceDelete deletes a network IP service
func resourceNetworkIpServiceDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	err := c.NetworkIPServices.Delete(data.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

// flattenNetworkIpServiceConfig flattens the network IP service configuration for Terraform state
func flattenNetworkIpServiceConfig(config *cloudconnexa.IPServiceConfig) interface{} {
	var data = map[string]interface{}{
		"custom_service_types": flattenNetworkIpServiceCustomServiceTypes(config.CustomServiceTypes),
		"service_types":        removeElement(config.ServiceTypes, "CUSTOM"),
	}
	return []interface{}{data}
}

// flattenNetworkIpServiceCustomServiceTypes flattens custom service types for Terraform state
func flattenNetworkIpServiceCustomServiceTypes(types []*cloudconnexa.CustomIPServiceType) interface{} {
	var cst []interface{}
	for _, t := range types {
		var ports = append(t.Port, t.IcmpType...)
		if len(ports) > 0 {
			for _, port := range ports {
				cst = append(cst, map[string]interface{}{
					"protocol":  t.Protocol,
					"from_port": port.LowerValue,
					"to_port":   port.UpperValue,
				})
			}
		} else {
			cst = append(cst, map[string]interface{}{
				"protocol": t.Protocol,
			})
		}
	}
	return cst
}

// flattenNetworkIpServiceRoutes flattens network IP service routes for Terraform state
func flattenNetworkIpServiceRoutes(routes []*cloudconnexa.Route) []string {
	var data []string
	for _, route := range routes {
		data = append(
			data,
			route.Subnet,
		)
	}
	return data
}

// resourceNetworkIpServiceCreate creates a new network IP service
func resourceNetworkIpServiceCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*cloudconnexa.Client)

	service := resourceDataToNetworkIpService(data)
	createdService, err := client.NetworkIPServices.Create(service)
	if err != nil {
		return diag.FromErr(err)
	}
	setNetworkIpServiceResourceData(data, createdService)
	return nil
}

// resourceDataToNetworkIpService converts Terraform resource data to a CloudConnexa network IP service configuration.
// It maps the Terraform schema data to the corresponding CloudConnexa API structure.
func resourceDataToNetworkIpService(data *schema.ResourceData) *cloudconnexa.IPService {
	// Extract and process routes from Terraform data
	routes := data.Get("routes").([]interface{})
	var configRoutes []*cloudconnexa.IPServiceRoute
	for _, r := range routes {
		configRoutes = append(
			configRoutes,
			&cloudconnexa.IPServiceRoute{
				Value:       r.(string),
				Description: "Managed by Terraform",
			},
		)
	}

	// Initialize and process service configuration
	config := cloudconnexa.IPServiceConfig{}
	configList := data.Get("config").([]interface{})
	if len(configList) > 0 && configList[0] != nil {
		config.CustomServiceTypes = []*cloudconnexa.CustomIPServiceType{}
		config.ServiceTypes = []string{}

		// Process main configuration block
		mainConfig := configList[0].(map[string]interface{})
		var cst = mainConfig["custom_service_types"].(*schema.Set)
		var groupedCst = make(map[string][]cloudconnexa.Range)

		// Group custom service types by protocol
		for _, item := range cst.List() {
			var cstItem = item.(map[string]interface{})
			var protocol = cstItem["protocol"].(string)
			var fromPort = cstItem["from_port"].(int)
			var toPort = cstItem["to_port"].(int)

			if groupedCst[protocol] == nil {
				groupedCst[protocol] = make([]cloudconnexa.Range, 0)
			}
			if fromPort > 0 || toPort > 0 {
				groupedCst[protocol] = append(groupedCst[protocol], cloudconnexa.Range{
					LowerValue: fromPort,
					UpperValue: toPort,
				})
			}
		}

		// Process grouped custom service types
		for protocol, ports := range groupedCst {
			if protocol == "ICMP" {
				config.CustomServiceTypes = append(
					config.CustomServiceTypes,
					&cloudconnexa.CustomIPServiceType{
						Protocol: protocol,
						IcmpType: ports,
					},
				)
			} else {
				config.CustomServiceTypes = append(
					config.CustomServiceTypes,
					&cloudconnexa.CustomIPServiceType{
						Protocol: protocol,
						Port:     ports,
					},
				)
			}
		}

		// Process service types and ensure CUSTOM is included if needed
		for _, r := range mainConfig["service_types"].([]interface{}) {
			config.ServiceTypes = append(config.ServiceTypes, r.(string))
		}
		if len(config.CustomServiceTypes) > 0 && !slices.Contains(config.ServiceTypes, "CUSTOM") {
			config.ServiceTypes = append(config.ServiceTypes, "CUSTOM")
		}
	}

	// Create and return the final IP service configuration
	s := &cloudconnexa.IPService{
		Name:            data.Get("name").(string),
		Description:     data.Get("description").(string),
		NetworkItemID:   data.Get("network_id").(string),
		NetworkItemType: "NETWORK",
		Type:            data.Get("type").(string),
		Routes:          configRoutes,
		Config:          &config,
	}
	return s
}
