package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)


type Config struct {
	AppHost        string
	AppPort        int
	LogLevel       string
	AllowedOrigins []string
	EnableMetrics  bool
}


func LoadConfig() (*Config, error) {
	cfg := &Config{}

	
	host := os.Getenv("APP_HOST")
	if host == "" {
		return nil, fmt.Errorf("APP_HOST is required but not set")
	}
	cfg.AppHost = host

	
	portStr := os.Getenv("APP_PORT")
	if portStr == "" {
		return nil, fmt.Errorf("APP_PORT is required but not set")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("APP_PORT must be a valid integer, got %q", portStr)
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("APP_PORT must be between 1 and 65535, got %d", port)
	}
	cfg.AppPort = port

	
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		return nil, fmt.Errorf("LOG_LEVEL is required but not set")
	}
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[logLevel] {
		return nil, fmt.Errorf("LOG_LEVEL must be one of [debug, info, warn, error], got %q", logLevel)
	}
	cfg.LogLevel = logLevel

	
	rawOrigins := os.Getenv("ALLOWED_ORIGINS")
	if rawOrigins != "" {
		parts := strings.Split(rawOrigins, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				cfg.AllowedOrigins = append(cfg.AllowedOrigins, trimmed)
			}
		}
	}

	
	enableMetricsStr := os.Getenv("ENABLE_METRICS")
	if enableMetricsStr != "" {
		val, err := strconv.ParseBool(enableMetricsStr)
		if err != nil {
			return nil, fmt.Errorf("ENABLE_METRICS must be a boolean (true/false), got %q", enableMetricsStr)
		}
		cfg.EnableMetrics = val
	}
	

	return cfg, nil
}
