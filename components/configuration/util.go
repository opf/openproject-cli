package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const configFilePath = ".config/openproject"
const configFileName = "config"

func WriteConfigFile(host, token string) error {
	if _, err := os.Stat(configFileDir()); os.IsNotExist(err) {
		_ = os.MkdirAll(configFileDir(), 0700)
	}
	
	bytes := []byte(fmt.Sprintf("%s %s", host, token))
	return os.WriteFile(configFile(), bytes, 0644)
}

func ReadConfigFile() (host, token string, err error) {
	file, err := os.ReadFile(configFile())
	if err != nil {
		return "", "", err
	}

	t := strings.Replace(string(file), "\n", "", -1)
	parts := strings.Split(t, " ")
	return parts[0], parts[1], nil
}

func configFile() string {
	return filepath.Join(os.Getenv("HOME"), configFilePath, configFileName)
}

func configFileDir() string {
	return filepath.Join(os.Getenv("HOME"), configFilePath)
}