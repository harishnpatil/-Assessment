package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)


type App struct {
	Config  *Config
	Metrics *Metrics
}


func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON encode error: %v", err)
	}
}


func (a *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
	a.Metrics.IncrHealth()
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}


type configResponse struct {
	AppHost        string   `json:"app_host"`
	AppPort        int      `json:"app_port"`
	LogLevel       string   `json:"log_level"`
	AllowedOrigins []string `json:"allowed_origins"`
	EnableMetrics  bool     `json:"enable_metrics"`
}


func (a *App) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	a.Metrics.IncrConfig()

	origins := a.Config.AllowedOrigins
	if origins == nil {
		origins = []string{} 
	}

	resp := configResponse{
		AppHost:        a.Config.AppHost,
		AppPort:        a.Config.AppPort,
		LogLevel:       a.Config.LogLevel,
		AllowedOrigins: origins,
		EnableMetrics:  a.Config.EnableMetrics,
	}
	writeJSON(w, http.StatusOK, resp)
}


type echoRequest struct {
	Message string `json:"message"`
}


type echoResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}


func (a *App) EchoHandler(w http.ResponseWriter, r *http.Request) {
	a.Metrics.IncrEcho()

	var req echoRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body: " + err.Error(),
		})
		return
	}

	if strings.TrimSpace(req.Message) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "message field is required and cannot be empty",
		})
		return
	}

	writeJSON(w, http.StatusOK, echoResponse{
		Message:   req.Message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}


type metricsResponse struct {
	HealthCount  int64 `json:"health_count"`
	ConfigCount  int64 `json:"config_count"`
	EchoCount    int64 `json:"echo_count"`
	MetricsCount int64 `json:"metrics_count"`
}


func (a *App) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if !a.Config.EnableMetrics {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "metrics endpoint is disabled; set ENABLE_METRICS=true to enable",
		})
		return
	}

	a.Metrics.IncrMetrics()
	h, c, e, m := a.Metrics.Snapshot()
	writeJSON(w, http.StatusOK, metricsResponse{
		HealthCount:  h,
		ConfigCount:  c,
		EchoCount:    e,
		MetricsCount: m,
	})
}
