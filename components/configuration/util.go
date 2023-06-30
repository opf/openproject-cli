package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opf/openproject-cli/components/errors"
)

const (
	envHost        = "OP_CLI_HOST"
	envToken       = "OP_CLI_TOKEN"
	configDirName  = "openproject"
	configFileName = "config"
)

func WriteConfigFile(host, token string) error {
	err := ensureConfigDir()
	if err != nil {
		return err
	}

	bytes := []byte(fmt.Sprintf("%s %s", host, token))
	return os.WriteFile(configFile(), bytes, 0644)
}

func ReadConfig() (host, token string, err error) {
	err = ensureConfigDir()
	if err != nil {
		return "", "", err
	}

	ok, h, t := readEnvironment()
	if ok {
		return h, t, nil
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

func readEnvironment() (ok bool, host, token string) {
	host, hasHost := os.LookupEnv(envHost)
	token, hasToken := os.LookupEnv(envToken)
	ok = hasHost && hasToken

	return
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
