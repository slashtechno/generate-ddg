package utils

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// LoadConfig loads the configuration from the config file and environment variables. If .env fails to load, it will log a warning (if charmLogger is not nil).
// The config file is read from the userPassedConfigPath if it is not empty, otherwise it is read from the XDG_CONFIG_HOME/<configRelativeToXdgConfig> path.
// If the config file is not found, a new one is created with default values and an error is returned.
// If the config file is found, it is read in and the path is logged (if charmLogger is not nil).
// If noFile is true, the config file is not read in and only environment variables are used.
func LoadConfig(passedViper *viper.Viper, userPassedConfigPath, configRelativeToXdgConfig string, charmLogger *log.Logger, noFile bool) error {

	err := godotenv.Load()
	if err != nil {
		if charmLogger != nil {
			charmLogger.Warn("Failed to load .env file", "error", err)
		}
	}
	// Read in environment variables that match
	passedViper.AutomaticEnv()
	passedViper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if noFile {
		return nil
	}

	// If the user specifies a config file, use that
	// Otherwise use $XDG_CONFIG_HOME/<configRelativeToXdgConfig>
	if userPassedConfigPath != "" {
		// Use config file from the flag.
		passedViper.SetConfigFile(userPassedConfigPath)
	} else {
		configPath, err := xdg.ConfigFile(configRelativeToXdgConfig)
		if err != nil {
			// https://pkg.go.dev/fmt#Errorf
			return fmt.Errorf("failed to get config file path relative to XDG_CONFIG_HOME: %w", err)
		}
		passedViper.SetConfigFile(configPath)
	}

	usedCfgFile := passedViper.ConfigFileUsed()

	// If a config file is found, read it in.
	if err := passedViper.ReadInConfig(); err == nil {
		if charmLogger != nil {
			// TODO: Allow "config file" to be replaced with a parameter
			// That way it can be called a secrets file or something else, depending on the use case
			charmLogger.Info("Loaded file", "path", passedViper.ConfigFileUsed())
		}
	} else {
		if _, ok := err.(*fs.PathError); ok {
			if charmLogger != nil {
				log.Debug("Configuration file not found, creating a new one", "file", usedCfgFile)
			}
			if err := passedViper.WriteConfigAs(usedCfgFile); err != nil {
				// if charmLogger != nil {
				// log.Error("Failed to write configuration file:", "error", err)
				// }

				return fmt.Errorf("failed to write configuration file: %w", err)

			}
			return fmt.Errorf("created new configuration file with default values; please edit the file and run the command again: %s", usedCfgFile)
		}

		return fmt.Errorf("failed to read configuration file: %w", err)
	}
	return nil
}
