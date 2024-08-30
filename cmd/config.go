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
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
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

		keysToEdit := []*KeyToEdit{
			{Key: "token", Title: "DuckDuckGo API token", Description: "Your DuckDuckGo API token. If not set now, the login process will start the first time the program is run. The token will then be stored in the secrets file.", ViperToEdit: internal.SecretViper},
			{Key: "duck-address-username", Title: "DuckDuckGo address username", Description: "Your DuckDuckGo address username. This is the part before the @duck.com in your email address.", ViperToEdit: internal.Viper},
			{Key: "log-level", Title: "Log level", Description: "The minimum log level to display",
				ViperToEdit: internal.Viper,
				// The log level, ideally should be multiple choice
				Options: []Option{
					{Display: "Debug", Value: "debug"},
					{Display: "Info", Value: "info"},
					{Display: "Warn", Value: "warn"},
					{Display: "Error", Value: "error"},
				},
			},
		}

		// fmt.Printf("To skip editing a key, press %s\nIf you enter a blank value, the key will be set to an empty string\n", color.YellowString("Ctrl+C"))
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

type Option struct {
	// The option to display to the user
	Display string
	// The value to set if this option is selected
	Value string
}

type KeyToEdit struct {
	// The key in the viper to edit
	Key string
	// The title to display to the user (defaults to the key)
	Title string
	// The description to display to the user
	Description string
	// The viper to edit
	ViperToEdit *viper.Viper
	// If true, Options will turn into a multiple choice input
	Options []Option
	value   *string
}

func EditKeys(keys []*KeyToEdit) error {
	var inputs []huh.Field
	for _, key := range keys {
		if key.Key == "" {
			return errors.New("key is empty")
		}
		if key.ViperToEdit == nil {
			log.Warn("Viper to edit is nil", "key", key.Key)
		}
		if key.value != nil {
			log.Warn("Value is not nil; overwriting", "key", key.Key)
		}

		// new() allocates memory for the value
		key.value = new(string)
		if key.Options != nil {
			inputs = append(inputs, GetSelectStringInput(key))
		} else {
			inputs = append(inputs, GetInputForKey(key))
		}
	}
	fmt.Printf("If you enter a blank value, the key will be set to an empty string\n")
	form := huh.NewForm(huh.NewGroup(inputs...))
	err := form.Run()
	if err != nil {
		// if err.Error() == "user aborted" {
		// 	return errors.New("user aborted")
		// }
		return err
	}

	for _, key := range keys {
		if key.Key != "" {
			log.Debug("Setting key", "key", key.Key, "value", *key.value)
			key.ViperToEdit.Set(key.Key, *key.value)
		} else {
			return errors.New("key is empty")
		}
	}

	return nil
}

func GetSelectStringInput(key *KeyToEdit) *huh.Select[string] {
	var title string

	var options []huh.Option[string]
	for _, option := range key.Options {
		options = append(options, huh.NewOption(option.Display, option.Value))
	}
	huhSelect := huh.NewSelect[string]().Options(options...).Value(key.value)

	if key.Title != "" {
		log.Debug("Using title from key", "key", key.Key)
		title = key.Title
	} else {
		log.Debug("Using key as title", "key", key.Key, "title", key.Key)
		title = key.Key
	}
	huhSelect.Title(title)

	if key.Description != "" {
		huhSelect.Description(key.Description)
	}

	return huhSelect
}

func GetInputForKey(key *KeyToEdit) *huh.Input {
	var title string

	huhInput := huh.NewInput().Value(key.value)

	if key.Title != "" {
		log.Debug("Using title from key", "key", key.Key)
		title = key.Title
	} else {
		log.Debug("Using key as title", "key", key.Key, "title", key.Key)
		title = key.Key
	}
	huhInput.Title(title)

	if key.Description != "" {
		huhInput.Description(key.Description)
	}

	return huhInput

}
