// Package config resolves the settings of the CLI. For every key the first
// value found wins: command line flag, environment variable, config file,
// default.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	yaml "go.yaml.in/yaml/v3"
)

// envPrefix is prepended to a key to get its environment variable, e.g.
// log-level is read from SUPERVISOR_LOG_LEVEL
const envPrefix = "SUPERVISOR_"

var (
	flags    = map[string]*pflag.Flag{}
	extraEnv = map[string][]string{}
	file     = map[string]string{}
	defaults = map[string]string{}
)

// BindFlag makes a command line flag available under its name
func BindFlag(flag *pflag.Flag) {
	flags[flag.Name] = flag
}

// BindEnv adds environment variables to read a key from, on top of the one
// derived from the name of the key
func BindEnv(key string, names ...string) {
	extraEnv[key] = append(extraEnv[key], names...)
}

// SetDefault sets the value of a key which is not set anywhere else
func SetDefault(key, value string) {
	defaults[key] = value
}

// ReadFile loads the keys of a YAML (or JSON) config file
func ReadFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var values map[string]any
	if err := yaml.Unmarshal(content, &values); err != nil {
		return fmt.Errorf("parsing %s: %w", path, err)
	}

	for key, value := range values {
		if value != nil {
			file[strings.ToLower(key)] = fmt.Sprint(value)
		}
	}
	return nil
}

// GetString returns the value of a key, or an empty string if it is not set
func GetString(key string) string {
	flag := flags[key]
	if flag != nil && flag.Changed {
		return flag.Value.String()
	}

	names := append([]string{envPrefix + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))}, extraEnv[key]...)
	for _, name := range names {
		if value := os.Getenv(name); value != "" {
			return value
		}
	}

	if value, ok := file[key]; ok {
		return value
	}
	if value, ok := defaults[key]; ok {
		return value
	}
	if flag != nil {
		return flag.DefValue
	}
	return ""
}

// GetBool returns the value of a key as boolean, false if it is not set or
// not a boolean
func GetBool(key string) bool {
	value, _ := strconv.ParseBool(GetString(key))
	return value
}
