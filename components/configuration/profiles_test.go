package configuration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/opf/openproject-cli/components/configuration"
)

// setupTempConfig points XDG_CONFIG_HOME at a temp dir and unsets credential
// env vars so every test starts from a clean slate.
func setupTempConfig(t *testing.T) {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	unsetEnv(t, "OP_CLI_HOST")
	unsetEnv(t, "OP_CLI_TOKEN")
}

// unsetEnv removes an env var for the duration of the test and restores the
// original value (or removes it again) in cleanup.
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	old, existed := os.LookupEnv(key)
	os.Unsetenv(key)
	t.Cleanup(func() {
		if existed {
			os.Setenv(key, old)
		} else {
			os.Unsetenv(key)
		}
	})
}

// writeRaw writes arbitrary bytes directly to the config file, bypassing the
// profile API – used to seed old-format files for migration tests.
func writeRaw(t *testing.T, content string) {
	t.Helper()
	dir := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "openproject")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// ---------- SanitizeProfileName ----------

func TestSanitizeProfileName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"default", "default"},
		{"work", "work"},
		{"my-profile", "my-profile"},
		{"my_profile", "my_profile"},
		{"MyProfile123", "MyProfile123"},
		{"my profile", "my-profile"},
		{"my@work", "my-work"},
		{"my--work", "my-work"},
		{"-mywork-", "mywork"},
		{"--", "default"},
		{"!@#$%", "default"},
		{"", "default"},
	}
	for _, c := range cases {
		got := configuration.SanitizeProfileName(c.input)
		if got != c.want {
			t.Errorf("SanitizeProfileName(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

// ---------- ValidateProfileName ----------

func TestValidateProfileName_valid(t *testing.T) {
	valid := []string{
		"default", "work", "my-profile", "my_profile",
		"MyProfile123", "a", "1", "a1", "a-b", "a_b",
	}
	for _, name := range valid {
		if err := configuration.ValidateProfileName(name); err != nil {
			t.Errorf("ValidateProfileName(%q) should be valid, got: %v", name, err)
		}
	}
}

func TestValidateProfileName_invalid(t *testing.T) {
	invalid := []string{
		"",
		"my profile",
		"my@work",
		"!@#",
		"-leading",
		"trailing-",
		"my--work",
		"--",
	}
	for _, name := range invalid {
		if err := configuration.ValidateProfileName(name); err == nil {
			t.Errorf("ValidateProfileName(%q) should be invalid, got no error", name)
		}
	}
}

// ---------- WriteConfigForProfile / ReadConfig ----------

func TestWriteAndReadConfigForProfile(t *testing.T) {
	setupTempConfig(t)

	if err := configuration.WriteConfigForProfile("default", "https://example.com", "token123"); err != nil {
		t.Fatal(err)
	}

	host, token, err := configuration.ReadConfig("default")
	if err != nil {
		t.Fatal(err)
	}
	if host != "https://example.com" {
		t.Errorf("host = %q, want %q", host, "https://example.com")
	}
	if token != "token123" {
		t.Errorf("token = %q, want %q", token, "token123")
	}
}

func TestWriteMultipleProfilesAndReadBack(t *testing.T) {
	setupTempConfig(t)

	if err := configuration.WriteConfigForProfile("default", "https://default.example.com", "tok-default"); err != nil {
		t.Fatal(err)
	}
	if err := configuration.WriteConfigForProfile("work", "https://work.example.com", "tok-work"); err != nil {
		t.Fatal(err)
	}

	host, token, err := configuration.ReadConfig("work")
	if err != nil {
		t.Fatal(err)
	}
	if host != "https://work.example.com" {
		t.Errorf("host = %q, want %q", host, "https://work.example.com")
	}
	if token != "tok-work" {
		t.Errorf("token = %q, want %q", token, "tok-work")
	}

	host, token, err = configuration.ReadConfig("default")
	if err != nil {
		t.Fatal(err)
	}
	if host != "https://default.example.com" {
		t.Errorf("host = %q, want %q", host, "https://default.example.com")
	}
	if token != "tok-default" {
		t.Errorf("token = %q, want %q", token, "tok-default")
	}
}

func TestReadConfig_missingProfile_returnsEmpty(t *testing.T) {
	setupTempConfig(t)

	// Write one profile but read a different one
	if err := configuration.WriteConfigForProfile("default", "https://example.com", "tok"); err != nil {
		t.Fatal(err)
	}

	host, token, err := configuration.ReadConfig("nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if host != "" || token != "" {
		t.Errorf("expected empty credentials for missing profile, got host=%q token=%q", host, token)
	}
}

func TestReadConfig_noFile_returnsEmpty(t *testing.T) {
	setupTempConfig(t)

	host, token, err := configuration.ReadConfig("default")
	if err != nil {
		t.Fatal(err)
	}
	if host != "" || token != "" {
		t.Errorf("expected empty credentials when no config file, got host=%q token=%q", host, token)
	}
}

func TestReadConfig_envVarsOverrideProfile(t *testing.T) {
	setupTempConfig(t)

	if err := configuration.WriteConfigForProfile("default", "https://file.example.com", "file-token"); err != nil {
		t.Fatal(err)
	}

	t.Setenv("OP_CLI_HOST", "https://env.example.com")
	t.Setenv("OP_CLI_TOKEN", "env-token")

	host, token, err := configuration.ReadConfig("default")
	if err != nil {
		t.Fatal(err)
	}
	if host != "https://env.example.com" {
		t.Errorf("host = %q, want env var value", host)
	}
	if token != "env-token" {
		t.Errorf("token = %q, want env var value", token)
	}
}

func TestWriteConfigForProfile_overwritesExisting(t *testing.T) {
	setupTempConfig(t)

	if err := configuration.WriteConfigForProfile("work", "https://old.example.com", "old-token"); err != nil {
		t.Fatal(err)
	}
	if err := configuration.WriteConfigForProfile("work", "https://new.example.com", "new-token"); err != nil {
		t.Fatal(err)
	}

	host, token, err := configuration.ReadConfig("work")
	if err != nil {
		t.Fatal(err)
	}
	if host != "https://new.example.com" {
		t.Errorf("host = %q, want %q", host, "https://new.example.com")
	}
	if token != "new-token" {
		t.Errorf("token = %q, want %q", token, "new-token")
	}
}

// ---------- AllProfiles ----------

func TestAllProfiles(t *testing.T) {
	setupTempConfig(t)

	if err := configuration.WriteConfigForProfile("default", "https://a.example.com", "tok-a"); err != nil {
		t.Fatal(err)
	}
	if err := configuration.WriteConfigForProfile("work", "https://b.example.com", "tok-b"); err != nil {
		t.Fatal(err)
	}

	profiles, err := configuration.AllProfiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}

	byName := make(map[string]*configuration.Profile)
	for _, p := range profiles {
		byName[p.Name] = p
	}

	if p, ok := byName["default"]; !ok {
		t.Error("missing 'default' profile")
	} else {
		if p.Host != "https://a.example.com" {
			t.Errorf("default host = %q", p.Host)
		}
	}
	if p, ok := byName["work"]; !ok {
		t.Error("missing 'work' profile")
	} else {
		if p.Host != "https://b.example.com" {
			t.Errorf("work host = %q", p.Host)
		}
	}
}

func TestAllProfiles_noFile_returnsEmpty(t *testing.T) {
	setupTempConfig(t)

	profiles, err := configuration.AllProfiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 0 {
		t.Errorf("expected 0 profiles, got %d", len(profiles))
	}
}

// ---------- DeleteProfile ----------

func TestDeleteProfile_removesProfile(t *testing.T) {
	setupTempConfig(t)

	if err := configuration.WriteConfigForProfile("default", "https://a.example.com", "tok-a"); err != nil {
		t.Fatal(err)
	}
	if err := configuration.WriteConfigForProfile("work", "https://b.example.com", "tok-b"); err != nil {
		t.Fatal(err)
	}

	if err := configuration.DeleteProfile("default"); err != nil {
		t.Fatal(err)
	}

	profiles, err := configuration.AllProfiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile after delete, got %d", len(profiles))
	}
	if profiles[0].Name != "work" {
		t.Errorf("remaining profile name = %q, want %q", profiles[0].Name, "work")
	}
}

func TestDeleteProfile_idempotent(t *testing.T) {
	setupTempConfig(t)

	// Delete on non-existent profile must not error
	if err := configuration.DeleteProfile("nonexistent"); err != nil {
		t.Errorf("DeleteProfile on missing profile should not error, got: %v", err)
	}

	// Delete on missing file must not error
	if err := configuration.DeleteProfile("default"); err != nil {
		t.Errorf("DeleteProfile with no config file should not error, got: %v", err)
	}
}

// ---------- Migration ----------

func TestMigration_oldFormatMigratedToDefault(t *testing.T) {
	setupTempConfig(t)
	writeRaw(t, "https://legacy.example.com legacytoken")

	host, token, err := configuration.ReadConfig("default")
	if err != nil {
		t.Fatal(err)
	}
	if host != "https://legacy.example.com" {
		t.Errorf("host = %q, want %q", host, "https://legacy.example.com")
	}
	if token != "legacytoken" {
		t.Errorf("token = %q, want %q", token, "legacytoken")
	}
}

func TestMigration_IPv6HostMigratedCorrectly(t *testing.T) {
	setupTempConfig(t)
	writeRaw(t, "http://[::1]:8080 mytoken")

	host, token, err := configuration.ReadConfig("default")
	if err != nil {
		t.Fatal(err)
	}
	if host != "http://[::1]:8080" {
		t.Errorf("host = %q, want %q", host, "http://[::1]:8080")
	}
	if token != "mytoken" {
		t.Errorf("token = %q, want %q", token, "mytoken")
	}
}

func TestMigration_oldFormatRewrittenAsIni(t *testing.T) {
	setupTempConfig(t)
	writeRaw(t, "https://legacy.example.com legacytoken")

	// Trigger migration by reading
	if _, _, err := configuration.ReadConfig("default"); err != nil {
		t.Fatal(err)
	}

	// Now AllProfiles should work correctly
	profiles, err := configuration.AllProfiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 1 || profiles[0].Name != "default" {
		t.Errorf("after migration: expected [default], got %v", profiles)
	}
}
