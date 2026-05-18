package cloudconnexa

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAccCloudConnexaHost_basic exercises the full create / read / update /
// import / destroy lifecycle of the host resource against the real
// CloudConnexa API. The update step changes every mutable field so each
// GetChange call in resourceHostUpdate observes a delta.
func TestAccCloudConnexaHost_basic(t *testing.T) {
	rn := "cloudconnexa_host.test"
	name1 := "tf-acc-" + acctest.RandStringFromCharSet(10, alphabet)
	name2 := "tf-acc-" + acctest.RandStringFromCharSet(10, alphabet)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaHostConfigBasic(name1, "first description", "host1.example.com", "SPLIT_TUNNEL_ON"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudConnexaHostExists(rn),
					resource.TestCheckResourceAttr(rn, "name", name1),
					resource.TestCheckResourceAttr(rn, "description", "first description"),
					resource.TestCheckResourceAttr(rn, "domain", "host1.example.com"),
					resource.TestCheckResourceAttr(rn, "internet_access", "SPLIT_TUNNEL_ON"),
					resource.TestCheckResourceAttrSet(rn, "id"),
				),
			},
			{
				Config: testAccCloudConnexaHostConfigBasic(name2, "second description", "host2.example.com", "RESTRICTED_INTERNET"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudConnexaHostExists(rn),
					resource.TestCheckResourceAttr(rn, "name", name2),
					resource.TestCheckResourceAttr(rn, "description", "second description"),
					resource.TestCheckResourceAttr(rn, "domain", "host2.example.com"),
					resource.TestCheckResourceAttr(rn, "internet_access", "RESTRICTED_INTERNET"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCheckCloudConnexaHostExists asserts the resource is present in state
// and the host can be fetched from the API.
func testAccCheckCloudConnexaHostExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}
		c := testAccProvider.Meta().(*providerMeta).Client
		h, err := c.Hosts.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if h == nil {
			return fmt.Errorf("host %s not found", rs.Primary.ID)
		}
		return nil
	}
}

// testAccCheckCloudConnexaHostDestroy confirms every host tracked in state
// has been removed from the API after Terraform destroys the resource.
func testAccCheckCloudConnexaHostDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*providerMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_host" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}
		h, err := c.Hosts.Get(rs.Primary.ID)
		if err == nil && h != nil {
			return fmt.Errorf("host %s still exists after destroy", rs.Primary.ID)
		}
	}
	return nil
}

// testAccCloudConnexaHostConfigBasic renders a Terraform configuration with a
// single host that uses the supplied attributes.
func testAccCloudConnexaHostConfigBasic(name, description, domain, internetAccess string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
  base_url = %q
}

resource "cloudconnexa_host" "test" {
  name            = %q
  description     = %q
  domain          = %q
  internet_access = %q
}
`, testBaseURL, name, description, domain, internetAccess)
}

// --- Unit tests -------------------------------------------------------------
//
// The acceptance test above covers every happy-path branch of resource_host.go
// when TF_ACC=1. The unit tests below cover the error branches and the
// gatewaysIdsFromResource helper using a local httptest server, so coverage
// remains meaningful even when the API is unreachable.

// TestUnitGatewaysIdsFromResource exercises both branches of
// gatewaysIdsFromResource: the empty-list short-circuit returning nil and the
// non-empty conversion of []interface{} into []string.
func TestUnitGatewaysIdsFromResource(t *testing.T) {
	t.Run("empty list returns nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceHost().Schema, map[string]interface{}{
			"name": "x",
		})
		require.Nil(t, gatewaysIdsFromResource(d))
	})

	t.Run("populated list converts to []string", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceHost().Schema, map[string]interface{}{
			"name":         "x",
			"gateways_ids": []interface{}{"gw-1", "gw-2"},
		})
		require.Equal(t, []string{"gw-1", "gw-2"}, gatewaysIdsFromResource(d))
	})
}

// newHostUnitTestClient returns a *cloudconnexa.Client wired to a local
// httptest server. The server answers /api/v1/oauth/token so
// NewClientWithOptions can complete its authentication handshake, and routes
// every other request to the supplied handler. The server is shut down via
// t.Cleanup.
func newHostUnitTestClient(t *testing.T, hostsHandler http.Handler) *cloudconnexa.Client {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/oauth/token", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"unit-test-token"}`))
	})
	mux.Handle("/", hostsHandler)

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	c, err := cloudconnexa.NewClientWithOptions(server.URL, "test-id", "test-secret", &cloudconnexa.ClientOptions{
		AllowInsecureHTTP: true,
	})
	require.NoError(t, err)
	return c
}

// hostsHandlerError returns a handler that fails every host request with the
// supplied status. It drives the CRUD functions into their error branches.
func hostsHandlerError(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{"error":"unit test failure"}`))
	})
}

// TestUnitResourceHostCreate_Success covers the happy path of
// resourceHostCreate: the server echoes the created host with an assigned ID
// and the resource persists that ID via d.SetId.
func TestUnitResourceHostCreate_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"new-host","name":"new","description":"d","domain":"e","internetAccess":"SPLIT_TUNNEL_ON","gatewaysIds":["gw-1"]}`))
	})
	c := newHostUnitTestClient(t, handler)
	d := schema.TestResourceDataRaw(t, resourceHost().Schema, map[string]interface{}{
		"name":         "new",
		"description":  "d",
		"domain":       "e",
		"gateways_ids": []interface{}{"gw-1"},
	})
	diags := resourceHostCreate(context.Background(), d, &providerMeta{Client: c, RetryConfig: defaultRetryConfig()})
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Equal(t, "new-host", d.Id())
}

