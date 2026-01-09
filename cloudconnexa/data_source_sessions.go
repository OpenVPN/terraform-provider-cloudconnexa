package cloudconnexa

import (
	"context"
	"time"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// dataSourceSessions returns a Terraform data source resource for CloudConnexa VPN sessions.
// This resource allows users to read information about VPN sessions with optional filtering.
func dataSourceSessions() *schema.Resource {
	return &schema.Resource{
		Description: "Use `cloudconnexa_sessions` data source to retrieve VPN session information.",
		ReadContext: dataSourceSessionsRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Filter sessions by status. Valid values are `ACTIVE`, `COMPLETED`, or `FAILED`.",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "COMPLETED", "FAILED"}, false),
			},
			"start_date": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Filter sessions starting from this date (RFC3339 format).",
				ValidateFunc: validation.IsRFC3339Time,
			},
			"end_date": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Filter sessions until this date (RFC3339 format).",
				ValidateFunc: validation.IsRFC3339Time,
			},
			"sessions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of VPN sessions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The session ID.",
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user ID associated with this session.",
						},
						"device_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device ID used in this session.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN region ID.",
						},
						"bytes_in": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of bytes received.",
						},
						"bytes_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of bytes sent.",
						},
						"connector_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connector name.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username.",
						},
						"device_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device name.",
						},
						"client_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The client IP address.",
						},
						"start_date_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The session start date and time.",
						},
						"vpn_ipv4": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN IPv4 address assigned to the client.",
						},
						"network_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network name.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region name.",
						},
						"connection_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connection status.",
						},
					},
				},
			},
		},
	}
}

// dataSourceSessionsRead handles the read operation for the sessions data source.
func dataSourceSessionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	options := cloudconnexa.SessionsListOptions{
		Size: 100,
	}

	if v, ok := d.GetOk("status"); ok {
		options.Status = cloudconnexa.SessionStatus(v.(string))
	}

	if v, ok := d.GetOk("start_date"); ok {
		t, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return diag.Errorf("Invalid start_date format: %s", err)
		}
		options.StartDate = &t
	}

	if v, ok := d.GetOk("end_date"); ok {
		t, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return diag.Errorf("Invalid end_date format: %s", err)
		}
		options.EndDate = &t
	}

	sessions, err := c.Sessions.ListAll(options)
	if err != nil {
		return diag.Errorf("Failed to get sessions: %s", err)
	}

	d.SetId("sessions")
	d.Set("sessions", flattenSessions(sessions))

	return diags
}

// flattenSessions converts a slice of CloudConnexa sessions into a slice of interface{}
func flattenSessions(sessions []cloudconnexa.Session) []interface{} {
	result := make([]interface{}, len(sessions))
	for i, s := range sessions {
		session := map[string]interface{}{
			"session_id":        s.SessionID,
			"user_id":           s.UserID,
			"device_id":         s.DeviceID,
			"region_id":         s.RegionID,
			"bytes_in":          s.BytesIn,
			"bytes_out":         s.BytesOut,
			"connector_name":    s.ConnectorName,
			"user_name":         s.UserName,
			"device_name":       s.DeviceName,
			"client_ip":         s.ClientIP,
			"start_date_time":   s.StartDateTime.Format(time.RFC3339),
			"vpn_ipv4":          s.VpnIPv4,
			"network_name":      s.NetworkName,
			"region_name":       s.RegionName,
			"connection_status": s.ConnectionStatus,
		}
		result[i] = session
	}
	return result
}
