# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

A Terraform provider for [CloudConnexa](https://openvpn.net/cloud-vpn/), built with the HashiCorp Terraform Plugin SDK v2. The provider wraps the `github.com/openvpn/cloudconnexa-go-client/v2` SDK, which talks to the `/v1` CloudConnexa REST API. Published to the Terraform Registry as `OpenVPN/cloudconnexa`.

## Common commands

- `make build` — compile the provider binary.
- `make install` — build and install to the local Terraform plugins dir at `~/.terraform.d/plugins/cloudconnexa.dev/openvpn/cloudconnexa/${VERSION}/${OS_ARCH}`. Local Terraform configs can then `source = "cloudconnexa.dev/openvpn/cloudconnexa"` to exercise an unreleased build. Bump `VERSION` in the Makefile when the local plugin doesn't get picked up.
- `make lint` — `golangci-lint run ./... --disable errcheck` (matches CI; CI pins v2.12.2). `errcheck` is intentionally disabled — don't re-enable it.
- `make test` — the Makefile target is fragile (`go test -i` is a Go ≤1.9 flag). CI's actual command is `go test -v -cover ./cloudconnexa` — prefer that.
- `make testacc` — full acceptance run (`TF_ACC=1`). Creates real CloudConnexa resources via the API; costs money.
- `make docs` / `make docs-check` — regenerate / validate `docs/` from schema descriptions + `templates/` via `tfplugindocs`. Run `make docs` after changing any resource schema or its `Description`.

Single test: `go test ./cloudconnexa -run TestUnitResourceHostCreate_Success -v`. Acceptance tests are gated by `TF_ACC`; unit tests (the `TestUnit*` family) always run.

## Test environment loading

`cloudconnexa/main_test.go` defines a `TestMain` that auto-loads `.env` and `../.env` (simple `KEY=VALUE` lines) before tests run. **Side effect:** if a `.env` was applied, `TF_ACC=1` is set automatically unless you exported `TF_ACC` yourself. So merely having a `.env` with `CLOUDCONNEXA_*` populated makes `go test ./cloudconnexa` hit the real API. Set `TF_ACC=0` explicitly to opt out.

Required env vars for acceptance tests:
- `CLOUDCONNEXA_CLIENT_ID`, `CLOUDCONNEXA_CLIENT_SECRET` (OAuth credentials)
- `CLOUDCONNEXA_BASE_URL` (full URL, e.g. `https://example.api.openvpn.com`) — separate from the runtime `base_url`/`cloud_id` provider arg; only the acceptance harness reads it.

## Architecture

Single Go package `cloudconnexa/` exports a `Provider()` (`provider.go`) that wires every resource and data source explicitly into `ResourcesMap` / `DataSourcesMap`. **A new resource or data source must be registered there or it won't be visible.** `main.go` is just the plugin entrypoint.

The provider accepts `client_id` + `client_secret` plus exactly one of `base_url` or `cloud_id` (enforced via `ExactlyOneOf`). `cloud_id` is shorthand: the configure function expands it to `https://<cloud_id>.api.openvpn.com`. After `providerConfigure` runs, the provider's "meta" passed into every CRUD function is a `*cloudconnexa.Client`; resources start with `c := m.(*cloudconnexa.Client)` and call `c.Hosts.Create`, `c.Users.Get`, `c.Devices.GetByID`, etc.

Each resource lives in its own `resource_<name>.go` / `data_source_<name>.go`. The shape is consistent: a builder returning `*schema.Resource` with `CreateContext` / `ReadContext` / `UpdateContext` / `DeleteContext` / `Importer`, then four CRUD funcs that translate between Terraform's `*schema.ResourceData` and the SDK's structs. `toStrings([]interface{}) []string` in `resource_network_connector.go` is the shared helper for converting list attributes — reuse it rather than reimplementing.

Composite import IDs: most resources use `schema.ImportStatePassthroughContext`, but `cloudconnexa_device` parses `user_id/device_id` in a custom `resourceDeviceImport` because the SDK's `Devices.GetByID(userID, id)` needs both. Mirror this pattern for any child resource whose parent isn't already in state.

The `devices` block on `cloudconnexa_user` is deprecated — managing devices inline is being removed in a future major release. New device work goes through the standalone `cloudconnexa_device` resource and `cloudconnexa_device`/`cloudconnexa_devices` data sources.

## Test patterns

Two styles coexist in `cloudconnexa/`:

- `TestAccXxx_*` — acceptance tests that drive `resource.Test` against the real API via `testAccProviderFactories` after `testAccPreCheck`. Cover at minimum a create step, an update step that flips every mutable field (so `d.GetChange` observes a delta in `Update`), and an `ImportState`/`ImportStateVerify` step. Pair each with a `testAccCheck<Resource>Destroy` that confirms the API no longer returns the resource.
- `TestUnitXxx_*` — unit tests that spin up an `httptest.NewServer` exposing `/api/v1/oauth/token` plus a per-test handler, then construct the SDK client via `cloudconnexa.NewClientWithOptions(server.URL, ..., &cloudconnexa.ClientOptions{AllowInsecureHTTP: true})`. See `newHostUnitTestClient` in `resource_host_test.go` for the canonical setup. Use these to cover the error branches (`*_Error` tests with `hostsHandlerError`) so coverage doesn't depend on the API being reachable.

`provider_test.go` initializes `testAccProvider` / `testAccProviderFactories` in `init()` and defines `testAccPreCheck`. The `alphabet` constant + `acctest.RandStringFromCharSet` is the project's standard for unique acceptance-test resource names.

## End-to-end suite (`e2e/`)

`e2e/integration_test.go` is separate from the provider-package tests: it uses Terratest to apply `./e2e/setup` (provisions an EC2 instance running an OpenVPN connector), polls the CloudConnexa API for up to 60 s waiting for the connector to come online, then destroys. It expects `OVPN_HOST` set in addition to the standard client ID/secret. Don't run it casually — it costs real AWS + CloudConnexa resources.

## Docs generation

`docs/` is generated, not hand-edited. Source of truth is the `Description` fields on each schema attribute plus templates in `templates/`. Only two resources have explicit templates (`host_connector`, `network_connector`) because their schemas are large; everything else uses the default. Examples that get inlined into the generated docs live at `examples/resources/<resource-name>/resource.tf` (and optional `import.sh`) — keep them runnable.

## Releases

Tag-driven (`v*`) via `.github/workflows/release.yml` → GoReleaser. Cross-builds darwin/linux/windows/freebsd × amd64/386/arm/arm64, signs checksums with GPG, and publishes a **draft** GitHub release for manual review before going live. The `version` constant in `cloudconnexa/provider.go` is the user-agent string sent to the API — bump it when cutting a release.
