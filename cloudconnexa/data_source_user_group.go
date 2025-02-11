package cloudconnexa

import (
	"context"
	"strings"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use an `cloudconnexa_user_group` data source to read an CloudConnexa user group.",
		ReadContext: dataSourceUserGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user group ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user group name.",
			},
			"vpn_region_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of region IDs this user group is associated with.",
			},
			"all_regions_included": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true all regions will be available for this user group.",
			},
			"internet_access": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of internet access provided. Valid values are `SPLIT_TUNNEL_ON`, `SPLIT_TUNNEL_OFF`, or `RESTRICTED_INTERNET`. Defaults to `SPLIT_TUNNEL_ON`.",
			},
			"max_device": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of devices per user.",
			},
			"system_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IPV4 and IPV6 addresses of the subnets associated with this user group.",
			},
			"connect_auth": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of connection authentication. Valid values are `NO_AUTH`, `ON_PRIOR_AUTH`, or `EVERY_TIME`.",
			},
		},
	}
}

func dataSourceUserGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	var userGroup *cloudconnexa.UserGroup
	var err error
	id := d.Get("id").(string)
	userGroup, err = c.UserGroups.Get(id)
	if err != nil {
		if strings.Contains(err.Error(), "user group not found") {
			return append(diags, diag.Errorf("User group with id %s was not found", id)...)
		} else {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if userGroup == nil {
		return append(diags, diag.Errorf("User group with id %s was not found", id)...)
	}
	d.SetId(userGroup.ID)
	d.Set("name", userGroup.Name)
	d.Set("vpn_region_ids", userGroup.VpnRegionIds)
	d.Set("all_regions_included", userGroup.AllRegionsIncluded)
	d.Set("internet_access", userGroup.InternetAccess)
	d.Set("max_device", userGroup.MaxDevice)
	d.Set("system_subnets", userGroup.SystemSubnets)
	d.Set("connect_auth", userGroup.ConnectAuth)
	return diags
}
