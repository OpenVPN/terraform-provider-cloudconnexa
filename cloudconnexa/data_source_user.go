package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceUser returns a Terraform data source resource for CloudConnexa users.
// This resource allows users to read information about a specific CloudConnexa user by their username.
//
// Returns:
//   - *schema.Resource: A Terraform resource definition for user data source
func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "Use a `cloudconnexa_user` data source to read a specific CloudConnexa user.",
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of this resource.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username of the user.",
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of user role. Valid values are `ADMIN`, `MEMBER`, or `OWNER`.",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of the user.",
			},
			"auth_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The authentication type of the user.",
			},
			"first_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's first name.",
			},
			"last_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's last name.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's group id.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's status.",
			},
			"connection_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's connection status.",
			},
			"devices": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of user devices.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device's id.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device's name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device's description.",
						},
						"ip_v4_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device's IPV4 address.",
						},
						"ip_v6_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device's IPV6 address.",
						},
					},
				},
			},
		},
	}
}

// dataSourceUserRead handles the read operation for the user data source.
// It retrieves user information from CloudConnexa using the provided username
// and updates the Terraform state with the retrieved data.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The interface containing the CloudConnexa client
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during the operation
func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	userName := d.Get("username").(string)
	user, err := c.Users.GetByUsername(userName)
	if err != nil {
		return diag.Errorf("Failed to get user with username: %s, %s", userName, err)
	}
	if user == nil {
		return append(diags, diag.Errorf("User with name %s was not found", userName)...)
	}

	d.SetId(user.ID)
	d.Set("username", user.Username)
	d.Set("role", user.Role)
	d.Set("email", user.Email)
	d.Set("auth_type", user.AuthType)
	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("group_id", user.GroupID)
	d.Set("status", user.Status)
	d.Set("devices", getUserDevicesSlice(&user.Devices))
	d.Set("connection_status", user.ConnectionStatus)
	return diags
}

// getUserDevicesSlice converts a slice of CloudConnexa devices into a slice of interface{}
// that can be used by Terraform. It maps each device's properties to a map structure.
//
// Parameters:
//   - userDevices: A pointer to a slice of CloudConnexa devices
//
// Returns:
//   - []interface{}: A slice of maps containing device information
func getUserDevicesSlice(userDevices *[]cloudconnexa.Device) []interface{} {
	devices := make([]interface{}, len(*userDevices))
	for i, d := range *userDevices {
		device := make(map[string]interface{})
		device["id"] = d.ID
		device["name"] = d.Name
		device["description"] = d.Description
		device["ip_v4_address"] = d.IPv4Address
		device["ip_v6_address"] = d.IPv6Address
		devices[i] = device
	}
	return devices
}
