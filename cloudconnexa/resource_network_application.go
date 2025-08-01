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

// resourceNetworkApplication returns a Terraform resource schema for managing network applications
func resourceNetworkApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkApplicationCreate,
		ReadContext:   resourceNetworkApplicationRead,
		DeleteContext: resourceNetworkApplicationDelete,
		UpdateContext: resourceNetworkApplicationUpdate,
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
				Elem:     resourceNetworkApplicationRoute(),
			},
			"config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     resourceNetworkApplicationConfig(),
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// resourceNetworkApplicationUpdate handles updates to an existing network application
func resourceNetworkApplicationUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)

	s, err := c.NetworkApplications.Update(data.Id(), resourceDataToNetworkApplication(data))
	if err != nil {
		return diag.FromErr(err)
	}
	setNetworkApplicationData(data, s)
	return nil
}

// resourceNetworkApplicationRoute returns the schema for network application routes
func resourceNetworkApplicationRoute() *schema.Resource {
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

// resourceNetworkApplicationConfig returns the schema for network application configuration
func resourceNetworkApplicationConfig() *schema.Resource {
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

// resourceNetworkApplicationRead retrieves and sets the state of a network application
func resourceNetworkApplicationRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := data.Id()
	application, err := c.NetworkApplications.Get(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get network application with ID: %s, %s", id, err)...)
	}
	if application == nil {
		data.SetId("")
		return diags
	}
	setNetworkApplicationData(data, application)
	return diags
}

// resourceNetworkApplicationDelete handles the deletion of a network application
func resourceNetworkApplicationDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	err := c.NetworkApplications.Delete(data.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

// resourceNetworkApplicationCreate handles the creation of a new network application
func resourceNetworkApplicationCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*cloudconnexa.Client)

	application := resourceDataToNetworkApplication(data)
	createdApplication, err := client.NetworkApplications.Create(application)
	if err != nil {
		return diag.FromErr(err)
	}
	setNetworkApplicationData(data, createdApplication)
	return nil
}

// resourceDataToNetworkApplication converts Terraform resource data to a CloudConnexa Application
func resourceDataToNetworkApplication(data *schema.ResourceData) *cloudconnexa.Application {
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
		NetworkItemID:   data.Get("network_id").(string),
		NetworkItemType: "NETWORK",
		Routes:          configRoutes,
		Config:          &config,
	}
	return s
}

// setNetworkApplicationData sets the Terraform state data from a network application response
func setNetworkApplicationData(data *schema.ResourceData, application *cloudconnexa.ApplicationResponse) {
	data.SetId(application.ID)
	_ = data.Set("name", application.Name)
	_ = data.Set("description", application.Description)
	_ = data.Set("routes", flattenApplicationRoutes(application.Routes))
	_ = data.Set("config", flattenApplicationConfig(application.Config))
	_ = data.Set("network_id", application.NetworkItemID)
}
