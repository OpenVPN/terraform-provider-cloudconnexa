package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceDevices returns a Terraform data source resource for CloudConnexa devices.
// This resource allows users to read information about devices with optional filtering by user.
func dataSourceDevices() *schema.Resource {
	return &schema.Resource{
		Description: "Use `cloudconnexa_devices` data source to retrieve device information.",
		ReadContext: dataSourceDevicesRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter devices by user ID.",
			},
			"devices": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of devices.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device description.",
						},
						"platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device platform (e.g., Windows, macOS, iOS, Android).",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The device status (ACTIVE, INACTIVE, BLOCKED, PENDING).",
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user ID associated with this device.",
						},
					},
				},
			},
		},
	}
}

// dataSourceDevicesRead handles the read operation for the devices data source.
func dataSourceDevicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*providerMeta)
	c := meta.Client
	var diags diag.Diagnostics

	var devices []cloudconnexa.DeviceDetail
	var err error

	if v, ok := d.GetOk("user_id"); ok {
		devices, err = withRetry(ctx, meta.RetryConfig, func() ([]cloudconnexa.DeviceDetail, error) {
			return c.Devices.ListByUserID(v.(string))
		})
	} else {
		devices, err = withRetry(ctx, meta.RetryConfig, func() ([]cloudconnexa.DeviceDetail, error) {
			return c.Devices.ListAll()
		})
	}

	if err != nil {
		return diag.Errorf("Failed to get devices: %s", err)
	}

	d.SetId("devices")
	d.Set("devices", flattenDevices(devices))

	return diags
}

// flattenDevices converts a slice of CloudConnexa devices into a slice of interface{}
func flattenDevices(devices []cloudconnexa.DeviceDetail) []interface{} {
	result := make([]interface{}, len(devices))
	for i, dev := range devices {
		device := map[string]interface{}{
			"id":          dev.ID,
			"name":        dev.Name,
			"description": dev.Description,
			"platform":    dev.Platform,
			"user_id":     dev.UserID,
		}
		result[i] = device
	}
	return result
}

// dataSourceDevice returns a Terraform data source resource for a single CloudConnexa device.
func dataSourceDevice() *schema.Resource {
	return &schema.Resource{
		Description: "Use `cloudconnexa_device` data source to retrieve a specific device by ID.",
		ReadContext: dataSourceDeviceRead,
		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The device ID.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user ID that owns the device. Required by the CloudConnexa API.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device description.",
			},
			"platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device platform.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device status.",
			},
		},
	}
}

// dataSourceDeviceRead handles the read operation for a single device data source.
func dataSourceDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*providerMeta)
	c := meta.Client
	var diags diag.Diagnostics

	deviceID := d.Get("device_id").(string)
	userID := d.Get("user_id").(string)
	device, err := withRetry(ctx, meta.RetryConfig, func() (*cloudconnexa.DeviceDetail, error) {
		return c.Devices.GetByID(userID, deviceID)
	})
	if err != nil {
		return diag.Errorf("Failed to get device with ID %s: %s", deviceID, err)
	}

	d.SetId(device.ID)
	d.Set("name", device.Name)
	d.Set("description", device.Description)
	d.Set("platform", device.Platform)
	d.Set("user_id", device.UserID)

	return diags
}
