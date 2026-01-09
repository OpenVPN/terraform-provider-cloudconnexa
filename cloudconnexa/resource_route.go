package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceRoute returns a Terraform resource for CloudConnexa routes.
// This resource allows users to create, read, update, and delete routes on a CloudConnexa network.
//
// Returns:
//   - *schema.Resource: A Terraform resource definition for routes
func resourceRoute() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_route` to create a route on an CloudConnexa network.",
		CreateContext: resourceRouteCreate,
		UpdateContext: resourceRouteUpdate,
		ReadContext:   resourceRouteRead,
		DeleteContext: resourceRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IP_V4", "IP_V6"}, false),
				Description:  "The type of route. Valid values are `IP_V4` and `IP_V6`.",
			},
			"subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target value of the default route.",
			},
			"network_item_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the network on which to create the route.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
		},
	}
}

// resourceRouteCreate handles the creation of a new CloudConnexa route.
// It creates a route with the specified type, subnet, and description on the given network.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func resourceRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	networkItemId := d.Get("network_item_id").(string)
	routeType := d.Get("type").(string)
	routeSubnet := d.Get("subnet").(string)
	descriptionValue := d.Get("description").(string)
	r := cloudconnexa.Route{
		Type:        routeType,
		Subnet:      routeSubnet,
		Description: descriptionValue,
	}
	route, err := c.Routes.Create(networkItemId, r)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(route.ID)
	if routeType == "IP_V4" || routeType == "IP_V6" {
		d.Set("subnet", route.Subnet)
	}
	return diags
}

// resourceRouteRead handles the read operation for a CloudConnexa route.
// It retrieves the route information and updates the Terraform state.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func resourceRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := d.Id()
	r, err := c.Routes.Get(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get route with ID: %s, %s", id, err)...)
	}
	if r == nil {
		d.SetId("")
	} else {
		d.Set("type", r.Type)
		if r.Type == "IP_V4" || r.Type == "IP_V6" {
			d.Set("subnet", r.Subnet)
		}
		d.Set("description", r.Description)
		d.Set("network_item_id", r.NetworkItemID)
	}
	return diags
}

// resourceRouteUpdate handles the update operation for a CloudConnexa route.
// It updates the route's description and subnet if they have changed.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func resourceRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	if !d.HasChanges("description", "subnet") {
		return diags
	}

	_, description := d.GetChange("description")
	_, subnet := d.GetChange("subnet")
	r := cloudconnexa.Route{
		ID:          d.Id(),
		Description: description.(string),
		Subnet:      subnet.(string),
	}

	err := c.Routes.Update(r)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

// resourceRouteDelete handles the deletion of a CloudConnexa route.
// It removes the route from the CloudConnexa network.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func resourceRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	routeId := d.Id()
	err := c.Routes.Delete(routeId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
