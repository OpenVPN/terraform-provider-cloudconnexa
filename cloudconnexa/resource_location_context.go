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
		Description:   "Use `cloudconnexa_location_context` to create a Location Context Check.",
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
			"ip_check": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				AtLeastOneOf: []string{"ip_check", "country_check"},
				Elem:         ipCheckConfig(),
			},
			"country_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     countryCheckConfig(),
			},
			"default_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem:     defaultCheckConfig(),
			},
		},
	}
}

func ipCheckConfig() *schema.Resource {
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

func countryCheckConfig() *schema.Resource {
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

func defaultCheckConfig() *schema.Resource {
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
	d.SetId(response.ID)
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
	d.SetId(lc.ID)
	d.Set("name", lc.Name)
	d.Set("description", lc.Description)
	d.Set("user_groups_ids", lc.UserGroupsIDs)

	if lc.IPCheck != nil {
		ipCheck := make(map[string]interface{})
		ipCheck["allowed"] = lc.IPCheck.Allowed
		var ips []interface{}
		for _, ip := range lc.IPCheck.Ips {
			ips = append(ips, map[string]interface{}{
				"ip":          ip.IP,
				"description": ip.Description,
			})
		}
		ipCheck["ips"] = ips
		d.Set("ip_check", []interface{}{ipCheck})
	}

	if lc.CountryCheck != nil {
		countryCheck := make(map[string]interface{})
		countryCheck["allowed"] = lc.CountryCheck.Allowed
		countryCheck["countries"] = lc.CountryCheck.Countries
		d.Set("country_check", []interface{}{countryCheck})
	}

	defaultCheck := make(map[string]interface{})
	defaultCheck["allowed"] = lc.DefaultCheck.Allowed
	d.Set("default_check", []interface{}{defaultCheck})
}

func resourceDataToLocationContext(data *schema.ResourceData) *cloudconnexa.LocationContext {
	defaultCheckData := data.Get("default_check").([]interface{})[0].(map[string]interface{})
	defaultCheck := &cloudconnexa.DefaultCheck{
		Allowed: defaultCheckData["allowed"].(bool),
	}

	response := &cloudconnexa.LocationContext{
		ID:           data.Id(),
		Name:         data.Get("name").(string),
		Description:  data.Get("description").(string),
		DefaultCheck: defaultCheck,
	}

	for _, id := range data.Get("user_groups_ids").([]interface{}) {
		response.UserGroupsIDs = append(response.UserGroupsIDs, id.(string))
	}

	ipCheckList := data.Get("ip_check").([]interface{})
	if len(ipCheckList) > 0 {
		ipCheck := &cloudconnexa.IPCheck{}
		ipCheckData := ipCheckList[0].(map[string]interface{})
		ipCheck.Allowed = ipCheckData["allowed"].(bool)
		for _, ip := range ipCheckData["ips"].([]interface{}) {
			ipCheck.Ips = append(ipCheck.Ips, cloudconnexa.IP{
				IP:          ip.(map[string]interface{})["ip"].(string),
				Description: ip.(map[string]interface{})["description"].(string),
			})
		}
		response.IPCheck = ipCheck
	}

	countryCheckList := data.Get("country_check").([]interface{})
	if len(countryCheckList) > 0 && countryCheckList[0] != nil {
		countryCheckData := data.Get("country_check").([]interface{})[0].(map[string]interface{})
		countryCheck := &cloudconnexa.CountryCheck{
			Allowed: countryCheckData["allowed"].(bool),
		}
		for _, country := range countryCheckData["countries"].([]interface{}) {
			countryCheck.Countries = append(countryCheck.Countries, country.(string))
		}
		response.CountryCheck = countryCheck
	}

	return response
}
