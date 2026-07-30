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

// fieldEdit pairs a ConfigField with the huh form's live edit buffer for it,
// pre-seeded from the field's current value so unedited fields round-trip
// instead of being wiped to "".
type fieldEdit struct {
	Field *internal.ConfigField
	Value string
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Edit the configuration file",
	Long:  `Interactively edit the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := EditKeys(internal.ConfigFields)
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

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the configuration values, including secrets",
	Long:  `Show the configuration values, including secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		args2 := []interface{}{"config-file-used", internal.Viper.ConfigFileUsed(), "secret-config-file-used", internal.SecretViper.ConfigFileUsed()}
		for _, f := range internal.ConfigFields {
			args2 = append(args2, f.Key, f.Viper.GetString(f.Key))
		}
		log.Info("Configuration values read successfully", args2...)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(showConfigCmd)
}

func EditKeys(fields []*internal.ConfigField) error {
	var inputs []huh.Field
	edits := make([]*fieldEdit, 0, len(fields))
	for _, field := range fields {
		if field.Key == "" {
			return errors.New("key is empty")
		}
		if field.Viper == nil {
			log.Warn("Viper to edit is nil", "key", field.Key)
		}

		// Seed from the current value so unedited fields round-trip
		// instead of being blanked out.
		edit := &fieldEdit{Field: field, Value: field.Viper.GetString(field.Key)}
		edits = append(edits, edit)
		if field.Options != nil {
			inputs = append(inputs, GetSelectStringInput(edit))
		} else {
			inputs = append(inputs, GetInputForKey(edit))
		}
	}
	fmt.Printf("Existing values are pre-filled; clear a field to set it to an empty string\n")
	form := huh.NewForm(huh.NewGroup(inputs...))
	err := form.Run()
	if err != nil {
		// if err.Error() == "user aborted" {
		// 	return errors.New("user aborted")
		// }
		return err
	}

	for _, edit := range edits {
		log.Debug("Setting key", "key", edit.Field.Key, "value", edit.Value)
		edit.Field.Viper.Set(edit.Field.Key, edit.Value)
	}

	return nil
}

// titledField sets the title and, if present, description on any huh field
// builder whose fluent methods return its own concrete type.
func titledField[T interface {
	Title(string) T
	Description(string) T
}](field T, config *internal.ConfigField) T {
	field = field.Title(config.ResolvedTitle())
	if config.Description != "" {
		field = field.Description(config.Description)
	}
	return field
}

func GetSelectStringInput(edit *fieldEdit) *huh.Select[string] {
	var options []huh.Option[string]
	for _, option := range edit.Field.Options {
		options = append(options, huh.NewOption(option.Display, option.Value))
	}
	huhSelect := huh.NewSelect[string]().Options(options...).Value(&edit.Value)
	return titledField(huhSelect, edit.Field)
}

func GetInputForKey(edit *fieldEdit) *huh.Input {
	huhInput := huh.NewInput().Value(&edit.Value)
	return titledField(huhInput, edit.Field)
}