// TestUnitResourceHostCreate_Error covers the error branch of
// resourceHostCreate: when Hosts.Create fails the function must surface a
// diagnostic and leave the resource ID empty.
func TestUnitResourceHostCreate_Error(t *testing.T) {
	c := newHostUnitTestClient(t, hostsHandlerError(http.StatusInternalServerError))
	d := schema.TestResourceDataRaw(t, resourceHost().Schema, map[string]interface{}{
		"name": "to-fail",
	})
	diags := resourceHostCreate(context.Background(), d, &providerMeta{Client: c, RetryConfig: defaultRetryConfig()})
	assert.True(t, diags.HasError())
	assert.Empty(t, d.Id(), "ID must remain empty when Create fails")
}

// TestUnitResourceHostRead_Success covers the happy path of resourceHostRead:
// every populated field returned by the API (including the computed
// system_subnets and gateways_ids) is persisted into state.
func TestUnitResourceHostRead_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
            "id":"host-id",
            "name":"my-host",
            "description":"desc",
            "domain":"example.com",
            "internetAccess":"SPLIT_TUNNEL_ON",
            "systemSubnets":["10.0.0.0/24"],
            "gatewaysIds":["gw-1"]
        }`))
	})
	c := newHostUnitTestClient(t, handler)
	d := schema.TestResourceDataRaw(t, resourceHost().Schema, map[string]interface{}{
		"name": "ignored",
	})
	d.SetId("host-id")
	diags := resourceHostRead(context.Background(), d, &providerMeta{Client: c, RetryConfig: defaultRetryConfig()})
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Equal(t, "my-host", d.Get("name"))
	assert.Equal(t, "desc", d.Get("description"))
	assert.Equal(t, "example.com", d.Get("domain"))
	assert.Equal(t, "SPLIT_TUNNEL_ON", d.Get("internet_access"))
	assert.Equal(t, []interface{}{"gw-1"}, d.Get("gateways_ids"))
}

// TestUnitResourceHostRead_Error covers the error branch of resourceHostRead:
// when Hosts.Get fails the function returns a descriptive diagnostic.
func TestUnitResourceHostRead_Error(t *testing.T) {
	c := newHostUnitTestClient(t, hostsHandlerError(http.StatusInternalServerError))
	d := schema.TestResourceDataRaw(t, resourceHost().Schema, map[string]interface{}{
		"name": "x",
	})
	d.SetId("host-id")
	diags := resourceHostRead(context.Background(), d, &providerMeta{Client: c, RetryConfig: defaultRetryConfig()})
	assert.True(t, diags.HasError())
}

// TestUnitResourceHostUpdate_Success covers the happy path of
// resourceHostUpdate: a successful Update is immediately followed by a Read
// that refreshes state. The handler dispatches on method so both calls
// succeed in a single test.
func TestUnitResourceHostUpdate_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPut:
			_, _ = w.Write([]byte(`{}`))
		case http.MethodGet:
			_, _ = w.Write([]byte(`{
                "id":"host-id",
                "name":"updated",
                "description":"d",
                "internetAccess":"SPLIT_TUNNEL_OFF"
            }`))
		default:
			t.Errorf("unexpected method: %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	c := newHostUnitTestClient(t, handler)
	d := schema.TestResourceDataRaw(t, resourceHost().Schema, map[string]interface{}{
		"name":            "updated",
		"description":     "d",
		"internet_access": "SPLIT_TUNNEL_OFF",
	})
	d.SetId("host-id")
	diags := resourceHostUpdate(context.Background(), d, &providerMeta{Client: c, RetryConfig: defaultRetryConfig()})
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Equal(t, "updated", d.Get("name"))
}

// TestUnitResourceHostUpdate_Error covers the error branch of
// resourceHostUpdate: when Hosts.Update fails the function returns a
// diagnostic without performing a subsequent Read.
func TestUnitResourceHostUpdate_Error(t *testing.T) {
	c := newHostUnitTestClient(t, hostsHandlerError(http.StatusInternalServerError))
	d := schema.TestResourceDataRaw(t, resourceHost().Schema, map[string]interface{}{
		"name": "x",
	})
	d.SetId("host-id")
	diags := resourceHostUpdate(context.Background(), d, &providerMeta{Client: c, RetryConfig: defaultRetryConfig()})
	assert.True(t, diags.HasError())
}

// TestUnitResourceHostDelete_Success covers the happy path of
// resourceHostDelete: a successful DELETE produces no diagnostics.
func TestUnitResourceHostDelete_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	c := newHostUnitTestClient(t, handler)
	d := schema.TestResourceDataRaw(t, resourceHost().Schema, map[string]interface{}{
		"name": "x",
	})
	d.SetId("host-id")
	diags := resourceHostDelete(context.Background(), d, &providerMeta{Client: c, RetryConfig: defaultRetryConfig()})
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
}

// TestUnitResourceHostDelete_Error covers the error branch of
// resourceHostDelete: when Hosts.Delete fails the function returns a
// diagnostic.
func TestUnitResourceHostDelete_Error(t *testing.T) {
	c := newHostUnitTestClient(t, hostsHandlerError(http.StatusInternalServerError))
	d := schema.TestResourceDataRaw(t, resourceHost().Schema, map[string]interface{}{
		"name": "x",
	})
	d.SetId("host-id")
	diags := resourceHostDelete(context.Background(), d, &providerMeta{Client: c, RetryConfig: defaultRetryConfig()})
	assert.True(t, diags.HasError())
}
