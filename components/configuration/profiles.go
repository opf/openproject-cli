package configuration

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
)

const (
	EnvProfile     = "OP_CLI_PROFILE"
	DefaultProfile = "default"
)

// Profile holds credentials for a named OpenProject instance.
type Profile struct {
	Name  string
	Host  string
	Token string
}

var invalidChars = regexp.MustCompile(`[^a-zA-Z0-9_]`)
var multiHyphen = regexp.MustCompile(`-{2,}`)

// SanitizeProfileName replaces invalid characters with hyphens, collapses
// consecutive hyphens, strips leading/trailing hyphens, and falls back to
// "default" when the result is empty.
func SanitizeProfileName(name string) string {
	result := invalidChars.ReplaceAllString(name, "-")
	result = multiHyphen.ReplaceAllString(result, "-")
	result = strings.Trim(result, "-")
	if result == "" {
		return DefaultProfile
	}
	return result
}

// ValidateProfileName returns an error when name is not already in sanitized
// form (only letters, digits, - and _, no leading/trailing hyphens, non-empty).
func ValidateProfileName(name string) error {
	if name == "" {
		return errors.Custom("profile name cannot be empty")
	}
	if SanitizeProfileName(name) != name {
		return errors.Custom(fmt.Sprintf(
			"invalid profile name %q: only letters, numbers, - and _ are allowed (no leading/trailing hyphens)",
			name,
		))
	}
	return nil
}

// ReadConfig returns host and token for profile.
// OP_CLI_HOST and OP_CLI_TOKEN always take precedence over the file.
func ReadConfig(profile string) (host, token string, err error) {
	if err = ensureConfigDir(); err != nil {
		return "", "", err
	}
	if ok, h, t := readEnvironment(); ok {
		return h, t, nil
	}
	return readConfigForProfile(profile)
}

// WriteConfigForProfile writes or updates profile in the config file.
func WriteConfigForProfile(profile, host, token string) error {
	if err := ensureConfigDir(); err != nil {
		return err
	}
	return writeProfile(profile, host, token)
}

// DeleteProfile removes profile from the config file, including any
// duplicate sections with the same name. It returns an error when the
// profile does not exist, so callers do not report success for a no-op.
func DeleteProfile(profile string) error {
	if err := ensureConfigDir(); err != nil {
		return err
	}
	return deleteProfile(profile)
}

// ConfigFilePath returns the path of the config file. It honours
// $XDG_CONFIG_HOME/$HOME, so it is absolute only when those are. Exposed so
// callers can name the file in user-facing messages.
func ConfigFilePath() string {
	return configFile()
}

// InsecureConfigPermissions reports whether the config file exists with
// permissions that let group or other users access it. The file stores API
// tokens (and is written with mode 0600), so callers should warn the user when
// this returns true. The second value is the file's permission bits. Unix
// permission bits are not meaningful on Windows, so it always reports false
// there; a missing or unreadable file is likewise not reported as insecure.
func InsecureConfigPermissions() (insecure bool, mode os.FileMode) {
	if runtime.GOOS == "windows" {
		return false, 0
	}
	info, err := os.Stat(configFile())
	if err != nil {
		return false, 0
	}
	mode = info.Mode().Perm()
	return mode&0077 != 0, mode
}

// AllProfiles returns every profile stored in the config file.
func AllProfiles() ([]*Profile, error) {
	if err := ensureConfigDir(); err != nil {
		return nil, err
	}
	return readAllProfiles()
}

// --- internal helpers --------------------------------------------------------

type iniSection struct {
	name string
	kv   map[string]string
}

type iniFile struct {
	sections []*iniSection
	index    map[string]int
}

func newIniFile() *iniFile {
	return &iniFile{index: make(map[string]int)}
}

func parseIni(data []byte) *iniFile {
	f := newIniFile()
	var current *iniSection

	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			name := line[1 : len(line)-1]
			s := &iniSection{name: name, kv: make(map[string]string)}
			f.index[name] = len(f.sections)
			f.sections = append(f.sections, s)
			current = s
			continue
		}
		if current != nil {
			if idx := strings.Index(line, "="); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				val := strings.TrimSpace(line[idx+1:])
				current.kv[key] = val
			}
		}
	}
	return f
}

func (f *iniFile) get(section, key string) (string, bool) {
	idx, ok := f.index[section]
	if !ok {
		return "", false
	}
	v, ok := f.sections[idx].kv[key]
	return v, ok
}

func (f *iniFile) hasSection(section string) bool {
	_, ok := f.index[section]
	return ok
}

func (f *iniFile) set(section, key, val string) {
	idx, ok := f.index[section]
	if !ok {
		s := &iniSection{name: section, kv: make(map[string]string)}
		f.index[section] = len(f.sections)
		f.sections = append(f.sections, s)
		idx = len(f.sections) - 1
	}
	f.sections[idx].kv[key] = val
}

// delete removes every section named section (duplicate sections can exist in
// hand-edited files; leaving one behind would silently keep its credentials
// active). Reports whether anything was removed.
func (f *iniFile) delete(section string) bool {
	kept := f.sections[:0]
	removed := false
	for _, s := range f.sections {
		if s.name == section {
			removed = true
			continue
		}
		kept = append(kept, s)
	}
	if !removed {
		return false
	}
	f.sections = kept
	f.index = make(map[string]int, len(f.sections))
	for i, s := range f.sections {
		f.index[s.name] = i
	}
	return true
}

