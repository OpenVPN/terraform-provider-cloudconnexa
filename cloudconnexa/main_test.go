package cloudconnexa

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestMain loads CloudConnexa credentials from a local .env file before
// running the package tests. The file is optional: when missing, tests fall
// back to whatever is already in the environment (CI populates the same
// variables via GitHub secrets).
//
// When a .env file is loaded, TF_ACC is set to "1" by default so acceptance
// tests participate in the run. Callers that want to opt out can export
// TF_ACC=0 explicitly before invoking `go test`.
func TestMain(m *testing.M) {
	var loaded bool
	for _, p := range []string{".env", "../.env"} {
		if loadDotEnv(p) {
			loaded = true
		}
	}
	if loaded {
		if _, set := os.LookupEnv(resource.EnvTfAcc); !set {
			_ = os.Setenv(resource.EnvTfAcc, "1")
		}
	}
	// testBaseURL was initialized at package-init time, before TestMain ran,
	// so reseat it now that .env has populated the environment.
	testBaseURL = os.Getenv(BaseURLEnvVar)
	os.Exit(m.Run())
}

// loadDotEnv reads simple KEY=VALUE pairs from the file at path and exports
// them as process environment variables. Values may be wrapped in single or
// double quotes; comments (lines beginning with #) and blank lines are
// skipped. Keys already present in the environment are left untouched so an
// explicit override (e.g. from CI) still wins.
//
// Returns true if the file was opened and at least one key/value pair was
// applied.
func loadDotEnv(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	f, err := os.Open(abs)
	if err != nil {
		return false
	}
	defer func() { _ = f.Close() }()

	applied := false
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		eq := strings.IndexByte(line, '=')
		if eq <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		val = strings.Trim(val, `"'`)
		if _, set := os.LookupEnv(key); !set {
			if err := os.Setenv(key, val); err == nil {
				applied = true
			}
		}
	}
	return applied
}
