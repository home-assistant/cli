package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
)

func TestPrecedence(t *testing.T) {
	t.Setenv("SUPERVISOR_ENDPOINT", "")
	t.Setenv("SUPERVISOR_API_TOKEN", "")
	t.Setenv("SUPERVISOR_TOKEN", "")
	t.Setenv("SUPERVISOR_RAW_JSON", "")

	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flagSet.String("endpoint", "", "")
	flagSet.Bool("raw-json", false, "")
	BindFlag(flagSet.Lookup("endpoint"))
	BindFlag(flagSet.Lookup("raw-json"))
	BindEnv("api-token", "SUPERVISOR_TOKEN")

	if got := GetString("endpoint"); got != "" {
		t.Errorf("unset: got %q", got)
	}
	if GetBool("raw-json") {
		t.Error("raw-json: got true from the flag default")
	}

	SetDefault("endpoint", "default")
	if got := GetString("endpoint"); got != "default" {
		t.Errorf("default: got %q", got)
	}

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("Endpoint: file\nraw-json: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ReadFile(path); err != nil {
		t.Fatal(err)
	}
	if got := GetString("endpoint"); got != "file" {
		t.Errorf("file: got %q", got)
	}
	if !GetBool("raw-json") {
		t.Error("raw-json: got false from the file")
	}

	t.Setenv("SUPERVISOR_ENDPOINT", "env")
	if got := GetString("endpoint"); got != "env" {
		t.Errorf("env: got %q", got)
	}

	if err := flagSet.Parse([]string{"--endpoint", "flag", "--raw-json=false"}); err != nil {
		t.Fatal(err)
	}
	if got := GetString("endpoint"); got != "flag" {
		t.Errorf("flag: got %q", got)
	}
	if GetBool("raw-json") {
		t.Error("raw-json: flag did not override the file")
	}
}

func TestBindEnv(t *testing.T) {
	t.Setenv("SUPERVISOR_API_TOKEN", "")
	t.Setenv("SUPERVISOR_TOKEN", "token")
	BindEnv("api-token", "SUPERVISOR_TOKEN")

	if got := GetString("api-token"); got != "token" {
		t.Errorf("got %q, want %q", got, "token")
	}

	t.Setenv("SUPERVISOR_API_TOKEN", "api-token")
	if got := GetString("api-token"); got != "api-token" {
		t.Errorf("got %q, want %q", got, "api-token")
	}
}

func TestReadFileErrors(t *testing.T) {
	if err := ReadFile(filepath.Join(t.TempDir(), "missing.yaml")); err == nil {
		t.Error("expected an error for a missing file")
	}

	path := filepath.Join(t.TempDir(), "broken.yaml")
	if err := os.WriteFile(path, []byte("endpoint: [\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ReadFile(path); err == nil {
		t.Error("expected an error for invalid YAML")
	}
}
