package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/rkoesters/xdg/basedir"
	"github.com/spf13/viper"
)

var configDir string
var configPath string

func Init() error {
	// Configuration location
	if runtime.GOOS != "windows" {
		configDir = filepath.Join(basedir.ConfigHome, "bougie")
	}

	configPath = filepath.Join(configDir, "config.toml")

	// Default configuration
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err == nil {
		_, err = f.Write(defaultConfiguration)
		if err != nil {
			f.Close()
			return nil
		}
		f.Close()
	}

	// Main configuration
	viper.SetDefault("general.home", "gopher://floodgap.com")

	viper.SetConfigFile(configPath)
	viper.SetConfigType("toml")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}
