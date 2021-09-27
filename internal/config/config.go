package config

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ConfigDir string
)

// Looks for where configuations are set, if not found creates one
func InitConfigDir(flagConfigDir string) (err error) {
	//attempt to resolve home directory
	if len(flagConfigDir) > 0 {
		if _, err = os.Stat(flagConfigDir); os.IsNotExist(err) {
			ConfigDir, err = loadConfigDir()
			if err != nil {
				return err
			}
			//just a warning... can still continue
			err = errors.New("Error: " + flagConfigDir + " does not exist. Defaulting to " + ConfigDir + ".")
		} else {
			ConfigDir = flagConfigDir
			return nil
		}
	} else {
		ConfigDir, err = loadConfigDir()
		if err != nil {
			return err
		}
	}

	//attempt to create home, does nothing if exists
	e := os.MkdirAll(ConfigDir, os.ModePerm)
	if e != nil {
		return errors.New("Error creating configuration directory: " + e.Error())
	}

	return err
}

func loadConfigDir() (path string, err error) {
	//attempt loading from environmental variables
	home := os.Getenv("STRANGELET_CONFIG_HOME")
	if home == "" {
		home, err = os.UserConfigDir()
		if err != nil {
			return home, errors.New("Error finding your home directory\nCan't load config files: " + err.Error())
		}
		home = filepath.Join(home, "strangelet")
	}
	return home, err
}
