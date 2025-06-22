package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// validSource contains the list of valid source types for access groups
var validSource = []string{"USER_GROUP", "NETWORK", "HOST"}

// validDestination contains the list of valid destination types for access groups
var validDestination = []string{"USER_GROUP", "NETWORK", "HOST", "PUBLISHED_APPLICATION"}

// sourceRequestConversions maps source types to their API request format
var sourceRequestConversions = map[string]string{
	"NETWORK": "NETWORK_SERVICE",
}

// destinationRequestConversions maps destination types to their API request format
var destinationRequestConversions = map[string]string{
	"NETWORK":               "NETWORK_SERVICE",
	"HOST":                  "HOST_SERVICE",
	"PUBLISHED_APPLICATION": "PUBLISHED_SERVICE",
}

// sourceResponseConversions maps API response source types back to their original format
var sourceResponseConversions = map[string]string{
	"NETWORK_SERVICE": "NETWORK",
}

// destinationResponseConversions maps API response destination types back to their original format
var destinationResponseConversions = map[string]string{
	"NETWORK_SERVICE":   "NETWORK",
	"HOST_SERVICE":      "HOST",
	"PUBLISHED_SERVICE": "PUBLISHED_APPLICATION",
}

// resourceAccessGroup returns a Terraform resource schema for CloudConnexa access groups
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
				Optional:    true,
				Description: "The Access group description.",
			},
			"source": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				Elem:     resourceSource(),
			},
			"destination": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				Elem:     resourceDestination(),
			},
		},
	}
}

// resourceSource returns a schema for access group source configuration
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
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "ID of child entities assigned to access group source.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// resourceDestination returns a schema for access group destination configuration
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
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "ID of child entities assigned to access group destination.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// resourceAccessGroupCreate creates a new access group in CloudConnexa
func resourceAccessGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	request := resourceDataToAccessGroup(d)
	accessGroup, err := c.AccessGroups.Create(request)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(accessGroup.ID)
	return diags
}

// resourceAccessGroupRead retrieves an access group from CloudConnexa
func resourceAccessGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := d.Id()
	ag, err := c.AccessGroups.Get(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get access group with ID: %s, %s", id, err)...)
	}
	if ag == nil {
		d.SetId("")
	} else {
		setAccessGroupData(d, ag)
	}
	return diags
}

// resourceAccessGroupUpdate updates an existing access group in CloudConnexa
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

// resourceAccessGroupDelete removes an access group from CloudConnexa
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

// setAccessGroupData sets the Terraform resource data from a CloudConnexa access group
func setAccessGroupData(d *schema.ResourceData, ag *cloudconnexa.AccessGroup) {
	d.SetId(ag.ID)
	d.Set("name", ag.Name)
	d.Set("description", ag.Description)

	var sources []interface{}
	for _, source := range ag.Source {
		var parent = ""
		if source.Parent != "" {
			parent = source.Parent
		}
		var children []interface{}
		if source.Type == "USER_GROUP" || !source.AllCovered {
			for _, child := range source.Children {
				children = append(children, child)
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
		if destination.Parent != "" {
			parent = destination.Parent
		}
		var children []interface{}
		if destination.Type == "USER_GROUP" || !destination.AllCovered {
			for _, child := range destination.Children {
				children = append(children, child)
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

// resourceDataToAccessGroup converts Terraform resource data to a CloudConnexa access group
func resourceDataToAccessGroup(data *schema.ResourceData) *cloudconnexa.AccessGroup {
	name := data.Get("name").(string)
	description := data.Get("description").(string)

	request := &cloudconnexa.AccessGroup{
		Name:        name,
		Description: description,
	}

	sources := data.Get("source").(*schema.Set).List()

	for _, source := range sources {
		var convertedSource = source.(map[string]interface{})
		newSource := cloudconnexa.AccessItem{
			Type:       convert(convertedSource["type"].(string), sourceRequestConversions),
			AllCovered: convertedSource["all_covered"].(bool),
			Parent:     convertedSource["parent"].(string),
		}
		for _, child := range convertedSource["children"].(*schema.Set).List() {
			newSource.Children = append(newSource.Children, child.(string))
		}
		request.Source = append(request.Source, newSource)
	}

	destinations := data.Get("destination").(*schema.Set).List()

	for _, destination := range destinations {
		var mappedDestination = destination.(map[string]interface{})
		newDestination := cloudconnexa.AccessItem{
			Type:       convert(mappedDestination["type"].(string), destinationRequestConversions),
			AllCovered: mappedDestination["all_covered"].(bool),
			Parent:     mappedDestination["parent"].(string),
		}
		for _, child := range mappedDestination["children"].(*schema.Set).List() {
			newDestination.Children = append(newDestination.Children, child.(string))
		}
		request.Destination = append(request.Destination, newDestination)
	}

	return request
}

// convert transforms an input string using a provided conversion map.
// If the input string exists as a key in the conversions map, it returns the mapped value.
// Otherwise, it returns the original input string unchanged.
//
// Parameters:
//   - input: The string to be converted
//   - conversions: A map containing string-to-string conversions
//
// Returns:
//   - string: The converted string or the original input if no conversion exists
func convert(input string, conversions map[string]string) string {
	if output, exists := conversions[input]; exists {
		return output
	}
	return input
}
