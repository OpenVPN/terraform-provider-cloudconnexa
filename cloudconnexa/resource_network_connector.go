package cloudconnexa

import (
	"context"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceNetworkConnector returns a Terraform resource schema for managing network connectors
func resourceNetworkConnector() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_connector` to create an CloudConnexa connector.\n\n~> NOTE: This only creates the CloudConnexa connector object. Additional manual steps are required to associate a host in your infrastructure with the connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
		CreateContext: resourceNetworkConnectorCreate,
		ReadContext:   resourceNetworkConnectorRead,
		DeleteContext: resourceNetworkConnectorDelete,
		UpdateContext: resourceNetworkConnectorUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The connector display name.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
				Description:  "The description for the UI. Defaults to `Managed by Terraform`.",
			},
			"vpn_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the region where the connector will be deployed. Actual list of available regions can be obtained from data_source_vpn_regions.",
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the network with which this connector is associated.",
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
			"profile": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "OpenVPN profile of the connector.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Connector token.",
			},
			"ipsec_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     ipSecConfigSchema(),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ACTIVE",
				Description:  "The status of the connector. Valid values are `ACTIVE` or `SUSPENDED`. When set to `SUSPENDED`, the connector will be suspended. Note: This is a write-only field - the API does not return connector status.",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "SUSPENDED"}, false),
			},
			"connection_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The connection status of the connector.",
			},
		},
	}
}

// ipSecConfigSchema returns a Terraform resource schema for managing ipsec config
func ipSecConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"platform": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "CISCO", "AZURE", "GCP", "OTHER"}, false),
				Required:     true,
			},
			"authentication_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"SHARED_SECRET", "CERTIFICATE"}, false),
				Required:     true,
			},
			"remote_site_public_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pre_shared_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ca_certificate": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"peer_certificate": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"remote_gateway_certificate": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"peer_certificate_private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"peer_certificate_key_passphrase": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"protocol_version": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"IKE_V1", "IKE_V2"}, false),
				Required:     true,
			},
			"startup_action": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"START", "ATTACH"}, false),
				Required:     true,
			},
			"phase_1_encryption_algorithms": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"AES128", "AES256", "AES128_GCM_16", "AES256_GCM_16"}, false),
				},
			},
			"phase_1_integrity_algorithms": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"SHA1", "SHA2_256", "SHA2_384", "SHA2_512"}, false),
				},
			},
			"phase_1_diffie_hellman_groups": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"G_1", "G_2", "G_5", "G_14", "G_15", "G_16", "G_19", "G_20", "G_24"}, false),
				},
			},
			"phase_1_lifetime_sec": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"phase_2_encryption_algorithms": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"AES128", "AES256", "AES128_GCM_16", "AES256_GCM_16"}, false),
				},
			},
			"phase_2_integrity_algorithms": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"SHA1", "SHA2_256", "SHA2_384", "SHA2_512"}, false),
				},
			},
			"phase_2_diffie_hellman_groups": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"G_1", "G_2", "G_5", "G_14", "G_15", "G_16", "G_19", "G_20", "G_24"}, false),
				},
			},
			"phase_2_lifetime_sec": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"margin_time_sec": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"fuzz_percent": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"replay_window_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"timeout_sec": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"dead_peer_handling": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"RESTART", "NONE"}, false),
				Required:     true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// resourceNetworkConnectorUpdate updates an existing network connector
func resourceNetworkConnectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	// Handle status change (suspend/activate)
	if d.HasChange("status") {
		_, newStatus := d.GetChange("status")
		switch newStatus.(string) {
		case "SUSPENDED":
			if err := c.NetworkConnectors.Suspend(d.Id()); err != nil {
				return append(diags, diag.FromErr(err)...)
			}
		case "ACTIVE":
			if err := c.NetworkConnectors.Activate(d.Id()); err != nil {
				return append(diags, diag.FromErr(err)...)
			}
		}
	}

	// Handle other field changes
	if d.HasChanges("name", "description", "vpn_region_id", "ipsec_config") {
		connector := resourceDataToNetworkConnector(d)
		_, err := c.NetworkConnectors.Update(connector)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		if connector.IPSecConfig != nil {
			err := c.NetworkConnectors.StartIPsec(connector.ID)
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
		}
	}

	return resourceNetworkConnectorRead(ctx, d, m)
}

// resourceNetworkConnectorCreate creates a new network connector
func resourceNetworkConnectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	connector := resourceDataToNetworkConnector(d)
	conn, err := c.NetworkConnectors.Create(connector, connector.NetworkItemID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(conn.ID)
	if conn.TunnelingProtocol == "OPENVPN" {
		profile, err := c.NetworkConnectors.GetProfile(conn.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("profile", profile)
		token, err := c.NetworkConnectors.GetToken(conn.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("token", token)
	}

	if conn.IPSecConfig != nil {
		err := c.NetworkConnectors.StartIPsec(conn.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Connector needs to be set up manually",
		Detail:   "Terraform only creates the CloudConnexa connector object, but additional manual steps are required to associate a host in your infrastructure with this connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
	})
}

// resourceNetworkConnectorRead reads the state of a network connector
func resourceNetworkConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	id := d.Id()
	connector, err := c.NetworkConnectors.GetByID(id)
	if err != nil {
		return append(diags, diag.Errorf("Failed to get network connector with ID: %s, %s", id, err)...)
	}
	setNetworkConnectorData(d, connector)

	if connector.TunnelingProtocol == "OPENVPN" {
		token, err := c.NetworkConnectors.GetToken(connector.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("token", token)
		profile, err := c.NetworkConnectors.GetProfile(connector.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		d.Set("profile", profile)
	}
	return diags
}

// resourceNetworkConnectorDelete deletes a network connector
func resourceNetworkConnectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics
	err := c.NetworkConnectors.Delete(d.Id(), d.Get("network_id").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

// resourceDataToNetworkConnector converts Terraform resource data to a CloudConnexa network connector
func resourceDataToNetworkConnector(data *schema.ResourceData) cloudconnexa.NetworkConnector {
	name := data.Get("name").(string)
	description := data.Get("description").(string)
	networkItemId := data.Get("network_id").(string)
	vpnRegionId := data.Get("vpn_region_id").(string)
	connector := cloudconnexa.NetworkConnector{
		ID:              data.Id(),
		Name:            name,
		NetworkItemID:   networkItemId,
		NetworkItemType: "NETWORK",
		VpnRegionID:     vpnRegionId,
		Description:     description,
	}
	ipSecConfigs := data.Get("ipsec_config").([]interface{})
	if len(ipSecConfigs) > 0 {
		ipSecConfigData := ipSecConfigs[0].(map[string]interface{})
		ipSecConfig := &cloudconnexa.IPSecConfig{
			Platform:                     ipSecConfigData["platform"].(string),
			AuthenticationType:           ipSecConfigData["authentication_type"].(string),
			RemoteSitePublicIP:           ipSecConfigData["remote_site_public_ip"].(string),
			PreSharedKey:                 ipSecConfigData["pre_shared_key"].(string),
			CaCertificate:                ipSecConfigData["ca_certificate"].(string),
			PeerCertificate:              ipSecConfigData["peer_certificate"].(string),
			RemoteGatewayCertificate:     ipSecConfigData["remote_gateway_certificate"].(string),
			PeerCertificatePrivateKey:    ipSecConfigData["peer_certificate_private_key"].(string),
			PeerCertificateKeyPassphrase: ipSecConfigData["peer_certificate_key_passphrase"].(string),
			IkeProtocol: cloudconnexa.IkeProtocol{
				ProtocolVersion: ipSecConfigData["protocol_version"].(string),
				Phase1: cloudconnexa.Phase{
					EncryptionAlgorithms: toStrings(ipSecConfigData["phase_1_encryption_algorithms"].([]interface{})),
					IntegrityAlgorithms:  toStrings(ipSecConfigData["phase_1_integrity_algorithms"].([]interface{})),
					DiffieHellmanGroups:  toStrings(ipSecConfigData["phase_1_diffie_hellman_groups"].([]interface{})),
					LifetimeSec:          ipSecConfigData["phase_1_lifetime_sec"].(int),
				},
				Phase2: cloudconnexa.Phase{
					EncryptionAlgorithms: toStrings(ipSecConfigData["phase_2_encryption_algorithms"].([]interface{})),
					IntegrityAlgorithms:  toStrings(ipSecConfigData["phase_2_integrity_algorithms"].([]interface{})),
					DiffieHellmanGroups:  toStrings(ipSecConfigData["phase_2_diffie_hellman_groups"].([]interface{})),
					LifetimeSec:          ipSecConfigData["phase_2_lifetime_sec"].(int),
				},
				Rekey: cloudconnexa.Rekey{
					MarginTimeSec:    ipSecConfigData["margin_time_sec"].(int),
					FuzzPercent:      ipSecConfigData["fuzz_percent"].(int),
					ReplayWindowSize: ipSecConfigData["replay_window_size"].(int),
				},
				DeadPeerDetection: cloudconnexa.DeadPeerDetection{
					TimeoutSec:       ipSecConfigData["timeout_sec"].(int),
					DeadPeerHandling: ipSecConfigData["dead_peer_handling"].(string),
				},
				StartupAction: ipSecConfigData["startup_action"].(string),
			},
			Hostname: ipSecConfigData["hostname"].(string),
			Domain:   ipSecConfigData["domain"].(string),
		}
		connector.IPSecConfig = ipSecConfig
	}
	return connector
}

// setNetworkConnectorData sets the Terraform resource data from a CloudConnexa network connector
func setNetworkConnectorData(d *schema.ResourceData, connector *cloudconnexa.NetworkConnector) {
	d.SetId(connector.ID)
	d.Set("name", connector.Name)
	d.Set("description", connector.Description)
	d.Set("vpn_region_id", connector.VpnRegionID)
	d.Set("network_id", connector.NetworkItemID)
	d.Set("ip_v4_address", connector.IPv4Address)
	d.Set("ip_v6_address", connector.IPv6Address)
	d.Set("connection_status", connector.ConnectionStatus)
	// Note: status is not read from API as SDK doesn't support it.
	// Terraform manages status locally via suspend/activate operations.
	if connector.IPSecConfig != nil {
		ipSecConfig := make(map[string]interface{})
		ipSecConfig["platform"] = connector.IPSecConfig.Platform
		ipSecConfig["authentication_type"] = connector.IPSecConfig.AuthenticationType
		ipSecConfig["remote_site_public_ip"] = connector.IPSecConfig.RemoteSitePublicIP
		ipSecConfig["pre_shared_key"] = connector.IPSecConfig.PreSharedKey
		ipSecConfig["ca_certificate"] = connector.IPSecConfig.CaCertificate
		ipSecConfig["peer_certificate"] = connector.IPSecConfig.PeerCertificate
		ipSecConfig["remote_gateway_certificate"] = connector.IPSecConfig.RemoteGatewayCertificate
		ipSecConfig["peer_certificate_private_key"] = connector.IPSecConfig.PeerCertificatePrivateKey
		ipSecConfig["peer_certificate_key_passphrase"] = connector.IPSecConfig.PeerCertificateKeyPassphrase
		ipSecConfig["protocol_version"] = connector.IPSecConfig.IkeProtocol.ProtocolVersion
		ipSecConfig["phase_1_encryption_algorithms"] = connector.IPSecConfig.IkeProtocol.Phase1.EncryptionAlgorithms
		ipSecConfig["phase_1_integrity_algorithms"] = connector.IPSecConfig.IkeProtocol.Phase1.IntegrityAlgorithms
		ipSecConfig["phase_1_diffie_hellman_groups"] = connector.IPSecConfig.IkeProtocol.Phase1.DiffieHellmanGroups
		ipSecConfig["phase_1_lifetime_sec"] = connector.IPSecConfig.IkeProtocol.Phase1.LifetimeSec
		ipSecConfig["phase_2_encryption_algorithms"] = connector.IPSecConfig.IkeProtocol.Phase2.EncryptionAlgorithms
		ipSecConfig["phase_2_integrity_algorithms"] = connector.IPSecConfig.IkeProtocol.Phase2.IntegrityAlgorithms
		ipSecConfig["phase_2_diffie_hellman_groups"] = connector.IPSecConfig.IkeProtocol.Phase2.DiffieHellmanGroups
		ipSecConfig["phase_2_lifetime_sec"] = connector.IPSecConfig.IkeProtocol.Phase2.LifetimeSec
		ipSecConfig["margin_time_sec"] = connector.IPSecConfig.IkeProtocol.Rekey.MarginTimeSec
		ipSecConfig["replay_window_size"] = connector.IPSecConfig.IkeProtocol.Rekey.ReplayWindowSize
		ipSecConfig["fuzz_percent"] = connector.IPSecConfig.IkeProtocol.Rekey.FuzzPercent
		ipSecConfig["timeout_sec"] = connector.IPSecConfig.IkeProtocol.DeadPeerDetection.TimeoutSec
		ipSecConfig["dead_peer_handling"] = connector.IPSecConfig.IkeProtocol.DeadPeerDetection.DeadPeerHandling
		ipSecConfig["startup_action"] = connector.IPSecConfig.IkeProtocol.StartupAction
		ipSecConfig["hostname"] = connector.IPSecConfig.Hostname
		ipSecConfig["domain"] = connector.IPSecConfig.Domain
		d.Set("ipsec_config", []interface{}{ipSecConfig})
	}
}

func toStrings(strings []interface{}) []string {
	array := make([]string, len(strings))
	for i, v := range strings {
		array[i] = v.(string)
	}
	return array
}
