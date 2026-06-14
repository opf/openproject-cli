package configuration_test

// Regression tests for review findings #1 and #2 on PR #15: a corrupt or
// malformed config file must be reported as an error rather than silently
// yielding empty (or bogus migrated) credentials. They reuse setupTempConfig /
// writeRaw from profiles_test.go (same package).

import (
	"testing"

	"github.com/opf/openproject-cli/components/configuration"
)

// Finding #1: a corrupt config file must now be reported as an error rather than
// silently yielding empty (or, for spaced garbage, bogus migrated) credentials.
func TestFinding_CorruptConfigReportsError(t *testing.T) {
	corrupt := []struct {
		name    string
		content string
	}{
		{"garbage with spaces", "%%% not ini and not host token %%%"},
		{"garbage no spaces", "@@@garbage@@@"},
		{"malformed ini section", "[default\nhost broken"},
		{"host only, no token", "https://example.com"},
	}

	for _, c := range corrupt {
		t.Run(c.name, func(t *testing.T) {
			setupTempConfig(t)
			writeRaw(t, c.content)

			host, token, err := configuration.ReadConfig("default")
			if err == nil {
				t.Errorf("expected an error for corrupt config %q, got host=%q token=%q", c.content, host, token)
			}
		})
	}
}

// Finding #2: a malformed old-format file (single field, no usable host) is now
// surfaced as an error instead of a silent logout.
func TestFinding_MalformedOldFormatReportsError(t *testing.T) {
	setupTempConfig(t)
	writeRaw(t, "justatokenwithnospace")

	if _, _, err := configuration.ReadConfig("default"); err == nil {
		t.Error("expected an error for malformed old-format config, got nil")
	}
}

// A well-formed old file ("host token") must still migrate cleanly.
func TestFinding_WellFormedOldFormatStillMigrates(t *testing.T) {
	setupTempConfig(t)
	writeRaw(t, "https://legacy.example.com legacytoken")

	host, token, err := configuration.ReadConfig("default")
	if err != nil {
		t.Fatal(err)
	}
	if host != "https://legacy.example.com" || token != "legacytoken" {
		t.Errorf("well-formed old file should migrate, got host=%q token=%q", host, token)
	}
}
