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

// resourceHostApplication returns a Terraform resource schema for managing host applications
func resourceHostApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostApplicationCreate,
		ReadContext:   resourceHostApplicationRead,
		DeleteContext: resourceHostApplicationDelete,
		UpdateContext: resourceHostApplicationUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 40),
			},
			"description": {
				Type:         schema.TypeString,
				Default:      "Managed by Terraform",
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 120),
			},
			"routes": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     resourceHostApplicationRoute(),
			},
			"config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     resourceHostApplicationConfig(),
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// resourceHostApplicationUpdate updates an existing host application
func resourceHostApplicationUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)

	s, err := c.HostApplications.Update(data.Id(), resourceDataToHostApplication(data))
	if err != nil {
		return diag.FromErr(err)
	}
	setApplicationData(data, s)
	return nil
}

// resourceHostApplicationRoute returns the schema for host application routes
func resourceHostApplicationRoute() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allow_embedded_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

// resourceHostApplicationConfig returns the schema for host application configuration
func resourceHostApplicationConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"custom_service_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: customApplicationTypesConfig(),
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

// customApplicationTypesConfig returns the schema for custom application type configuration
func customApplicationTypesConfig() map[string]*schema.Schema {
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

// resourceHostApplicationRead reads the state of a host application
func resourceHostApplicationRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := data.Id()
	application, err := c.HostApplications.Get(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get host application with ID: %s, %s", id, err)...)
	}
	if application == nil {
		data.SetId("")
		return diags
	}
	setApplicationData(data, application)
	return diags
}

// setApplicationData sets the Terraform state data from an application response
func setApplicationData(data *schema.ResourceData, application *cloudconnexa.ApplicationResponse) {
	data.SetId(application.ID)
	_ = data.Set("name", application.Name)
	_ = data.Set("description", application.Description)
	_ = data.Set("routes", flattenApplicationRoutes(application.Routes))
	_ = data.Set("config", flattenApplicationConfig(application.Config))
	_ = data.Set("host_id", application.NetworkItemID)
}

// resourceHostApplicationDelete deletes a host application
func resourceHostApplicationDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	err := c.HostApplications.Delete(data.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

// flattenApplicationConfig converts an application config to Terraform state format
func flattenApplicationConfig(config *cloudconnexa.ApplicationConfig) interface{} {
	var data = map[string]interface{}{
		"custom_service_types": flattenApplicationCustomTypes(config.CustomServiceTypes),
		"service_types":        removeElement(config.ServiceTypes, "CUSTOM"),
	}
	return []interface{}{data}
}

// flattenApplicationCustomTypes converts custom application types to Terraform state format
func flattenApplicationCustomTypes(types []*cloudconnexa.CustomApplicationType) interface{} {
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

// flattenApplicationRoutes converts application routes to Terraform state format
func flattenApplicationRoutes(routes []*cloudconnexa.Route) []map[string]interface{} {
	var data []map[string]interface{}
	for _, route := range routes {
		data = append(data, map[string]interface{}{
			"domain":            route.Domain,
			"allow_embedded_ip": route.AllowEmbeddedIP,
		})
	}
	return data
}

// resourceHostApplicationCreate creates a new host application
func resourceHostApplicationCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*cloudconnexa.Client)

	application := resourceDataToHostApplication(data)
	createdApplication, err := client.HostApplications.Create(application)
	if err != nil {
		return diag.FromErr(err)
	}
	setApplicationData(data, createdApplication)
	return nil
}

// resourceDataToHostApplication converts Terraform resource data to a host application.
// It processes routes, custom service types, and service types from the Terraform state
// and converts them into a CloudConnexa Application structure.
//
// Parameters:
//   - data: The Terraform resource data containing the application configuration
//
// Returns:
//   - *cloudconnexa.Application: A pointer to the created CloudConnexa Application
func resourceDataToHostApplication(data *schema.ResourceData) *cloudconnexa.Application {
	routes := data.Get("routes").([]interface{})
	var configRoutes []*cloudconnexa.ApplicationRoute
	for _, r := range routes {
		var route = r.(map[string]interface{})
		configRoutes = append(
			configRoutes,
			&cloudconnexa.ApplicationRoute{
				Value:           route["domain"].(string),
				AllowEmbeddedIP: route["allow_embedded_ip"].(bool),
			},
		)
	}

	config := cloudconnexa.ApplicationConfig{}
	configList := data.Get("config").([]interface{})
	if len(configList) > 0 && configList[0] != nil {
		config.CustomServiceTypes = []*cloudconnexa.CustomApplicationType{}
		config.ServiceTypes = []string{}

		mainConfig := configList[0].(map[string]interface{})
		var cst = mainConfig["custom_service_types"].(*schema.Set)
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

		for protocol, ports := range groupedCst {
			if protocol == "ICMP" {
				config.CustomServiceTypes = append(
					config.CustomServiceTypes,
					&cloudconnexa.CustomApplicationType{
						Protocol: protocol,
						IcmpType: ports,
					},
				)
			} else {
				config.CustomServiceTypes = append(
					config.CustomServiceTypes,
					&cloudconnexa.CustomApplicationType{
						Protocol: protocol,
						Port:     ports,
					},
				)
			}
		}

		for _, r := range mainConfig["service_types"].([]interface{}) {
			config.ServiceTypes = append(config.ServiceTypes, r.(string))
		}
		if len(config.CustomServiceTypes) > 0 && !slices.Contains(config.ServiceTypes, "CUSTOM") {
			config.ServiceTypes = append(config.ServiceTypes, "CUSTOM")
		}
	}

	s := &cloudconnexa.Application{
		Name:            data.Get("name").(string),
		Description:     data.Get("description").(string),
		NetworkItemID:   data.Get("host_id").(string),
		NetworkItemType: "HOST",
		Routes:          configRoutes,
		Config:          &config,
	}
	return s
}

// removeElement removes a specific element from a string slice.
//
// Parameters:
//   - slice: The input slice of strings
//   - element: The string element to remove
//
// Returns:
//   - []string: A new slice with the specified element removed
func removeElement(slice []string, element string) []string {
	var result []string
	for _, item := range slice {
		if item != element {
			result = append(result, item)
		}
	}
	return result
}
