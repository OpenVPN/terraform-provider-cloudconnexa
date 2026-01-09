package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceHostRoute returns a Terraform resource for CloudConnexa host routes.
// This resource allows users to create, read, update, and delete routes on a CloudConnexa host.
func resourceHostRoute() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_host_route` to create a route on a CloudConnexa host.",
		CreateContext: resourceHostRouteCreate,
		UpdateContext: resourceHostRouteUpdate,
		ReadContext:   resourceHostRouteRead,
		DeleteContext: resourceHostRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the host on which to create the route.",
			},
			"subnet": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "The subnet CIDR for the route.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "The description for the route. Defaults to `Managed by Terraform`.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of route.",
			},
		},
	}
}

// resourceHostRouteCreate handles the creation of a new CloudConnexa host route.
func resourceHostRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	hostID := d.Get("host_id").(string)
	route := cloudconnexa.HostRoute{
		Subnet:      d.Get("subnet").(string),
		Description: d.Get("description").(string),
	}

	createdRoute, err := c.HostRoutes.Create(hostID, route)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdRoute.ID)
	d.Set("type", createdRoute.Type)
	d.Set("subnet", createdRoute.Subnet)

	return diags
}

// resourceHostRouteRead handles the read operation for a CloudConnexa host route.
func resourceHostRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	routeID := d.Id()
	route, err := c.HostRoutes.GetByID(routeID)
	if err != nil {
		return diag.Errorf("Failed to get host route with ID: %s, %s", routeID, err)
	}

	if route == nil {
		d.SetId("")
		return diags
	}

	d.Set("subnet", route.Subnet)
	d.Set("description", route.Description)
	d.Set("type", route.Type)
	d.Set("host_id", route.NetworkItemID)

	return diags
}

// resourceHostRouteUpdate handles the update operation for a CloudConnexa host route.
func resourceHostRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	if !d.HasChanges("description", "subnet") {
		return diags
	}

	route := cloudconnexa.HostRoute{
		ID:          d.Id(),
		Subnet:      d.Get("subnet").(string),
		Description: d.Get("description").(string),
	}

	err := c.HostRoutes.Update(route)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceHostRouteRead(ctx, d, m)
}

// resourceHostRouteDelete handles the deletion of a CloudConnexa host route.
func resourceHostRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	err := c.HostRoutes.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
