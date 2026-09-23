package cloudconnexa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// connectorEphemeral fetches a connector credential, a setup token or the OpenVPN profile, on
// every Open. The value is returned to Terraform as an ephemeral value, so it is never written to
// the plan or the state. A token Open mints a new token and leaves earlier ones valid; a profile
// Open returns the connector's current profile. Neither needs Renew or Close.
//
// The same implementation serves network and host connectors and both credentials; only the type
// name suffix, the client call and the wording differ.
type connectorEphemeral struct {
	kind                 string // "network" or "host"
	attribute            string // "token" or "profile": the result attribute and the type name suffix
	verb                 string // "mint" or "retrieve"
	description          string
	attributeDescription string
	client               *cloudconnexa.Client
	fetch                func(c *cloudconnexa.Client, connectorID string) (string, error)
}

var _ ephemeral.EphemeralResource = (*connectorEphemeral)(nil)
var _ ephemeral.EphemeralResourceWithConfigure = (*connectorEphemeral)(nil)

func newNetworkConnectorTokenEphemeral() ephemeral.EphemeralResource {
	return newConnectorTokenEphemeral("network", func(c *cloudconnexa.Client, id string) (string, error) {
		if err := requireOpenVPNNetworkConnector(c, id); err != nil {
			return "", err
		}
		return c.NetworkConnectors.GetToken(id)
	})
}

func newHostConnectorTokenEphemeral() ephemeral.EphemeralResource {
	return newConnectorTokenEphemeral("host", func(c *cloudconnexa.Client, id string) (string, error) {
		return c.HostConnectors.GetToken(id)
	})
}

func newNetworkConnectorProfileEphemeral() ephemeral.EphemeralResource {
	return newConnectorProfileEphemeral("network", func(c *cloudconnexa.Client, id string) (string, error) {
		if err := requireOpenVPNNetworkConnector(c, id); err != nil {
			return "", err
		}
		return c.NetworkConnectors.GetProfile(id)
	})
}

func newHostConnectorProfileEphemeral() ephemeral.EphemeralResource {
	return newConnectorProfileEphemeral("host", func(c *cloudconnexa.Client, id string) (string, error) {
		return c.HostConnectors.GetProfile(id)
	})
}

func newConnectorTokenEphemeral(kind string, fetch func(c *cloudconnexa.Client, connectorID string) (string, error)) *connectorEphemeral {
	return &connectorEphemeral{
		kind:      kind,
		attribute: "token",
		verb:      "mint",
		fetch:     fetch,
		description: fmt.Sprintf("Mints a setup token for an existing OpenVPN %s connector. "+
			"The token is an ephemeral value: it is never stored in the plan or the state and can only be used in "+
			"ephemeral contexts such as write-only arguments, provisioner and connection blocks, or other ephemeral resources. "+
			"A new token is minted every time Terraform opens this resource, which happens once during plan and once during apply, "+
			"so declare it only where a token is actually consumed. Minting a token does not invalidate tokens minted earlier.", kind),
		attributeDescription: "The setup token, as accepted by `openvpn-connector-setup --token` and the connector deployment templates.",
	}
}

func newConnectorProfileEphemeral(kind string, fetch func(c *cloudconnexa.Client, connectorID string) (string, error)) *connectorEphemeral {
	return &connectorEphemeral{
		kind:      kind,
		attribute: "profile",
		verb:      "retrieve",
		fetch:     fetch,
		description: fmt.Sprintf("Retrieves the `.ovpn` profile of an existing OpenVPN %s connector. "+
			"The profile is an ephemeral value: it is never stored in the plan or the state and can only be used in "+
			"ephemeral contexts such as write-only arguments, provisioner and connection blocks, or other ephemeral resources. "+
			"Terraform opens this resource once during plan and once during apply; each open returns the connector's current "+
			"profile, it does not generate a new one.", kind),
		attributeDescription: "The OpenVPN profile of the connector, in `.ovpn` format.",
	}
}

// requireOpenVPNNetworkConnector fails unless the network connector uses the OpenVPN tunneling
// protocol: IPsec connectors have neither a profile nor a token, and both endpoints fail for them.
func requireOpenVPNNetworkConnector(c *cloudconnexa.Client, id string) error {
	connector, err := c.NetworkConnectors.GetByID(id)
	if err != nil {
		return err
	}
	if connector == nil {
		return fmt.Errorf("network connector not found")
	}
	if connector.TunnelingProtocol != "OPENVPN" {
		return fmt.Errorf("network connector uses tunneling protocol %s; profiles and tokens exist only for OPENVPN connectors", connector.TunnelingProtocol)
	}
	return nil
}

func (e *connectorEphemeral) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + e.kind + "_connector_" + e.attribute
}

func (e *connectorEphemeral) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: e.description,
		Attributes: map[string]schema.Attribute{
			"connector_id": schema.StringAttribute{
				Description: fmt.Sprintf("The ID of the %s connector to %s a %s for.", e.kind, e.verb, e.attribute),
				Required:    true,
			},
			e.attribute: schema.StringAttribute{
				Description: e.attributeDescription,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (e *connectorEphemeral) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return // provider not configured yet (e.g. during validation)
	}
	client, ok := req.ProviderData.(*cloudconnexa.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data",
			fmt.Sprintf("Expected *cloudconnexa.Client, got %T. This is a bug in the provider.", req.ProviderData))
		return
	}
	e.client = client
}

// Open only sets the credential attribute: the framework starts the result as a copy of the
// config, so connector_id is already in place.
func (e *connectorEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var connectorID types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("connector_id"), &connectorID)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if e.client == nil {
		resp.Diagnostics.AddError("Provider not configured",
			fmt.Sprintf("The CloudConnexa provider must be configured before a connector %s can be requested.", e.attribute))
		return
	}

	id := connectorID.ValueString()
	value, err := e.fetch(e.client, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to %s %s connector %s", e.verb, e.kind, e.attribute),
			fmt.Sprintf("Connector %s: %s", id, err))
		return
	}
	// Both endpoints answer 202 with an empty body while the connector profile is still being generated.
	if value == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s connector profile is not ready", e.kind),
			fmt.Sprintf("The API accepted the %s request for connector %s but the profile is still being generated. Retry in a few seconds.", e.attribute, id))
		return
	}

	resp.Diagnostics.Append(resp.Result.SetAttribute(ctx, path.Root(e.attribute), types.StringValue(value))...)
}
