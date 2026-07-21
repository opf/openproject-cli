package configuration

import (
	"os"
	"path/filepath"
)

const (
	envHost        = "OP_CLI_HOST"
	envToken       = "OP_CLI_TOKEN"
	configDirName  = "openproject"
	configFileName = "config"
)

func readEnvironment() (ok bool, host, token string) {
	host, hasHost := os.LookupEnv(envHost)
	token, hasToken := os.LookupEnv(envToken)
	ok = hasHost && hasToken
	return
}

// HasEnvironmentConfig reports whether OP_CLI_HOST and OP_CLI_TOKEN are both
// set, in which case they override every stored profile.
func HasEnvironmentConfig() bool {
	ok, _, _ := readEnvironment()
	return ok
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
	return filepath.Join(homeDir(), ".config", configDirName)
}

func homeDir() string {
	if home, ok := os.LookupEnv("HOME"); ok {
		return home
	}
	return filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"))
}
