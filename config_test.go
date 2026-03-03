package main

import (
	"os"
	"testing"
)


func setEnv(t *testing.T, pairs map[string]string) {
	t.Helper()
	for k, v := range pairs {
		t.Setenv(k, v)
	}
}


var validEnv = map[string]string{
	"APP_HOST":  "0.0.0.0",
	"APP_PORT":  "8080",
	"LOG_LEVEL": "info",
}

func TestLoadConfig_ValidMinimal(t *testing.T) {
	setEnv(t, validEnv)
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.AppHost != "0.0.0.0" {
		t.Errorf("AppHost: got %q, want %q", cfg.AppHost, "0.0.0.0")
	}
	if cfg.AppPort != 8080 {
		t.Errorf("AppPort: got %d, want 8080", cfg.AppPort)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel: got %q, want %q", cfg.LogLevel, "info")
	}
	if cfg.EnableMetrics {
		t.Error("EnableMetrics should default to false")
	}
	if len(cfg.AllowedOrigins) != 0 {
		t.Errorf("AllowedOrigins should be empty, got %v", cfg.AllowedOrigins)
	}
}

func TestLoadConfig_MissingAppHost(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_PORT":  "8080",
		"LOG_LEVEL": "info",
	})
	os.Unsetenv("APP_HOST")
	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for missing APP_HOST, got nil")
	}
}

func TestLoadConfig_MissingAppPort(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_HOST":  "0.0.0.0",
		"LOG_LEVEL": "info",
	})
	os.Unsetenv("APP_PORT")
	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for missing APP_PORT, got nil")
	}
}

func TestLoadConfig_MissingLogLevel(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_HOST": "0.0.0.0",
		"APP_PORT": "8080",
	})
	os.Unsetenv("LOG_LEVEL")
	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for missing LOG_LEVEL, got nil")
	}
}

func TestLoadConfig_InvalidAppPortNonNumeric(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_HOST":  "0.0.0.0",
		"APP_PORT":  "not-a-number",
		"LOG_LEVEL": "info",
	})
	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for non-numeric APP_PORT, got nil")
	}
}

func TestLoadConfig_AppPortTooLow(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_HOST":  "0.0.0.0",
		"APP_PORT":  "0",
		"LOG_LEVEL": "info",
	})
	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for APP_PORT=0, got nil")
	}
}

func TestLoadConfig_AppPortTooHigh(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_HOST":  "0.0.0.0",
		"APP_PORT":  "99999",
		"LOG_LEVEL": "info",
	})
	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for APP_PORT=99999, got nil")
	}
}

func TestLoadConfig_InvalidLogLevel(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_HOST":  "0.0.0.0",
		"APP_PORT":  "8080",
		"LOG_LEVEL": "verbose",
	})
	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for invalid LOG_LEVEL, got nil")
	}
}

func TestLoadConfig_AllLogLevels(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		t.Run(level, func(t *testing.T) {
			setEnv(t, map[string]string{
				"APP_HOST":  "localhost",
				"APP_PORT":  "3000",
				"LOG_LEVEL": level,
			})
			cfg, err := LoadConfig()
			if err != nil {
				t.Fatalf("LOG_LEVEL=%q should be valid, got error: %v", level, err)
			}
			if cfg.LogLevel != level {
				t.Errorf("LogLevel: got %q, want %q", cfg.LogLevel, level)
			}
		})
	}
}

func TestLoadConfig_EnableMetricsTrue(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_HOST":       "0.0.0.0",
		"APP_PORT":       "8080",
		"LOG_LEVEL":      "debug",
		"ENABLE_METRICS": "true",
	})
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.EnableMetrics {
		t.Error("EnableMetrics should be true")
	}
}

func TestLoadConfig_EnableMetricsInvalid(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_HOST":       "0.0.0.0",
		"APP_PORT":       "8080",
		"LOG_LEVEL":      "info",
		"ENABLE_METRICS": "yes",
	})
	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for invalid ENABLE_METRICS value, got nil")
	}
}

func TestLoadConfig_AllowedOrigins(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_HOST":        "0.0.0.0",
		"APP_PORT":        "8080",
		"LOG_LEVEL":       "info",
		"ALLOWED_ORIGINS": "https://a.com,https://b.com, https://c.com ",
	})
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.AllowedOrigins) != 3 {
		t.Errorf("expected 3 origins, got %d: %v", len(cfg.AllowedOrigins), cfg.AllowedOrigins)
	}
	if cfg.AllowedOrigins[0] != "https://a.com" {
		t.Errorf("first origin: got %q", cfg.AllowedOrigins[0])
	}
}
