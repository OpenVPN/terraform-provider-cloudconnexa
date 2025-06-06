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

// resourceHostIPService returns a Terraform resource schema for managing host IP services
func resourceHostIPService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostIpServiceCreate,
		ReadContext:   resourceHostIpServiceRead,
		DeleteContext: resourceHostIpServiceDelete,
		UpdateContext: resourceHostIpServiceUpdate,
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
				Elem:     resourceHostIpServiceConfig(),
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// resourceHostIpServiceUpdate updates an existing host IP service
func resourceHostIpServiceUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)

	s, err := c.HostIPServices.Update(data.Id(), resourceDataToHostIpService(data))
	if err != nil {
		return diag.FromErr(err)
	}
	setHostIpServiceResourceData(data, s)
	return nil
}

// resourceHostIpServiceConfig returns the schema for host IP service configuration
func resourceHostIpServiceConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"custom_service_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: customHostServiceTypesConfig(),
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

// customHostServiceTypesConfig returns the schema for custom host service types configuration
func customHostServiceTypesConfig() map[string]*schema.Schema {
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

// resourceHostIpServiceRead reads the state of a host IP service
func resourceHostIpServiceRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	service, err := c.HostIPServices.Get(data.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if service == nil {
		data.SetId("")
		return diags
	}
	setHostIpServiceResourceData(data, service)
	return diags
}

// setHostIpServiceResourceData sets the resource data from a host IP service response
func setHostIpServiceResourceData(data *schema.ResourceData, service *cloudconnexa.IPServiceResponse) {
	data.SetId(service.ID)
	_ = data.Set("name", service.Name)
	_ = data.Set("description", service.Description)
	_ = data.Set("routes", flattenHostIpServiceRoutes(service.Routes))
	_ = data.Set("config", flattenHostServiceConfig(service.Config))
	_ = data.Set("host_id", service.NetworkItemID)
}

// resourceHostIpServiceDelete deletes a host IP service
func resourceHostIpServiceDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	err := c.HostIPServices.Delete(data.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

// flattenHostServiceConfig flattens the host service configuration into a format suitable for Terraform
func flattenHostServiceConfig(config *cloudconnexa.IPServiceConfig) interface{} {
	var data = map[string]interface{}{
		"custom_service_types": flattenCustomHostServiceTypes(config.CustomServiceTypes),
		"service_types":        removeElement(config.ServiceTypes, "CUSTOM"),
	}
	return []interface{}{data}
}

// flattenCustomHostServiceTypes flattens custom host service types into a format suitable for Terraform
func flattenCustomHostServiceTypes(types []*cloudconnexa.CustomIPServiceType) interface{} {
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

// flattenHostIpServiceRoutes flattens host IP service routes into a slice of strings
func flattenHostIpServiceRoutes(routes []*cloudconnexa.Route) []string {
	var data []string
	for _, route := range routes {
		data = append(
			data,
			route.Subnet,
		)
	}
	return data
}

// resourceHostIpServiceCreate creates a new host IP service
func resourceHostIpServiceCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*cloudconnexa.Client)

	service := resourceDataToHostIpService(data)
	createdService, err := client.HostIPServices.Create(service)
	if err != nil {
		return diag.FromErr(err)
	}
	setHostIpServiceResourceData(data, createdService)
	return nil
}

// resourceDataToHostIpService converts Terraform resource data to a host IP service configuration
func resourceDataToHostIpService(data *schema.ResourceData) *cloudconnexa.IPService {
	// Get routes from Terraform data and convert to IPServiceRoute slice
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

	// Initialize service configuration
	config := cloudconnexa.IPServiceConfig{}
	configList := data.Get("config").([]interface{})
	if len(configList) > 0 && configList[0] != nil {
		// Initialize empty slices for custom service types and service types
		config.CustomServiceTypes = []*cloudconnexa.CustomIPServiceType{}
		config.ServiceTypes = []string{}

		// Get main configuration from Terraform data
		mainConfig := configList[0].(map[string]interface{})
		var cst = mainConfig["custom_service_types"].(*schema.Set)
		// Group custom service types by protocol
		var groupedCst = make(map[string][]cloudconnexa.Range)
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

		// Add service types from configuration
		for _, r := range mainConfig["service_types"].([]interface{}) {
			config.ServiceTypes = append(config.ServiceTypes, r.(string))
		}
		// Add CUSTOM service type if custom service types exist
		if len(config.CustomServiceTypes) > 0 && !slices.Contains(config.ServiceTypes, "CUSTOM") {
			config.ServiceTypes = append(config.ServiceTypes, "CUSTOM")
		}
	}

	// Create and return the IP service configuration
	s := &cloudconnexa.IPService{
		Name:            data.Get("name").(string),
		Description:     data.Get("description").(string),
		NetworkItemID:   data.Get("host_id").(string),
		NetworkItemType: "HOST",
		Type:            "SERVICE_DESTINATION",
		Routes:          configRoutes,
		Config:          &config,
	}
	return s
}
