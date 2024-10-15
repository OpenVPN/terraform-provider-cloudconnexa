package cloudconnexa

import (
	"context"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceLocationContext() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_location_context` to create a Location Context policy.",
		CreateContext: resourceLocationContextCreate,
		ReadContext:   resourceLocationContextRead,
		DeleteContext: resourceLocationContextDelete,
		UpdateContext: resourceLocationContextUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Location Context name.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "The description for the UI. Defaults to `Managed by Terraform`.",
			},
			"user_groups_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of User Group IDs assigned to this policy.",
			},
			"ip_policy": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				AtLeastOneOf: []string{"ip_policy", "country_policy"},
				Elem:         ipPolicyConfig(),
			},
			"country_policy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     countryPolicyConfig(),
			},
			"default_policy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem:     defaultPolicyConfig(),
			},
		},
	}
}

func ipPolicyConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allowed": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"ips": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     ipConfig(),
			},
		},
	}
}

func ipConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func countryPolicyConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allowed": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"countries": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func defaultPolicyConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allowed": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceLocationContextCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	dr := resourceDataToLocationContext(d)
	response, err := c.LocationContexts.Create(dr)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(response.Id)
	return diags
}

func resourceLocationContextRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	recordId := d.Id()
	lc, err := c.LocationContexts.Get(recordId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if lc == nil {
		d.SetId("")
	} else {
		setLocationContextData(d, lc)
	}
	return diags
}

func resourceLocationContextUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	lc := resourceDataToLocationContext(d)
	_, err := c.LocationContexts.Update(d.Id(), lc)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func resourceLocationContextDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	routeId := d.Id()
	err := c.LocationContexts.Delete(routeId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func setLocationContextData(d *schema.ResourceData, lc *cloudconnexa.LocationContext) {
	d.SetId(lc.Id)
	d.Set("name", lc.Name)
	d.Set("description", lc.Description)
	d.Set("user_groups_ids", lc.UserGroupsIds)

	if lc.IpPolicy != nil {
		ipPolicy := make(map[string]interface{})
		ipPolicy["allowed"] = lc.IpPolicy.Allowed
		var ips []interface{}
		for _, ip := range lc.IpPolicy.Ips {
			ips = append(ips, map[string]interface{}{
				"ip":          ip.Ip,
				"description": ip.Description,
			})
		}
		ipPolicy["ips"] = ips
		d.Set("ip_policy", []interface{}{ipPolicy})
	}

	if lc.CountryPolicy != nil {
		countryPolicy := make(map[string]interface{})
		countryPolicy["allowed"] = lc.CountryPolicy.Allowed
		countryPolicy["countries"] = lc.CountryPolicy.Countries
		d.Set("country_policy", []interface{}{countryPolicy})
	}

	defaultPolicy := make(map[string]interface{})
	defaultPolicy["allowed"] = lc.DefaultPolicy.Allowed
	d.Set("default_policy", []interface{}{defaultPolicy})
}

func resourceDataToLocationContext(data *schema.ResourceData) *cloudconnexa.LocationContext {
	defaultPolicyData := data.Get("default_policy").([]interface{})[0].(map[string]interface{})
	defaultPolicy := &cloudconnexa.DefaultPolicy{
		Allowed: defaultPolicyData["allowed"].(bool),
	}

	response := &cloudconnexa.LocationContext{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		Description:   data.Get("description").(string),
		DefaultPolicy: defaultPolicy,
	}

	for _, id := range data.Get("user_groups_ids").([]interface{}) {
		response.UserGroupsIds = append(response.UserGroupsIds, id.(string))
	}

	ipPolicyList := data.Get("ip_policy").([]interface{})
	if len(ipPolicyList) > 0 {
		ipPolicy := &cloudconnexa.IpPolicy{}
		ipPolicyData := ipPolicyList[0].(map[string]interface{})
		ipPolicy.Allowed = ipPolicyData["allowed"].(bool)
		for _, ip := range ipPolicyData["ips"].([]interface{}) {
			ipPolicy.Ips = append(ipPolicy.Ips, cloudconnexa.Ip{
				Ip:          ip.(map[string]interface{})["ip"].(string),
				Description: ip.(map[string]interface{})["description"].(string),
			})
		}
		response.IpPolicy = ipPolicy
	}

	countryPolicyList := data.Get("country_policy").([]interface{})
	if len(countryPolicyList) > 0 && countryPolicyList[0] != nil {
		countryPolicyData := data.Get("country_policy").([]interface{})[0].(map[string]interface{})
		countryPolicy := &cloudconnexa.CountryPolicy{
			Allowed: countryPolicyData["allowed"].(bool),
		}
		for _, country := range countryPolicyData["countries"].([]interface{}) {
			countryPolicy.Countries = append(countryPolicy.Countries, country.(string))
		}
		response.CountryPolicy = countryPolicy
	}

	return response
}
