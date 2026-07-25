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

func TestReadConfig_EnvironmentDoesNotRequireConfigDirectory(t *testing.T) {
	configRoot := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(configRoot, []byte("occupied"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", configRoot)
	t.Setenv("OP_CLI_HOST", "https://env.example.com")
	t.Setenv("OP_CLI_TOKEN", "env-token")

	host, token, err := configuration.ReadConfig("default")
	if err != nil {
		t.Fatalf("ReadConfig with environment credentials: %v", err)
	}
	if host != "https://env.example.com" || token != "env-token" {
		t.Errorf("got host %q token %q", host, token)
	}
}

func TestReadConfig_rejectsUnsupportedEnvironmentScheme(t *testing.T) {
	setupTempConfig(t)
	t.Setenv("OP_CLI_HOST", "ftp://env.example.com")
	t.Setenv("OP_CLI_TOKEN", "env-token")

	_, _, err := configuration.ReadConfig("default")
	if err == nil {
		t.Fatal("expected error for unsupported environment host scheme")
	}
	if !strings.Contains(err.Error(), "only http and https") {
		t.Errorf("error should describe supported schemes, got: %v", err)
	}
	if !strings.Contains(err.Error(), "OP_CLI_HOST") {
		t.Errorf("error should name OP_CLI_HOST as the source, got: %v", err)
	}
}

// Scheme comparison must be case-insensitive: URL schemes are defined as
// case-insensitive, and hosts pasted from other tools may be upper-cased.
func TestReadConfig_acceptsMixedCaseScheme(t *testing.T) {
	setupTempConfig(t)
	writeRaw(t, "[default]\nhost = HTTPS://file.example.com\ntoken = file-token\n")

	host, token, err := configuration.ReadConfig("default")
	if err != nil {
		t.Fatalf("ReadConfig with mixed-case scheme: %v", err)
	}
	if host != "HTTPS://file.example.com" || token != "file-token" {
		t.Errorf("got host %q token %q", host, token)
	}
}

func TestReadConfig_rejectsUnsupportedProfileScheme(t *testing.T) {
	setupTempConfig(t)
	writeRaw(t, "[default]\nhost = ftp://file.example.com\ntoken = file-token\n")

	_, _, err := configuration.ReadConfig("default")
	if err == nil {
		t.Fatal("expected error for unsupported profile host scheme")
	}
	if !strings.Contains(err.Error(), "only http and https") {
		t.Errorf("error should describe supported schemes, got: %v", err)
	}
}
