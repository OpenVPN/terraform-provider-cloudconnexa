package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// resourceUser returns a Terraform resource schema for managing CloudConnexa users
func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_user` to create an CloudConnexa user.",
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "A username for the user.",
			},
			"email": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "An invitation to CloudConnexa account will be sent to this email. It will include an initial password and a VPN setup guide.",
			},
			"first_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20),
				Description:  "User's first name.",
			},
			"last_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20),
				Description:  "User's last name.",
			},
			"group_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The UUID of a user's group.",
				ValidateFunc: validation.IsUUID,
			},
			"secondary_groups_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The UUIDs of secondary user's groups.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "MEMBER",
				Description: "The type of user role. Valid values are `ADMIN`, `MEMBER`, or `OWNER`.",
			},
			"devices": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "When a user signs in, the device that they use will be added to their account. You can read more at [CloudConnexa Device](https://openvpn.net/cloud-docs/device/).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 32),
							Description:  "A device name.",
						},
						"description": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 120),
							Description:  "A device description.",
						},
						"ipv4_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An IPv4 address of the device.",
						},
						"ipv6_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An IPv6 address of the device.",
						},
					},
				},
			},
		},
	}
}

// resourceUserCreate creates a new CloudConnexa user
func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	username := d.Get("username").(string)
	email := d.Get("email").(string)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	groupId := d.Get("group_id").(string)
	secondaryGroupsIds := d.Get("secondary_groups_ids").([]interface{})
	role := d.Get("role").(string)
	configDevices := d.Get("devices").([]interface{})
	var devices []cloudconnexa.Device
	for _, d := range configDevices {
		device := d.(map[string]interface{})
		devices = append(
			devices,
			cloudconnexa.Device{
				Name:        device["name"].(string),
				Description: device["description"].(string),
				IPv4Address: device["ipv4_address"].(string),
				IPv6Address: device["ipv6_address"].(string),
			},
		)

	}
	u := cloudconnexa.User{
		Username:          username,
		Email:             email,
		FirstName:         firstName,
		LastName:          lastName,
		GroupID:           groupId,
		SecondaryGroupIds: toStrings(secondaryGroupsIds),
		Devices:           devices,
		Role:              role,
	}
	user, err := c.Users.Create(u)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(user.ID)
	return diags
}

// resourceUserRead retrieves information about an existing CloudConnexa user
func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := d.Id()
	u, err := c.Users.Get(id)

	if err != nil {
		return append(diags, diag.Errorf("Failed to get user with ID: %s, %s", id, err)...)
	}
	if u == nil {
		d.SetId("")
	} else {
		d.Set("username", u.Username)
		d.Set("email", u.Email)
		d.Set("first_name", u.FirstName)
		d.Set("last_name", u.LastName)
		d.Set("group_id", u.GroupID)
		d.Set("secondary_groups_ids", u.SecondaryGroupIds)
		d.Set("devices", u.Devices)
		d.Set("role", u.Role)
	}
	return diags
}

// resourceUserUpdate updates an existing CloudConnexa user's information
func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	if !d.HasChanges("first_name", "last_name", "group_id", "email", "role", "secondary_groups_ids") {
		return diags
	}

	u, err := c.Users.Get(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	_, email := d.GetChange("email")
	_, firstName := d.GetChange("first_name")
	_, lastName := d.GetChange("last_name")
	_, role := d.GetChange("role")
	_, groupId := d.GetChange("group_id")
	_, secondaryGroupsIds := d.GetChange("secondary_groups_ids")
	status := u.Status

	err = c.Users.Update(cloudconnexa.User{
		ID:                d.Id(),
		Email:             email.(string),
		FirstName:         firstName.(string),
		LastName:          lastName.(string),
		GroupID:           groupId.(string),
		SecondaryGroupIds: toStrings(secondaryGroupsIds.([]interface{})),
		Role:              role.(string),
		Status:            status,
	})

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// resourceUserDelete removes an existing CloudConnexa user
func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	userId := d.Id()
	err := c.Users.Delete(userId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
