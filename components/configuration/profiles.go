package configuration

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/errors"
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

// DeleteProfile removes profile from the config file.
// It is idempotent: returns nil even when the profile does not exist.
func DeleteProfile(profile string) error {
	if err := ensureConfigDir(); err != nil {
		return err
	}
	return deleteProfile(profile)
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

func (f *iniFile) delete(section string) {
	idx, ok := f.index[section]
	if !ok {
		return
	}
	f.sections = append(f.sections[:idx], f.sections[idx+1:]...)
	delete(f.index, section)
	for i, s := range f.sections {
		f.index[s.name] = i
	}
}

func (f *iniFile) marshal() []byte {
	var sb strings.Builder
	for i, s := range f.sections {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("[%s]\n", s.name))
		for _, key := range []string{"host", "token"} {
			if v, ok := s.kv[key]; ok {
				sb.WriteString(fmt.Sprintf("%s = %s\n", key, v))
			}
		}
	}
	return []byte(sb.String())
}

// readOrMigrate reads the config file and migrates it from the old
// single-line "host token" format when needed. Returns (ini, migrated, error).
func readOrMigrate(data []byte) (*iniFile, bool) {
	content := strings.TrimSpace(string(data))

	// Old format: no section headers
	if !strings.Contains(content, "[") && content != "" {
		clean := common.SanitizeLineBreaks(content)
		parts := strings.SplitN(clean, " ", 2)
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			f := newIniFile()
			f.set(DefaultProfile, "host", parts[0])
			f.set(DefaultProfile, "token", parts[1])
			return f, true
		}
	}

	return parseIni(data), false
}

func readOrMigrateFile() (*iniFile, error) {
	data, err := os.ReadFile(configFile())
	if os.IsNotExist(err) {
		return newIniFile(), nil
	}
	if err != nil {
		return nil, err
	}

	f, migrated := readOrMigrate(data)
	if migrated {
		if err := os.WriteFile(configFile(), f.marshal(), 0644); err != nil {
			return nil, err
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
	return os.WriteFile(configFile(), f.marshal(), 0644)
}

func deleteProfile(profile string) error {
	data, err := os.ReadFile(configFile())
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	f, _ := readOrMigrate(data)
	f.delete(profile)
	return os.WriteFile(configFile(), f.marshal(), 0644)
}
