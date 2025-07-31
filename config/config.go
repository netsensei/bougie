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
	viper.SetDefault("keybindings.quit", []string{"ctrl+c", "ctrl+q"})
	viper.SetDefault("keybindings.nav", "ctrl+n")
	viper.SetDefault("keybindings.view", "ctrl+v")
	viper.SetDefault("keybindings.home", "ctrl+h")
	viper.SetDefault("keybindings.reload", "ctrl+r")
	viper.SetDefault("keybindings.enter", "enter")
	viper.SetDefault("keybindings.item_forward", "tab")
	viper.SetDefault("keybindings.item_backward", "shift+tab")
	viper.SetDefault("keybindings.page_forward", "f")
	viper.SetDefault("keybindings.page_backward", "b")
	viper.SetDefault("keybindings.component_forward", "tab")

	viper.SetConfigFile(configPath)
	viper.SetConfigType("toml")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	keysInit()

	return nil
}
