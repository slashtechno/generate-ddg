/*
Copyright © 2024 Angad Behl

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"io/fs"
	"os"
	"strings"

	// "github.com/joho/godotenv"

	"github.com/slashtechno/generate-ddg/internal"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
	"github.com/slashtechno/generate-ddg/pkg/duckduckgoapi"
	"github.com/slashtechno/generate-ddg/pkg/utils"
	"github.com/spf13/cobra"

	// Load .env
	"github.com/charmbracelet/huh"
	_ "github.com/joho/godotenv/autoload"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "generate-ddg",
	Short: "Generate DuckDuckGo email addresses from the command line",
	Run: func(cmd *cobra.Command, args []string) {
		err := duckduckgoapi.InitiateLogin(internal.Viper.GetString("duck-address-username"))
		if err != nil {
			log.Fatal("Failed to initiate login", "error", err)
		}
		var otp string
		huh.NewInput().
			Title("One-time passphrase").
			Value(&otp).
			Run()
		err = duckduckgoapi.LoginWithOtp(internal.Viper.GetString("duck-address-username"), otp)
		if err != nil {
			log.Fatal("Failed to login with OTP", "error", err)
		}
		log.Info("Successfully logged in")
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		utils.SetupLogger(internal.Viper.GetString("log-level"))
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/generate-ddg/config.yaml)")

	rootCmd.PersistentFlags().String("log-level", "", "log level")
	internal.Viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	internal.Viper.SetDefault("log-level", "info")

	rootCmd.PersistentFlags().StringP("duck-address-username", "d", "", "The username in your DuckDuckGo email address (<username>@duck.com)")
	internal.Viper.BindPFlag("duck-address-username", rootCmd.PersistentFlags().Lookup("duck-address=username"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// If the user specifies a config file, use that
	// Otherwise use $XDG_CONFIG_HOME/generate-ddg/config.yaml
	if cfgFile != "" {
		// Use config file from the flag.
		internal.Viper.SetConfigFile(cfgFile)
	} else {
		configPath, err := xdg.ConfigFile("generate-ddg/config.yaml")
		if err != nil {
			log.Fatal("Failed to get config file path:", "error", err)
		}
		internal.Viper.SetConfigFile(configPath)
	}

	// Read in environment variables that match
	internal.Viper.AutomaticEnv()
	internal.Viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	usedCfgFile := internal.Viper.ConfigFileUsed()

	// If a config file is found, read it in.
	if err := internal.Viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", internal.Viper.ConfigFileUsed())
	} else {
		if _, ok := err.(*fs.PathError); ok {
			log.Debug("Configuration file not found, creating a new one", "file", usedCfgFile)
			if err := internal.Viper.WriteConfigAs(usedCfgFile); err != nil {
				log.Fatal("Failed to write configuration file:", "error", err)
			}
			log.Fatal("Failed to read config file. Created a config file with default values. Please edit the file and run the command again.", "path", usedCfgFile)
		}

		log.Fatal("Failed to read configurationfile:", "error", err)
	}
}