func (f *iniFile) marshal() []byte {
	var sb strings.Builder
	for i, s := range f.sections {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("[%s]\n", s.name))
		// host and token first for stable, readable ordering.
		written := make(map[string]bool)
		for _, key := range []string{"host", "token"} {
			if v, ok := s.kv[key]; ok {
				sb.WriteString(fmt.Sprintf("%s = %s\n", key, v))
				written[key] = true
			}
		}
		// Preserve any other keys (sorted for deterministic output) so data the
		// CLI does not recognise is not silently dropped on rewrite.
		rest := make([]string, 0, len(s.kv))
		for k := range s.kv {
			if !written[k] {
				rest = append(rest, k)
			}
		}
		sort.Strings(rest)
		for _, k := range rest {
			sb.WriteString(fmt.Sprintf("%s = %s\n", k, s.kv[k]))
		}
	}
	return []byte(sb.String())
}

// looksLikeHost reports whether s parses as an absolute URL with a scheme and
// host. Used to distinguish a genuine old-format "host token" line from a
// corrupt file, so garbage is not silently migrated into bogus credentials.
func looksLikeHost(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// readOrMigrate reads the config file and migrates it from the old
// single-line "host token" format when needed. Returns (ini, migrated, error).
// A non-empty file that is neither valid INI nor a well-formed old-format line
// is reported as corrupt rather than silently yielding empty/bogus credentials.
func readOrMigrate(data []byte) (*iniFile, bool, error) {
	content := strings.TrimSpace(string(data))
	if content == "" {
		return newIniFile(), false, nil
	}

	// Old format: no section headers. Check prefix rather than Contains so that
	// IPv6 host URLs (e.g. http://[::1]) are not mis-detected as INI files.
	if !strings.HasPrefix(content, "[") {
		clean := common.SanitizeLineBreaks(content)
		parts := strings.SplitN(clean, " ", 2)
		// Require the first field to look like a real host URL, so a corrupt
		// file is not migrated into bogus credentials just because it happens
		// to contain a space.
		if len(parts) == 2 && looksLikeHost(parts[0]) && parts[1] != "" {
			f := newIniFile()
			f.set(DefaultProfile, "host", parts[0])
			f.set(DefaultProfile, "token", parts[1])
			return f, true, nil
		}
		return nil, false, invalidConfigError()
	}

	f := parseIni(data)
	if len(f.sections) == 0 {
		return nil, false, invalidConfigError()
	}
	return f, false, nil
}

func invalidConfigError() error {
	return errors.Custom(fmt.Sprintf(
		"invalid config file at %s. Please remove the file and run `op login` again.",
		configFile(),
	))
}

// writeConfigFile atomically replaces the config file: it writes a temporary
// file in the same directory (created with mode 0600) and renames it over the
// target. An interrupted write cannot leave a truncated config, and a
// permissive mode on the old file (e.g. 0644 from earlier CLI versions) is
// not carried over.
func writeConfigFile(data []byte) error {
	target := configFile()
	tmp, err := os.CreateTemp(filepath.Dir(target), "config-*.tmp")
	if err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmp.Name())
		return err
	}
	if err := os.Rename(tmp.Name(), target); err != nil {
		_ = os.Remove(tmp.Name())
		return err
	}
	return nil
}

func readOrMigrateFile() (*iniFile, error) {
	data, err := os.ReadFile(configFile())
	if os.IsNotExist(err) {
		return newIniFile(), nil
	}
	if err != nil {
		return nil, err
	}

	f, migrated, err := readOrMigrate(data)
	if err != nil {
		return nil, err
	}
	if migrated {
		// Best effort: a config that cannot be rewritten (e.g. read-only
		// location) must not block reading the credentials it holds.
		if err := writeConfigFile(f.marshal()); err != nil {
			printer.Warning(fmt.Sprintf(
				"could not rewrite config file %s in the new format: %s", configFile(), err,
			))
		}
	}
	return f, nil
}

func readConfigForProfile(profile string) (host, token string, err error) {
	f, err := readOrMigrateFile()
	if err != nil {
		return "", "", err
	}
	host, _ = f.get(profile, "host")
	token, _ = f.get(profile, "token")
	// An absent profile is a normal "not logged in" state. But a section that
	// exists yet is missing host or token is corrupt (e.g. a malformed line was
	// dropped during parsing) and must be reported rather than handed back as
	// empty credentials.
	if f.hasSection(profile) && (host == "" || token == "") {
		return "", "", invalidConfigError()
	}
	// Reject junk hosts here with a clear message; url.Parse alone accepts
	// almost anything as a relative URL, which would surface later as a
	// confusing request failure.
	if host != "" && !looksLikeHost(host) {
		return "", "", errors.Custom(fmt.Sprintf(
			"profile %q has an invalid host %q: expected an absolute URL like https://example.openproject.com",
			profile, host,
		))
	}
	return host, token, nil
}

func readAllProfiles() ([]*Profile, error) {
	f, err := readOrMigrateFile()
	if err != nil {
		return nil, err
	}
	profiles := make([]*Profile, 0, len(f.sections))
	for _, s := range f.sections {
		profiles = append(profiles, &Profile{
			Name:  s.name,
			Host:  s.kv["host"],
			Token: s.kv["token"],
		})
	}
	return profiles, nil
}

func writeProfile(profile, host, token string) error {
	f, err := readOrMigrateFile()
	if err != nil {
		return err
	}
	f.set(profile, "host", host)
	f.set(profile, "token", token)
	return writeConfigFile(f.marshal())
}

func profileNotFoundError(profile string) error {
	return errors.Custom(fmt.Sprintf("profile %q not found", profile))
}

func deleteProfile(profile string) error {
	data, err := os.ReadFile(configFile())
	if os.IsNotExist(err) {
		return profileNotFoundError(profile)
	}
	if err != nil {
		return err
	}
	f, _, err := readOrMigrate(data)
	if err != nil {
		return err
	}
	if !f.delete(profile) {
		return profileNotFoundError(profile)
	}
	return writeConfigFile(f.marshal())
}
