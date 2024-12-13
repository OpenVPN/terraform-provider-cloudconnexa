package cloudconnexa

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

var (
	validSource              = []string{"USER_GROUP", "NETWORK", "HOST"}
	validDestination         = []string{"USER_GROUP", "NETWORK", "HOST", "PUBLISHED_APPLICATION"}
	sourceRequestConversions = map[string]string{
		"NETWORK": "NETWORK_SERVICE",
	}
	destinationRequestConversions = map[string]string{
		"NETWORK":               "NETWORK_SERVICE",
		"HOST":                  "HOST_SERVICE",
		"PUBLISHED_APPLICATION": "PUBLISHED_SERVICE",
	}
	sourceResponseConversions = map[string]string{
		"NETWORK_SERVICE": "NETWORK",
	}
	destinationResponseConversions = map[string]string{
		"NETWORK_SERVICE":   "NETWORK",
		"HOST_SERVICE":      "HOST",
		"PUBLISHED_SERVICE": "PUBLISHED_APPLICATION",
	}
)

func resourceAccessGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_access_group` to create an Access group.",
		CreateContext: resourceAccessGroupCreate,
		ReadContext:   resourceAccessGroupRead,
		DeleteContext: resourceAccessGroupDelete,
		UpdateContext: resourceAccessGroupUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Access group name.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Access group description.",
			},
			"source": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				ForceNew: true,
				Elem:     resourceSource(),
			},
			"destination": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				ForceNew: true,
				Elem:     resourceDestination(),
			},
		},
	}
}

func resourceSource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Source type.",
				ValidateFunc: validation.StringInSlice(validSource, false),
			},
			"all_covered": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Allows to select all items of specific type or all children under specific parent",
			},
			"parent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the entity assigned to access group source.",
			},
			"children": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "ID of child entities assigned to access group source.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDestination() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Destination type.",
				ValidateFunc: validation.StringInSlice(validDestination, false),
			},
			"all_covered": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Allows to select all items of specific type or all children under specific parent",
			},
			"parent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the entity assigned to access group destination.",
			},
			"children": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "ID of child entities assigned to access group destination.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAccessGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	request := resourceDataToAccessGroup(d)
	accessGroup, err := c.AccessGroups.Create(request)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(accessGroup.Id)
	return diags
}

func resourceAccessGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := d.Id()
	ag, err := c.AccessGroups.Get(id)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if ag == nil {
		d.SetId("")
	} else {
		setAccessGroupData(d, ag)
	}
	return diags
}

func resourceAccessGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	ag := resourceDataToAccessGroup(d)
	savedAccessGroup, err := c.AccessGroups.Update(d.Id(), ag)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	setAccessGroupData(d, savedAccessGroup)
	return diags
}

func resourceAccessGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := d.Id()
	err := c.AccessGroups.Delete(id)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func setAccessGroupData(d *schema.ResourceData, ag *cloudconnexa.AccessGroupResponse) {
	d.SetId(ag.Id)
	d.Set("name", ag.Name)
	d.Set("description", ag.Description)

	var sources []interface{}
	for _, source := range ag.Source {
		var parent = ""
		if source.Parent != nil {
			parent = source.Parent.Id
		}
		var children []interface{}
		if source.Type == "USER_GROUP" || !source.AllCovered {
			for _, child := range source.Children {
				children = append(children, child.Id)
			}
		}

		sources = append(sources, map[string]interface{}{
			"type":        convert(source.Type, sourceResponseConversions),
			"all_covered": source.AllCovered,
			"parent":      parent,
			"children":    children,
		})
	}
	d.Set("source", sources)

	var destinations []interface{}
	for _, destination := range ag.Destination {
		var parent = ""
		if destination.Parent != nil {
			parent = destination.Parent.Id
		}
		var children []interface{}
		if destination.Type == "USER_GROUP" || !destination.AllCovered {
			for _, child := range destination.Children {
				children = append(children, child.Id)
			}
		}

		destinations = append(destinations, map[string]interface{}{
			"type":        convert(destination.Type, destinationResponseConversions),
			"all_covered": destination.AllCovered,
			"parent":      parent,
			"children":    children,
		})
	}
	d.Set("destination", destinations)
}

func resourceDataToAccessGroup(data *schema.ResourceData) *cloudconnexa.AccessGroupRequest {
	name := data.Get("name").(string)
	description := data.Get("description").(string)

	request := &cloudconnexa.AccessGroupRequest{
		Name:        name,
		Description: description,
	}

	sources := data.Get("source").(*schema.Set).List()

	for _, source := range sources {
		var convertedSource = source.(map[string]interface{})
		newSource := cloudconnexa.AccessItemRequest{
			Type:       convert(convertedSource["type"].(string), sourceRequestConversions),
			AllCovered: convertedSource["all_covered"].(bool),
			Parent:     convertedSource["parent"].(string),
		}
		for _, child := range convertedSource["children"].([]interface{}) {
			newSource.Children = append(newSource.Children, child.(string))
		}
		request.Source = append(request.Source, newSource)
	}

	destinations := data.Get("destination").(*schema.Set).List()

	for _, destination := range destinations {
		var mappedDestination = destination.(map[string]interface{})
		newDestination := cloudconnexa.AccessItemRequest{
			Type:       convert(mappedDestination["type"].(string), destinationRequestConversions),
			AllCovered: mappedDestination["all_covered"].(bool),
			Parent:     mappedDestination["parent"].(string),
		}
		for _, child := range mappedDestination["children"].([]interface{}) {
			newDestination.Children = append(newDestination.Children, child.(string))
		}
		request.Destination = append(request.Destination, newDestination)
	}

	return request
}

func convert(input string, conversions map[string]string) string {
	if output, exists := conversions[input]; exists {
		return output
	}
	return input
}
