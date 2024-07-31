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

var (
	validValues = []string{"ANY", "BGP", "DHCP", "DNS", "FTP", "HTTP", "HTTPS", "IMAP", "IMAPS", "NTP", "POP3", "POP3S", "SMTP", "SMTPS", "SNMP", "SSH", "TELNET", "TFTP"}
)

func resourceIPService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPServiceCreate,
		ReadContext:   resourceServiceRead,
		DeleteContext: resourceServiceDelete,
		UpdateContext: resourceServiceUpdate,
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
				Elem:     resourceServiceConfig(),
			},
			"network_item_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"NETWORK", "HOST"}, false),
			},
			"network_item_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServiceUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)

	s, err := c.IPServices.Update(data.Id(), resourceDataToService(data))
	if err != nil {
		return diag.FromErr(err)
	}
	setResourceData(data, s)
	return nil
}

func resourceServiceConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"custom_service_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: customServiceTypesConfig(),
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

func customServiceTypesConfig() map[string]*schema.Schema {
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

func resourceServiceRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	service, err := c.IPServices.Get(data.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if service == nil {
		data.SetId("")
		return diags
	}
	setResourceData(data, service)
	return diags
}

func setResourceData(data *schema.ResourceData, service *cloudconnexa.IPServiceResponse) {
	data.SetId(service.Id)
	_ = data.Set("name", service.Name)
	_ = data.Set("description", service.Description)
	_ = data.Set("type", service.Type)
	_ = data.Set("routes", flattenRoutes(service.Routes))
	_ = data.Set("config", flattenServiceConfig(service.Config))
	_ = data.Set("network_item_type", service.NetworkItemType)
	_ = data.Set("network_item_id", service.NetworkItemId)
}

func resourceServiceDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	err := c.IPServices.Delete(data.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func flattenServiceConfig(config *cloudconnexa.IPServiceConfig) interface{} {
	var data = map[string]interface{}{
		"custom_service_types": flattenCustomServiceTypes(config.CustomServiceTypes),
		"service_types":        removeElement(config.ServiceTypes, "CUSTOM"),
	}
	return []interface{}{data}
}

func flattenCustomServiceTypes(types []*cloudconnexa.CustomIPServiceType) interface{} {
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

func flattenRoutes(routes []*cloudconnexa.Route) []string {
	var data []string
	for _, route := range routes {
		data = append(
			data,
			route.Subnet,
		)
	}
	return data
}

func resourceIPServiceCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*cloudconnexa.Client)

	service := resourceDataToService(data)
	createdService, err := client.IPServices.Create(service)
	if err != nil {
		return diag.FromErr(err)
	}
	setResourceData(data, createdService)
	return nil
}

func resourceDataToService(data *schema.ResourceData) *cloudconnexa.IPService {
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

	config := cloudconnexa.IPServiceConfig{}
	configList := data.Get("config").([]interface{})
	if len(configList) > 0 && configList[0] != nil {

		config.CustomServiceTypes = []*cloudconnexa.CustomIPServiceType{}
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

		for _, r := range mainConfig["service_types"].([]interface{}) {
			config.ServiceTypes = append(config.ServiceTypes, r.(string))
		}
		if len(config.CustomServiceTypes) > 0 && !slices.Contains(config.ServiceTypes, "CUSTOM") {
			config.ServiceTypes = append(config.ServiceTypes, "CUSTOM")
		}
	}

	s := &cloudconnexa.IPService{
		Name:            data.Get("name").(string),
		Description:     data.Get("description").(string),
		NetworkItemId:   data.Get("network_item_id").(string),
		NetworkItemType: data.Get("network_item_type").(string),
		Type:            data.Get("type").(string),
		Routes:          configRoutes,
		Config:          &config,
	}
	return s
}
