package configuration_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/configuration"
	"github.com/opf/openproject-cli/components/printer"
)

func configPath() string {
	return filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "openproject", "config")
}

// Migration used to rewrite the config in place, which kept a permissive
// pre-existing mode (e.g. 0644) on a file that stores API tokens. The
// atomic rewrite must leave the migrated file with mode 0600.
func TestMigration_resetsInsecureFileMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix permission bits are not meaningful on windows")
	}
	setupTempConfig(t)
	writeRaw(t, "https://legacy.example.com legacytoken")

	if _, _, err := configuration.ReadConfig("default"); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(configPath())
	if err != nil {
		t.Fatal(err)
	}
	if mode := info.Mode().Perm(); mode != 0600 {
		t.Errorf("migrated config mode = %#o, want 0600", mode)
	}
}

// A readable but non-writable legacy config must still yield its credentials;
// persisting the migrated format is best effort.
func TestMigration_readOnlyConfigStillReadable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory write permissions are not meaningful on windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root bypasses permission checks")
	}
	printer.Init(&printer.TestingPrinter{})
	setupTempConfig(t)
	writeRaw(t, "https://legacy.example.com legacytoken")

	dir := filepath.Dir(configPath())
	if err := os.Chmod(dir, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0700) })

	host, token, err := configuration.ReadConfig("default")
	if err != nil {
		t.Fatalf("ReadConfig on read-only legacy config: %v", err)
	}
	if host != "https://legacy.example.com" || token != "legacytoken" {
		t.Errorf("got host %q token %q", host, token)
	}
}

// Hand-edited files can contain duplicate sections; logout must remove all of
// them, or an earlier section's token silently stays active.
func TestDeleteProfile_removesDuplicateSections(t *testing.T) {
	setupTempConfig(t)
	writeRaw(t, "[work]\nhost = https://a.example.com\ntoken = tok-a\n\n[work]\nhost = https://b.example.com\ntoken = tok-b\n")

	if err := configuration.DeleteProfile("work"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(configPath())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "[work]") {
		t.Errorf("config still contains a [work] section after delete:\n%s", data)
	}
}

// A host that is not an absolute URL must be rejected on read with a clear
// message instead of surfacing later as a confusing request failure.
func TestReadConfig_rejectsInvalidHost(t *testing.T) {
	setupTempConfig(t)
	writeRaw(t, "[default]\nhost = garbage\ntoken = tok\n")

	_, _, err := configuration.ReadConfig("default")
	if err == nil {
		t.Fatal("expected error for non-URL host")
	}
	if !strings.Contains(err.Error(), "invalid host") {
		t.Errorf("error should mention invalid host, got: %v", err)
	}
}
