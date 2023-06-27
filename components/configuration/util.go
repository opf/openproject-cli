package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const configDirName = "openproject"
const configFileName = "config"

func WriteConfigFile(host, token string) error {
	err := ensureConfigFile()
	if err != nil {
		return err
	}

	bytes := []byte(fmt.Sprintf("%s %s", host, token))
	return os.WriteFile(configFile(), bytes, 0644)
}

func ReadConfigFile() (host, token string, err error) {
	err = ensureConfigFile()
	if err != nil {
		return "", "", err
	}

	file, err := os.ReadFile(configFile())
	if err != nil {
		return "", "", err
	}

	t := strings.Replace(string(file), "\n", "", -1)
	parts := strings.Split(t, " ")
	if len(parts) < 2 {
		return parts[0], "", nil
	}
	return parts[0], parts[1], nil
}

func ensureConfigFile() error {
	if _, err := os.Stat(configFileDir()); os.IsNotExist(err) {
		err = os.MkdirAll(configFileDir(), 0700)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(configFile()); os.IsNotExist(err) {
		file, err := os.Create(configFile())
		if err != nil {
			return err
		}
		err = os.Chmod(configFile(), 0644)
		if err != nil {
			return err
		}
		defer file.Close()
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
