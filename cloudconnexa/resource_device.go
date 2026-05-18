package cloudconnexa

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// resourceDevice returns a Terraform resource schema for managing CloudConnexa
// devices as a standalone child of a user. Devices created via this resource
// are owned by Terraform: they are provisioned in Create and removed in Delete.
func resourceDevice() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_device` to manage a CloudConnexa device attached to a user. The resource creates the device in CloudConnexa and removes it on destroy.",
		CreateContext: resourceDeviceCreate,
		ReadContext:   resourceDeviceRead,
		UpdateContext: resourceDeviceUpdate,
		DeleteContext: resourceDeviceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDeviceImport,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the user that owns the device. Changing this recreates the device.",
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
			"client_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The client UUID of the device. Set by the OpenVPN client when the device first connects; you can also pre-assign it here.",
			},
			"device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CloudConnexa-assigned device ID.",
			},
			"platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The platform of the device (e.g., Windows, macOS, iOS, Android).",
			},
			"ipv4_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv4 address assigned to the device.",
			},
			"ipv6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv6 address assigned to the device.",
			},
			"connection_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current connection status of the device.",
			},
		},
	}
}

// resourceDeviceCreate provisions a new device for the given user.
func resourceDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*providerMeta)
	c := meta.Client

	userID := d.Get("user_id").(string)
	req := cloudconnexa.DeviceCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ClientUUID:  d.Get("client_uuid").(string),
	}

	device, err := withRetry(ctx, meta.RetryConfig, func() (*cloudconnexa.DeviceDetail, error) {
		return c.Devices.Create(userID, req)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(device.ID)
	return resourceDeviceRead(ctx, d, m)
}

// resourceDeviceRead refreshes Terraform state from the CloudConnexa API.
func resourceDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*providerMeta)
	c := meta.Client

	userID := d.Get("user_id").(string)
	device, err := withRetry(ctx, meta.RetryConfig, func() (*cloudconnexa.DeviceDetail, error) {
		return c.Devices.GetByID(userID, d.Id())
	})
	if err != nil {
		return diag.Errorf("Failed to get device with ID %s: %s", d.Id(), err)
	}

	d.Set("device_id", device.ID)
	d.Set("name", device.Name)
	d.Set("description", device.Description)
	d.Set("user_id", device.UserID)
	d.Set("client_uuid", device.ClientUUID)
	d.Set("platform", device.Platform)
	d.Set("ipv4_address", device.IPV4Address)
	d.Set("ipv6_address", device.IPV6Address)
	d.Set("connection_status", device.ConnectionStatus)

	return nil
}

// resourceDeviceUpdate pushes name/description changes back to CloudConnexa.
func resourceDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*providerMeta)
	c := meta.Client

	if d.HasChanges("name", "description") {
		// Both fields are always sent so omitempty on the SDK struct doesn't blank a value the user kept.
		req := cloudconnexa.DeviceUpdateRequest{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}
		userID := d.Get("user_id").(string)
		if _, err := withRetry(ctx, meta.RetryConfig, func() (*cloudconnexa.DeviceDetail, error) {
			return c.Devices.Update(userID, d.Id(), req)
		}); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDeviceRead(ctx, d, m)
}

// resourceDeviceDelete removes the device from CloudConnexa.
func resourceDeviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*providerMeta)
	c := meta.Client

	userID := d.Get("user_id").(string)
	if err := withRetryNoBody(ctx, meta.RetryConfig, func() error {
		return c.Devices.Delete(userID, d.Id())
	}); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

// resourceDeviceImport parses the import ID of the form "user_id/device_id" and
// populates both attributes so the subsequent Read can call the API.
func resourceDeviceImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("expected import ID in the form \"user_id/device_id\", got %q", d.Id())
	}
	d.Set("user_id", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
