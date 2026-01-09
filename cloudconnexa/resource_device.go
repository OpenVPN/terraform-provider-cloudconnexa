package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// resourceDevice returns a Terraform resource schema for managing CloudConnexa devices.
// Note: Devices are created automatically when users connect. This resource allows
// managing existing devices (updating name, description).
func resourceDevice() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_device` to manage an existing CloudConnexa device. Devices are created automatically when users connect to the VPN. This resource allows you to update device properties like name and description.",
		CreateContext: resourceDeviceCreate,
		ReadContext:   resourceDeviceRead,
		UpdateContext: resourceDeviceUpdate,
		DeleteContext: resourceDeviceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The ID of the device to manage. Use the `cloudconnexa_devices` data source to find device IDs.",
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the device.",
				ValidateFunc: validation.StringLenBetween(1, 32),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The description of the device.",
				ValidateFunc: validation.StringLenBetween(0, 120),
			},
			"platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The platform of the device (e.g., Windows, macOS, iOS, Android).",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the device (ACTIVE, INACTIVE, BLOCKED, PENDING).",
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the user who owns the device.",
			},
		},
	}
}

// resourceDeviceCreate "creates" a device resource by adopting an existing device.
// Since devices are created automatically when users connect, this function
// simply reads and validates the device exists.
func resourceDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	deviceID := d.Get("device_id").(string)

	// Verify the device exists
	device, err := c.Devices.GetByID(deviceID)
	if err != nil {
		return diag.Errorf("Failed to get device with ID %s: %s", deviceID, err)
	}
	if device == nil {
		return diag.Errorf("Device with ID %s not found", deviceID)
	}

	d.SetId(deviceID)

	// Update the device with the provided name and description
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	updateRequest := cloudconnexa.DeviceUpdateRequest{
		Name:        name,
		Description: description,
	}

	_, err = c.Devices.Update(deviceID, updateRequest)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return resourceDeviceRead(ctx, d, m)
}

// resourceDeviceRead retrieves the current state of a device.
func resourceDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	id := d.Id()
	device, err := c.Devices.GetByID(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get device with ID %s: %s", id, err)...)
	}

	if device == nil {
		d.SetId("")
		return diags
	}

	d.Set("device_id", device.ID)
	d.Set("name", device.Name)
	d.Set("description", device.Description)
	d.Set("platform", device.Platform)
	d.Set("status", device.Status)
	d.Set("user_id", device.UserID)

	return diags
}

// resourceDeviceUpdate updates an existing device's name and/or description.
func resourceDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	if d.HasChanges("name", "description") {
		// Always send both fields to ensure we don't lose data.
		// The SDK uses omitempty, so we need to send both even if only one changed.
		updateRequest := cloudconnexa.DeviceUpdateRequest{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}

		_, err := c.Devices.Update(d.Id(), updateRequest)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	return resourceDeviceRead(ctx, d, m)
}

// resourceDeviceDelete "deletes" the device resource.
// Note: This does not actually delete the device from CloudConnexa,
// it only removes it from Terraform state. Devices are tied to user accounts
// and should be managed through user lifecycle.
func resourceDeviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// We don't actually delete the device - we just remove it from Terraform state
	// Devices are created when users connect and should be managed through user lifecycle
	d.SetId("")
	return nil
}
