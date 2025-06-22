package cloudconnexa

import (
	"context"
	"fmt"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceDnsRecord returns a Terraform resource for managing DNS records in CloudConnexa.
// It defines the schema and CRUD operations for DNS record management.
func resourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_dns_record` to create a DNS record on your VPN.",
		CreateContext: resourceDnsRecordCreate,
		ReadContext:   resourceDnsRecordRead,
		DeleteContext: resourceDnsRecordDelete,
		UpdateContext: resourceDnsRecordUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: validateAtLeastOneNonEmptyList,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The DNS record name.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "The description for the UI. Defaults to `Managed by Terraform`.",
			},
			"ip_v4_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv4Address,
				},
				Description:  "The list of IPV4 addresses to which this record will resolve.",
				AtLeastOneOf: []string{"ip_v4_addresses", "ip_v6_addresses"},
			},
			"ip_v6_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv6Address,
				},
				Description: "The list of IPV6 addresses to which this record will resolve.",
			},
		},
	}
}

// resourceDnsRecordCreate creates a new DNS record in CloudConnexa.
// It converts the Terraform resource data into a DNS record and sends it to the API.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The provider meta interface
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred
func resourceDnsRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	domain := d.Get("domain").(string)
	description := d.Get("description").(string)
	ipV4Addresses := d.Get("ip_v4_addresses").([]interface{})
	ipV4AddressesSlice := make([]string, 0)
	for _, a := range ipV4Addresses {
		ipV4AddressesSlice = append(ipV4AddressesSlice, a.(string))
	}
	ipV6Addresses := d.Get("ip_v6_addresses").([]interface{})
	ipV6AddressesSlice := make([]string, 0)
	for _, a := range ipV6Addresses {
		ipV6AddressesSlice = append(ipV6AddressesSlice, a.(string))
	}
	dr := cloudconnexa.DNSRecord{
		Domain:        domain,
		Description:   description,
		IPV4Addresses: ipV4AddressesSlice,
		IPV6Addresses: ipV6AddressesSlice,
	}
	dnsRecord, err := c.DNSRecords.Create(dr)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(dnsRecord.ID)
	return diags
}

// resourceDnsRecordRead retrieves a DNS record from CloudConnexa and updates the Terraform state.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The provider meta interface
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred
func resourceDnsRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := d.Id()
	r, err := c.DNSRecords.GetDNSRecord(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get DNS record with ID: %s, %s", id, err)...)
	}
	if r == nil {
		d.SetId("")
	} else {
		d.Set("domain", r.Domain)
		d.Set("description", r.Description)
		d.Set("ip_v4_addresses", r.IPV4Addresses)
		d.Set("ip_v6_addresses", r.IPV6Addresses)
	}
	return diags
}

// resourceDnsRecordUpdate updates an existing DNS record in CloudConnexa.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The provider meta interface
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred
func resourceDnsRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	_, domain := d.GetChange("domain")
	_, description := d.GetChange("description")
	_, ipV4Addresses := d.GetChange("ip_v4_addresses")
	ipV4AddressesSlice := getAddressesSlice(ipV4Addresses.([]interface{}))
	_, ipV6Addresses := d.GetChange("ip_v6_addresses")
	ipV6AddressesSlice := getAddressesSlice(ipV6Addresses.([]interface{}))
	dr := cloudconnexa.DNSRecord{
		ID:            d.Id(),
		Domain:        domain.(string),
		Description:   description.(string),
		IPV4Addresses: ipV4AddressesSlice,
		IPV6Addresses: ipV6AddressesSlice,
	}
	err := c.DNSRecords.Update(dr)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

// resourceDnsRecordDelete removes a DNS record from CloudConnexa.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data
//   - m: The provider meta interface
//
// Returns:
//   - diag.Diagnostics: Diagnostics containing any errors that occurred
func resourceDnsRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	routeId := d.Id()
	err := c.DNSRecords.Delete(routeId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

// getAddressesSlice converts a slice of interface{} to a slice of strings.
//
// Parameters:
//   - addresses: A slice of interface{} containing address strings
//
// Returns:
//   - []string: A slice of string addresses
func getAddressesSlice(addresses []interface{}) []string {
	addressesSlice := make([]string, 0)
	for _, a := range addresses {
		addressesSlice = append(addressesSlice, a.(string))
	}
	return addressesSlice
}

// validateAtLeastOneNonEmptyList ensures that at least one of the IP address lists is not empty.
//
// Parameters:
//   - c: The context for the operation
//   - diff: The Terraform resource diff
//   - i: The provider meta interface
//
// Returns:
//   - error: An error if both IP address lists are empty
func validateAtLeastOneNonEmptyList(c context.Context, diff *schema.ResourceDiff, i interface{}) error {
	listA := diff.Get("ip_v4_addresses").([]interface{})
	listB := diff.Get("ip_v6_addresses").([]interface{})
	if len(listA) == 0 && len(listB) == 0 {
		return fmt.Errorf("either 'ip_v4_addresses' or 'ip_v6_addresses' must contain at least one item")
	}
	return nil
}
