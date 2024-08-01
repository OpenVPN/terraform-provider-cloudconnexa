package cloudconnexa

import (
	"context"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	userName := d.Get("username").(string)
	user, err := c.Users.GetByUsername(userName)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if user == nil {
		return append(diags, diag.Errorf("User with name %s was not found", userName)...)
	}

	d.SetId(user.Id)
	d.Set("username", user.Username)
	d.Set("role", user.Role)
	d.Set("email", user.Email)
	d.Set("auth_type", user.AuthType)
	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("group_id", user.GroupId)
	d.Set("status", user.Status)
	d.Set("devices", getUserDevicesSlice(&user.Devices))
	d.Set("connection_status", user.ConnectionStatus)
	return diags
}

func getUserDevicesSlice(userDevices *[]cloudconnexa.Device) []interface{} {
	devices := make([]interface{}, len(*userDevices))
	for i, d := range *userDevices {
		device := make(map[string]interface{})
		device["id"] = d.Id
		device["name"] = d.Name
		device["description"] = d.Description
		device["ip_v4_address"] = d.IPv4Address
		device["ip_v6_address"] = d.IPv6Address
		devices[i] = device
	}
	return devices
}
