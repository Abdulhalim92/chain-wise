package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_NoConfigFile_ReturnsErrNoConfigFile(t *testing.T) {
	dir := t.TempDir()
	prev, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(prev) }()
	prevEnv := os.Getenv("CONFIG_FILE")
	_ = os.Unsetenv("CONFIG_FILE")
	defer func() {
		if prevEnv != "" {
			_ = os.Setenv("CONFIG_FILE", prevEnv)
		}
	}()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when no config file")
	}
	if !errors.Is(err, ErrNoConfigFile) {
		t.Fatalf("got %v, want ErrNoConfigFile", err)
	}
}

func TestLoad_FromFile_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	const body = `
port: "8080"
grpc_port: "9090"
env: dev
log_level: info
log_format: json
`
	if err := os.WriteFile(path, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	prev, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(prev) }()
	prevEnv := os.Getenv("CONFIG_FILE")
	_ = os.Unsetenv("CONFIG_FILE")
	defer func() {
		if prevEnv != "" {
			_ = os.Setenv("CONFIG_FILE", prevEnv)
		}
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "8080" || cfg.GrpcPort != "9090" || cfg.Env != "dev" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoad_EnvOverridesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	const body = `
port: "8080"
grpc_port: "9090"
env: dev
log_level: info
log_format: json
`
	if err := os.WriteFile(path, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	prevWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(prevWd) }()
	prevPort := os.Getenv("PORT")
	_ = os.Setenv("PORT", "9999")
	defer func() {
		if prevPort != "" {
			_ = os.Setenv("PORT", prevPort)
		} else {
			_ = os.Unsetenv("PORT")
		}
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "9999" {
		t.Fatalf("expected PORT env to override file, got port %q", cfg.Port)
	}
}

func TestLoad_MissingRequiredField_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	const body = `
grpc_port: "9090"
env: dev
log_level: info
log_format: json
`
	if err := os.WriteFile(path, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	prevWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(prevWd) }()
	prevConfigFile := os.Getenv("CONFIG_FILE")
	_ = os.Unsetenv("CONFIG_FILE")
	defer func() {
		if prevConfigFile != "" {
			_ = os.Setenv("CONFIG_FILE", prevConfigFile)
		}
	}()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when port is missing")
	}
	if err != nil && err.Error() != "config: port is required" {
		t.Fatalf("got %v", err)
	}
}

func TestLoad_InvalidPort_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	const body = `
port: "99999"
grpc_port: "9090"
env: dev
log_level: info
log_format: json
`
	if err := os.WriteFile(path, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	prevWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(prevWd) }()
	_ = os.Unsetenv("CONFIG_FILE")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when port out of range")
	}
}

func TestLoad_InvalidLogLevel_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	const body = `
port: "8080"
grpc_port: "9090"
env: dev
log_level: invalid
log_format: json
`
	if err := os.WriteFile(path, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	prevWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(prevWd) }()
	_ = os.Unsetenv("CONFIG_FILE")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when log_level invalid")
	}
}
