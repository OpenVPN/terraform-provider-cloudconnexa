package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// resourceUserGroup returns a Terraform resource for managing CloudConnexa user groups.
// It defines the schema and CRUD operations for user group resources.
//
// Returns:
//   - *schema.Resource: A configured Terraform resource for user groups
func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_user_group` to create an CloudConnexa user group.",
		CreateContext: resourceUserGroupCreate,
		ReadContext:   resourceUserGroupRead,
		UpdateContext: resourceUserGroupUpdate,
		DeleteContext: resourceUserGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the user group.",
			},
			"connect_auth": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "AUTO",
				ValidateFunc: validation.StringInSlice([]string{"NO_AUTH", "ON_PRIOR_AUTH", "EVERY_TIME"}, false),
			},
			"internet_access": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "SPLIT_TUNNEL_ON",
				ValidateFunc: validation.StringInSlice([]string{"SPLIT_TUNNEL_ON", "SPLIT_TUNNEL_OFF", "RESTRICTED_INTERNET"}, false),
			},
			"max_device": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "The maximum number of devices that can be connected to the user group.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 40),
				Description:  "The name of the user group.",
			},
			"system_subnets": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Default:     nil,
				Description: "A list of subnets that are accessible to the user group.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vpn_region_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of regions IDs that are accessible to the user group. Actual list of available regions can be obtained from data_source_vpn_regions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"all_regions_included": {
				Type:         schema.TypeBool,
				Optional:     true,
				ExactlyOneOf: []string{"vpn_region_ids", "all_regions_included"},
				Description:  "If true all regions will be available for this user group.",
			},
		},
	}
}

// resourceUserGroupUpdate updates an existing CloudConnexa user group.
// It handles the update operation and updates the Terraform state accordingly.
//
// Parameters:
//   - ctx: The context for the operation
//   - data: The Terraform resource data
//   - i: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred
func resourceUserGroupUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	ug := resourceDataToUserGroup(data)

	userGroup, err := c.UserGroups.Update(data.Id(), ug)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if userGroup == nil {
		data.SetId("")
	} else {
		updateUserGroupData(data, userGroup)
	}

	return diags
}

// resourceDataToUserGroup converts Terraform resource data to a CloudConnexa UserGroup.
// It extracts all relevant fields from the Terraform state and creates a UserGroup object.
//
// Parameters:
//   - data: The Terraform resource data
//
// Returns:
//   - *cloudconnexa.UserGroup: A configured UserGroup object
func resourceDataToUserGroup(data *schema.ResourceData) *cloudconnexa.UserGroup {
	name := data.Get("name").(string)
	connectAuth := data.Get("connect_auth").(string)
	maxDevice := data.Get("max_device").(int)
	internetAccess := data.Get("internet_access").(string)
	configSystemSubnets := data.Get("system_subnets").([]interface{})
	var systemSubnets []string
	for _, s := range configSystemSubnets {
		systemSubnets = append(systemSubnets, s.(string))
	}
	configVpnRegionIds := data.Get("vpn_region_ids").([]interface{})
	var vpnRegionIds []string
	for _, r := range configVpnRegionIds {
		vpnRegionIds = append(vpnRegionIds, r.(string))
	}
	allRegionsIncluded := data.Get("all_regions_included").(bool)

	ug := &cloudconnexa.UserGroup{
		Name:               name,
		ConnectAuth:        connectAuth,
		MaxDevice:          maxDevice,
		SystemSubnets:      systemSubnets,
		VpnRegionIDs:       vpnRegionIds,
		InternetAccess:     internetAccess,
		AllRegionsIncluded: allRegionsIncluded,
	}
	return ug
}

// updateUserGroupData updates the Terraform resource data with values from a CloudConnexa UserGroup.
// It sets all relevant fields in the Terraform state based on the provided UserGroup.
//
// Parameters:
//   - data: The Terraform resource data to update
//   - userGroup: The CloudConnexa UserGroup containing the values to set
func updateUserGroupData(data *schema.ResourceData, userGroup *cloudconnexa.UserGroup) {
	data.SetId(userGroup.ID)
	_ = data.Set("connect_auth", userGroup.ConnectAuth)
	_ = data.Set("max_device", userGroup.MaxDevice)
	_ = data.Set("name", userGroup.Name)
	_ = data.Set("system_subnets", userGroup.SystemSubnets)
	_ = data.Set("internet_access", userGroup.InternetAccess)
	if !userGroup.AllRegionsIncluded {
		_ = data.Set("vpn_region_ids", userGroup.VpnRegionIDs)
	}
	_ = data.Set("all_regions_included", userGroup.AllRegionsIncluded)
}

// resourceUserGroupDelete deletes a CloudConnexa user group.
// It handles the deletion operation and updates the Terraform state accordingly.
//
// Parameters:
//   - ctx: The context for the operation
//   - data: The Terraform resource data
//   - i: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred
func resourceUserGroupDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	err := c.UserGroups.Delete(data.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	data.SetId("")
	return diags
}

// resourceUserGroupRead reads a CloudConnexa user group.
// It retrieves the current state of the user group and updates the Terraform state accordingly.
//
// Parameters:
//   - ctx: The context for the operation
//   - data: The Terraform resource data
//   - i: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred
func resourceUserGroupRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := data.Id()
	userGroup, err := c.UserGroups.Get(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get user group with ID: %s, %s", id, err)...)
	}

	if userGroup == nil {
		data.SetId("")
	} else {
		updateUserGroupData(data, userGroup)
	}
	return diags
}

// resourceUserGroupCreate creates a new CloudConnexa user group.
// It handles the creation operation and updates the Terraform state accordingly.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred
func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	ug := resourceDataToUserGroup(d)

	userGroup, err := c.UserGroups.Create(ug)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	updateUserGroupData(d, userGroup)
	return diags
}
