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
	"fmt"
	"os"

	// "github.com/joho/godotenv"

	"github.com/slashtechno/generate-ddg/internal"

	"github.com/charmbracelet/log"
	"github.com/slashtechno/generate-ddg/pkg/duckduckgoapi"
	"github.com/slashtechno/generate-ddg/pkg/utils"
	"github.com/spf13/cobra"

	"github.com/charmbracelet/huh"
)

var cfgFile string
var secretsFile string
var otp string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "generate-ddg [flags]",
	Short: "Generate DuckDuckGo email addresses from the command line",
	Run: func(cmd *cobra.Command, args []string) {

		if internal.SecretViper.GetString("token") == "" {
			if otp == "" {
				err := duckduckgoapi.InitiateLogin(internal.Viper.GetString("duck-address-username"))
				if err != nil {
					log.Fatal("Failed to initiate login", "error", err)
				}
				err = huh.NewInput().
					Title("One-time passphrase").
					Value(&otp).
					Run()
				if err != nil {
					if err.Error() == "user aborted" {
						log.Fatal("User aborted")
					} else {
						log.Fatal("Failed to get OTP input", "error", err)
					}
				}
			} else {
				log.Info("Using OTP from flag")
			}

			token, err := duckduckgoapi.LoginWithOtp(internal.Viper.GetString("duck-address-username"), otp)
			if err != nil {
				log.Fatal("Failed to login with OTP", "error", err)
			}
			log.Info("Successfully logged in with OTP")

			internal.SecretViper.Set("token", token)
			err = internal.SecretViper.WriteConfig()
			if err != nil {
				log.Fatal("Failed to write token to secrets file", "error", err)
			}
			log.Info("Token written to secrets file")
		}

		accessToken, err := duckduckgoapi.GetAccessToken(internal.SecretViper.GetString("token"))
		if err != nil {
			log.Fatal("Failed to get access token", "error", err)
		}
		log.Debug("Access token:", "token", accessToken)
		log.Info("Logged in")

		email, err := duckduckgoapi.GetEmail(accessToken)
		if err != nil {
			log.Fatal("Failed to get email", "error", err)
		}
		log.Info("Generated email", "email", fmt.Sprintf("%s@duck.com", email))
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $XDG_CONFIG_HOME/generate-ddg/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&secretsFile, "secrets", "", "Secrets file (default is $XDG_CONFIG_HOME/generate-ddg/secrets.yaml). This file will have the token written to it if it's not passed via an environment variable.")

	rootCmd.PersistentFlags().String("log-level", "", "Log level")
	internal.Viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	internal.Viper.SetDefault("log-level", "info")

	rootCmd.Flags().StringVarP(&otp, "otp", "o", "", "One-time passphrase")

	internal.Viper.SetDefault("duck-address-username", "...")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	utils.LoadConfig(
		internal.Viper,
		cfgFile,
		"generate-ddg/config.yaml",
		log.Default(),
		false,
	)
	utils.LoadConfig(
		internal.SecretViper,
		secretsFile,
		"generate-ddg/secrets.yaml",
		log.Default(),
		false,
	)
}
