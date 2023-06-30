package configuration

import (
	"fmt"
	"github.com/opf/openproject-cli/components/errors"
	"os"
	"path/filepath"
	"strings"
)

const configDirName = "openproject"
const configFileName = "config"

func WriteConfigFile(host, token string) error {
	err := ensureConfigDir()
	if err != nil {
		return err
	}

	bytes := []byte(fmt.Sprintf("%s %s", host, token))
	return os.WriteFile(configFile(), bytes, 0644)
}

func ReadConfigFile() (host, token string, err error) {
	err = ensureConfigDir()
	if err != nil {
		return "", "", err
	}

	file, err := os.ReadFile(configFile())
	if os.IsNotExist(err) {
		// Empty config file is no error,
		// user just has to run login command first
		return "", "", nil
	}

	sanitized := strings.Replace(string(file), "\n", "", -1)
	parts := strings.Split(sanitized, " ")
	if len(parts) != 2 {
		return "", "", errors.Custom(fmt.Sprintf("Invalid config file at %s. Please remove the file and run `op login` again.", configFile()))
	}

	return parts[0], parts[1], nil
}

func ensureConfigDir() error {
	if _, err := os.Stat(configFileDir()); os.IsNotExist(err) {
		err = os.MkdirAll(configFileDir(), 0700)
		if err != nil {
			return err
		}
	}

	return nil
}

func configFile() string {
	return filepath.Join(configFileDir(), configFileName)
}

func configFileDir() string {
	xdgConfigDir, present := os.LookupEnv("XDG_CONFIG_HOME")
	if present {
		return filepath.Join(xdgConfigDir, configDirName)
	}

	return filepath.Join(os.Getenv("HOME"), ".config", configDirName)
}
