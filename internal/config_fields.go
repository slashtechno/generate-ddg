package internal

import "github.com/spf13/viper"

// Option is one choice in a multiple-choice ConfigField.
type Option struct {
	// The option to display to the user
	Display string
	// The value to set if this option is selected
	Value string
}

// ConfigField describes a single configuration value: which Viper instance
// it lives in, how to present it for editing, and (optionally) its default.
// This is the single source of truth for config keys -- both `config`
// (editing) and `config show` (reading) iterate over ConfigFields instead
// of each hardcoding the list of keys.
type ConfigField struct {
	// The key in the viper to edit
	Key string
	// The title to display to the user (defaults to Key)
	Title string
	// The description to display to the user
	Description string
	// The viper this key is stored in
	Viper *viper.Viper
	// If set, turns editing into a multiple-choice input
	Options []Option
	// If non-empty, registered as the key's default via SetDefault
	Default string
}

// ResolvedTitle returns Title if set, falling back to Key.
func (f *ConfigField) ResolvedTitle() string {
	if f.Title != "" {
		return f.Title
	}
	return f.Key
}

// Get and Set close over Viper+Key together, so callers never need to name
// a raw key string and separately remember which Viper it lives in.
func (f *ConfigField) Get() string      { return f.Viper.GetString(f.Key) }
func (f *ConfigField) Set(value string) { f.Viper.Set(f.Key, value) }

// Individual fields are named so call sites (root.go, cmd/web/main.go) refer
// to e.g. internal.Token instead of a key string paired with a Viper.
var (
	Token = &ConfigField{
		Key:         "token",
		Title:       "DuckDuckGo API token",
		Description: "Your DuckDuckGo API token. If not set now, the login process will start the first time the program is run. The token will then be stored in the secrets file.",
		Viper:       SecretViper,
	}
	JwtSecret = &ConfigField{
		Key:         "jwt-secret",
		Title:       "JWT signing secret",
		Description: "Secret used to sign and verify web UI session tokens. Generate one with `openssl rand -hex 32`. Changing this invalidates all existing sessions.",
		Viper:       SecretViper,
	}
	WebPassword = &ConfigField{
		Key:         "web-password",
		Title:       "Web UI password",
		Description: "Password required to log in to the web UI.",
		Viper:       SecretViper,
	}
	DuckAddressUsername = &ConfigField{
		Key:         "duck-address-username",
		Title:       "DuckDuckGo address username",
		Description: "Your DuckDuckGo address username. This is the part before the @duck.com in your email address.",
		Viper:       Viper,
	}
	LogLevel = &ConfigField{
		Key:         "log-level",
		Title:       "Log level",
		Description: "The minimum log level to display",
		Viper:       Viper,
		Default:     "info",
		Options: []Option{
			{Display: "Debug", Value: "debug"},
			{Display: "Info", Value: "info"},
			{Display: "Warn", Value: "warn"},
			{Display: "Error", Value: "error"},
		},
	}
)

// ConfigFields lists every field for code that needs to iterate them all
// (the `config` edit form and `config show`).
var ConfigFields = []*ConfigField{Token, JwtSecret, WebPassword, DuckAddressUsername, LogLevel}

// ApplyDefaults registers each ConfigField's Default with its Viper.
func ApplyDefaults() {
	for _, f := range ConfigFields {
		if f.Default != "" {
			f.Viper.SetDefault(f.Key, f.Default)
		}
	}
}
