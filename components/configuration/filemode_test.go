package configuration_test

// Copilot review (PR #15): the config file stores API tokens, so it must not be
// world-readable. WriteConfigForProfile must create it with mode 0600.

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opf/openproject-cli/components/configuration"
)

func TestWriteConfigForProfile_FileModeIs0600(t *testing.T) {
	setupTempConfig(t)

	if err := configuration.WriteConfigForProfile("default", "https://example.com", "secret-token"); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "openproject", "config")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if mode := info.Mode().Perm(); mode != 0600 {
		t.Errorf("config file mode = %04o, want 0600 (token must not be world-readable)", mode)
	}
}

// A config written by the CLI (mode 0600) must not be reported as insecure.
func TestInsecureConfigPermissions_SecureFile(t *testing.T) {
	setupTempConfig(t)

	if err := configuration.WriteConfigForProfile("default", "https://example.com", "secret-token"); err != nil {
		t.Fatal(err)
	}

	if insecure, mode := configuration.InsecureConfigPermissions(); insecure {
		t.Errorf("0600 config reported insecure (mode %#o)", mode)
	}
}

// A config readable by group or other users must be reported as insecure so the
// CLI can warn that the API token may leak.
func TestInsecureConfigPermissions_WorldReadableFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permission bits are not meaningful on Windows")
	}
	setupTempConfig(t)

	if err := configuration.WriteConfigForProfile("default", "https://example.com", "secret-token"); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(configuration.ConfigFilePath(), 0644); err != nil {
		t.Fatal(err)
	}

	insecure, mode := configuration.InsecureConfigPermissions()
	if !insecure {
		t.Errorf("0644 config not reported insecure (mode %#o)", mode)
	}
	if mode.Perm() != 0644 {
		t.Errorf("reported mode = %#o, want 0644", mode.Perm())
	}
}

// A missing config file is not insecure.
func TestInsecureConfigPermissions_MissingFile(t *testing.T) {
	setupTempConfig(t)

	if insecure, _ := configuration.InsecureConfigPermissions(); insecure {
		t.Error("missing config file reported insecure")
	}
}
