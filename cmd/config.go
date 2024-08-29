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

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/fatih/color"
	"github.com/slashtechno/generate-ddg/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Edit the configuration file",
	Long:  `Interactively edit the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {

		keysToEdit := []KeyToEdit{
			{key: "token", title: "DuckDuckGo API token", description: "Your DuckDuckGo API token. If not set now, the login process will start the first time the program is run. The token will then be stored in the secrets file.", viperToEdit: internal.SecretViper},
			{key: "duck-address-username", title: "DuckDuckGo address username", description: "Your DuckDuckGo address username. This is the part before the @duck.com in your email address.", viperToEdit: internal.Viper},
			{key: "log-level", title: "Log level", description: "The log level to use. Possible values are debug, info, warn, and error.", viperToEdit: internal.Viper},
		}

		fmt.Printf("To skip editing a key, press %s\nIf you enter a blank value, the key will be set to an empty string\n", color.YellowString("Ctrl+C"))
		err := EditKeys(keysToEdit)
		if err != nil {
			log.Fatal("Failed to edit keys", "error", err)
		}

		vipers := []*viper.Viper{internal.SecretViper, internal.Viper}
		for _, v := range vipers {
			err = v.WriteConfig()
			if err != nil {
				log.Fatal("Failed to write configuration file", "error", err)
			}
			log.Infof("Wrote to %s", v.ConfigFileUsed())
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

type KeyToEdit struct {
	// The key in the viper to edit
	key string
	// The title to display to the user (defaults to the key)
	title string
	// The description to display to the user
	description string
	// The viper to edit
	viperToEdit *viper.Viper
}

func EditKeys(keys []KeyToEdit) error {
	for _, key := range keys {
		value, err := GetValueForKey(key)

		if err != nil {
			if err.Error() == "user aborted" {
				log.Debug("User aborted input")
				continue
			}
			return err
		}
		key.viperToEdit.Set(key.key, value)
		log.Debug("Set key", "key", key.key, "value", key.viperToEdit.Get(key.key))
	}
	return nil
}

// Get a new value for a key via huh
func GetValueForKey(key KeyToEdit) (string, error) {
	var value, title string

	huhInput := huh.NewInput().Value(&value)

	if key.title != "" {
		log.Debug("Using title from key", "key", key.key)
		title = key.title
	} else {
		log.Debug("Using key as title", "key", key.key, "title", key.key)
		title = key.key
	}
	huhInput.Title(title)

	if key.description != "" {
		huhInput.Description(key.description)
	}

	err := huhInput.Run()
	if err != nil {
		return "", err
	}

	return value, nil
}
