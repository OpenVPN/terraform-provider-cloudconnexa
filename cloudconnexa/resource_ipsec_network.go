package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// resourceIpsecNetwork returns a Terraform resource for managing CloudConnexa IPsec networks
func resourceIpsecNetwork() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_ipsec_network` to create and manage CloudConnexa IPsec networks with dedicated IPsec tunnel support.",
		CreateContext: resourceIpsecNetworkCreate,
		ReadContext:   resourceIpsecNetworkRead,
		DeleteContext: resourceIpsecNetworkDelete,
		UpdateContext: resourceIpsecNetworkUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 40),
				Description:  "The name of the IPsec network.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "The description for the IPsec network. Defaults to `Managed by Terraform`.",
			},
			"egress": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Boolean to control if this network can be used as an egress network or not.",
			},
			"internet_access": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "SPLIT_TUNNEL_ON",
				ValidateFunc: validation.StringInSlice([]string{"SPLIT_TUNNEL_ON", "SPLIT_TUNNEL_OFF", "RESTRICTED_INTERNET"}, false),
				Description:  "The type of internet access provided. Valid values are `SPLIT_TUNNEL_ON`, `SPLIT_TUNNEL_OFF`, or `RESTRICTED_INTERNET`. Defaults to `SPLIT_TUNNEL_ON`.",
			},
			"ipsec_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "IPsec tunnel configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remote_gateway": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Remote IPsec gateway IP address or FQDN.",
						},
						"remote_networks": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of remote network subnets accessible through the IPsec tunnel.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"pre_shared_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Pre-shared key for IPsec authentication.",
						},
						"ike_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "IKEv2",
							ValidateFunc: validation.StringInSlice([]string{"IKEv1", "IKEv2"}, false),
							Description:  "IKE version to use. Valid values are `IKEv1` or `IKEv2`. Defaults to `IKEv2`.",
						},
						"encryption_algorithm": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "AES256",
							ValidateFunc: validation.StringInSlice([]string{"AES128", "AES256", "3DES"}, false),
							Description:  "Encryption algorithm. Valid values are `AES128`, `AES256`, or `3DES`. Defaults to `AES256`.",
						},
						"hash_algorithm": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "SHA256",
							ValidateFunc: validation.StringInSlice([]string{"SHA1", "SHA256", "SHA384", "SHA512", "MD5"}, false),
							Description:  "Hash algorithm. Valid values are `SHA1`, `SHA256`, `SHA384`, `SHA512`, or `MD5`. Defaults to `SHA256`.",
						},
						"dh_group": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "14",
							ValidateFunc: validation.StringInSlice([]string{"1", "2", "5", "14", "15", "16", "17", "18", "19", "20", "21"}, false),
							Description:  "Diffie-Hellman group. Valid values are `1`, `2`, `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`. Defaults to `14`.",
						},
						"pfs_group": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "14",
							ValidateFunc: validation.StringInSlice([]string{"1", "2", "5", "14", "15", "16", "17", "18", "19", "20", "21"}, false),
							Description:  "Perfect Forward Secrecy group. Valid values are `1`, `2`, `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`. Defaults to `14`.",
						},
						"ike_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      28800,
							ValidateFunc: validation.IntBetween(3600, 86400),
							Description:  "IKE lifetime in seconds. Must be between 3600 and 86400. Defaults to 28800.",
						},
						"ipsec_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntBetween(1800, 86400),
							Description:  "IPsec lifetime in seconds. Must be between 1800 and 86400. Defaults to 3600.",
						},
						"dpd_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      30,
							ValidateFunc: validation.IntBetween(10, 300),
							Description:  "Dead Peer Detection timeout in seconds. Must be between 10 and 300. Defaults to 30.",
						},
						"nat_traversal": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Enable NAT traversal. Defaults to true.",
						},
					},
				},
			},
			"connector": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Network connector configuration for the IPsec network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The connector display name.",
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "IPsec connector managed by Terraform",
							ValidateFunc: validation.StringLenBetween(1, 120),
							Description:  "The description for the connector. Defaults to `IPsec connector managed by Terraform`.",
						},
						"vpn_region_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the region where the connector will be deployed.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the created connector.",
						},
						"ip_v4_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPV4 address of the connector.",
						},
						"ip_v6_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPV6 address of the connector.",
						},
					},
				},
			},
			"tunnel_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current IPsec tunnel status.",
			},
		},
	}
}

// resourceIpsecNetworkCreate creates a new IPsec network with connector and IPsec configuration
func resourceIpsecNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	// Create the network first
	network := cloudconnexa.Network{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Egress:         d.Get("egress").(bool),
		InternetAccess: d.Get("internet_access").(string),
	}

	createdNetwork, err := c.Networks.Create(network)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId(createdNetwork.ID)

	// Create the connector
	connectorConfig := d.Get("connector").([]interface{})[0].(map[string]interface{})
	connector := cloudconnexa.NetworkConnector{
		Name:            connectorConfig["name"].(string),
		Description:     connectorConfig["description"].(string),
		NetworkItemID:   createdNetwork.ID,
		NetworkItemType: "NETWORK",
		VpnRegionID:     connectorConfig["vpn_region_id"].(string),
	}

	createdConnector, err := c.NetworkConnectors.Create(connector, createdNetwork.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// Start IPsec tunnel
	_, err = c.NetworkConnectors.StartIPsec(createdConnector.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// Update the connector information in the state
	connectorData := []map[string]interface{}{
		{
			"name":          createdConnector.Name,
			"description":   createdConnector.Description,
			"vpn_region_id": createdConnector.VpnRegionID,
			"id":            createdConnector.ID,
			"ip_v4_address": createdConnector.IPv4Address,
			"ip_v6_address": createdConnector.IPv6Address,
		},
	}
	d.Set("connector", connectorData)

	// Set tunnel status after starting IPsec
	d.Set("tunnel_status", "active")

	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "IPsec network created",
		Detail:   "IPsec network and connector have been created. Additional configuration may be required on the remote IPsec gateway to establish the tunnel.",
	})
}

// resourceIpsecNetworkRead reads the current state of an IPsec network
func resourceIpsecNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	network, err := c.Networks.Get(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if network == nil {
		d.SetId("")
		return diags
	}

	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("egress", network.Egress)
	d.Set("internet_access", network.InternetAccess)

	// Get connector information
	connectors, err := c.NetworkConnectors.ListByNetworkID(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if len(connectors) > 0 {
		connector := connectors[0]
		connectorData := []map[string]interface{}{
			{
				"name":          connector.Name,
				"description":   connector.Description,
				"vpn_region_id": connector.VpnRegionID,
				"id":            connector.ID,
				"ip_v4_address": connector.IPv4Address,
				"ip_v6_address": connector.IPv6Address,
			},
		}
		d.Set("connector", connectorData)

		// Set tunnel status - IPsec networks are always considered active since they auto-start IPsec
		d.Set("tunnel_status", "active")
	} else {
		d.Set("tunnel_status", "inactive")
	}

	return diags
}

// resourceIpsecNetworkUpdate updates an existing IPsec network
func resourceIpsecNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	// Update network if needed
	if d.HasChanges("name", "description", "egress", "internet_access") {
		network := cloudconnexa.Network{
			ID:             d.Id(),
			Name:           d.Get("name").(string),
			Description:    d.Get("description").(string),
			Egress:         d.Get("egress").(bool),
			InternetAccess: d.Get("internet_access").(string),
		}

		err := c.Networks.Update(network)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	// Update connector if needed
	if d.HasChange("connector") {
		connectorConfig := d.Get("connector").([]interface{})[0].(map[string]interface{})
		connectorID := connectorConfig["id"].(string)

		if connectorID != "" {
			connector := cloudconnexa.NetworkConnector{
				ID:          connectorID,
				Name:        connectorConfig["name"].(string),
				Description: connectorConfig["description"].(string),
				VpnRegionID: connectorConfig["vpn_region_id"].(string),
			}

			_, err := c.NetworkConnectors.Update(connector)
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
		}
	}

	return diags
}

// resourceIpsecNetworkDelete deletes an IPsec network and its associated connector
func resourceIpsecNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	// Get and delete connectors first
	connectors, err := c.NetworkConnectors.ListByNetworkID(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	for _, connector := range connectors {
		// Stop IPsec tunnel first
		_, err := c.NetworkConnectors.StopIPsec(connector.ID)
		if err != nil {
			// Log warning but continue with deletion
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to stop IPsec tunnel",
				Detail:   err.Error(),
			})
		}

		// Delete the connector
		err = c.NetworkConnectors.Delete(connector.ID, d.Id())
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	// Delete the network
	err = c.Networks.Delete(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId("")
	return diags
}
