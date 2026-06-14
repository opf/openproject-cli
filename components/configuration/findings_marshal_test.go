package configuration_test

// Regression test for review finding #3 on PR #15: rewriting the config must
// preserve keys other than host/token. Reuses setupTempConfig / writeRaw from
// profiles_test.go (same package).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/configuration"
)

// readRawConfig returns the raw bytes of the config file the package writes to.
func readRawConfig(t *testing.T) string {
	t.Helper()
	path := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "openproject", "config")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading raw config: %v", err)
	}
	return string(data)
}

func TestFinding_MarshalPreservesUnknownKeys(t *testing.T) {
	setupTempConfig(t)
	writeRaw(t, "[default]\nhost = https://example.com\ntoken = tok\nextra = keepme\n")

	// Trigger a read-modify-write cycle on a different profile so the default
	// section is preserved-and-rewritten rather than edited.
	if err := configuration.WriteConfigForProfile("work", "https://work.example.com", "tok-work"); err != nil {
		t.Fatal(err)
	}

	raw := readRawConfig(t)
	if !strings.Contains(raw, "extra = keepme") {
		t.Errorf("expected unknown key 'extra = keepme' to survive rewrite, file was:\n%s", raw)
	}
}
