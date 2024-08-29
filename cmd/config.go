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

		keysToEdit := []KeyToEdit{
			{key: "token", viperToEdit: internal.SecretViper},
		}

		err := EditKeys(keysToEdit)
		if err != nil {
			log.Fatal("Failed to edit keys", "error", err)
		}

		internal.SecretViper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

type KeyToEdit struct {
	key         string
	viperToEdit *viper.Viper
}

func EditKeys(keys []KeyToEdit) error {
	for _, key := range keys {
		value, err := GetValueForKey(key.key)
		if err != nil {
			return err
		}
		if value == "" {
			continue
		}

		key.viperToEdit.Set(key.key, value)
	}
	return nil
}

// Get a new value for a key via huh
func GetValueForKey(key string) (string, error) {
	var value string
	err := huh.NewInput().
		Title("Edit " + key).
		Value(&value).
		Run()
	if err != nil {
		if err.Error() == "user aborted" {
			log.Debug("User aborted input")
			return "", nil
		} else {
			return "", err
		}

	}

	return value, nil
}
