package configuration_test

// Copilot review (PR #15): the config file stores API tokens, so it must not be
// world-readable. WriteConfigForProfile must create it with mode 0600.

import (
	"os"
	"path/filepath"
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
