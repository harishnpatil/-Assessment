package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


func newTestApp(enableMetrics bool) *App {
	return &App{
		Config: &Config{
			AppHost:        "localhost",
			AppPort:        8080,
			LogLevel:       "info",
			AllowedOrigins: []string{"https://example.com"},
			EnableMetrics:  enableMetrics,
		},
		Metrics: &Metrics{},
	}
}


func TestHealthHandler_Returns200(t *testing.T) {
	app := newTestApp(false)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.HealthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status: got %d, want 200", rr.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf(`body["status"]: got %q, want "ok"`, body["status"])
	}
}

func TestHealthHandler_IncrementsCounter(t *testing.T) {
	app := newTestApp(false)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rr := httptest.NewRecorder()
		app.HealthHandler(rr, req)
	}
	h, _, _, _ := app.Metrics.Snapshot()
	if h != 3 {
		t.Errorf("health counter: got %d, want 3", h)
	}
}



func TestConfigHandler_Returns200WithConfig(t *testing.T) {
	app := newTestApp(true)
	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	rr := httptest.NewRecorder()

	app.ConfigHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status: got %d, want 200", rr.Code)
	}

	var resp configResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.AppHost != "localhost" {
		t.Errorf("app_host: got %q, want %q", resp.AppHost, "localhost")
	}
	if resp.AppPort != 8080 {
		t.Errorf("app_port: got %d, want 8080", resp.AppPort)
	}
	if !resp.EnableMetrics {
		t.Error("enable_metrics should be true")
	}
}



func TestEchoHandler_ValidMessage(t *testing.T) {
	app := newTestApp(false)
	body := `{"message":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.EchoHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status: got %d, want 200", rr.Code)
	}

	var resp echoResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Message != "hello" {
		t.Errorf("message: got %q, want %q", resp.Message, "hello")
	}
	if resp.Timestamp == "" {
		t.Error("timestamp should not be empty")
	}
}

func TestEchoHandler_EmptyMessage(t *testing.T) {
	app := newTestApp(false)
	body := `{"message":""}`
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.EchoHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", rr.Code)
	}
}

func TestEchoHandler_MissingMessageField(t *testing.T) {
	app := newTestApp(false)
	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.EchoHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", rr.Code)
	}
}

func TestEchoHandler_InvalidJSON(t *testing.T) {
	app := newTestApp(false)
	body := `not json`
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.EchoHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", rr.Code)
	}
}

func TestEchoHandler_WhitespaceOnlyMessage(t *testing.T) {
	app := newTestApp(false)
	body := `{"message":"   "}`
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.EchoHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400 for whitespace-only message", rr.Code)
	}
}



func TestMetricsHandler_EnabledReturns200(t *testing.T) {
	app := newTestApp(true)

	
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rr := httptest.NewRecorder()
		app.HealthHandler(rr, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	app.MetricsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status: got %d, want 200", rr.Code)
	}

	var resp metricsResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.HealthCount != 2 {
		t.Errorf("health_count: got %d, want 2", resp.HealthCount)
	}
	if resp.MetricsCount != 1 {
		t.Errorf("metrics_count: got %d, want 1", resp.MetricsCount)
	}
}

func TestMetricsHandler_DisabledReturns404(t *testing.T) {
	app := newTestApp(false)
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	app.MetricsHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", rr.Code)
	}
}

func TestMetricsHandler_CounterIncrements(t *testing.T) {
	app := newTestApp(true)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		rr := httptest.NewRecorder()
		app.MetricsHandler(rr, req)
	}
	_, _, _, m := app.Metrics.Snapshot()
	if m != 5 {
		t.Errorf("metrics counter: got %d, want 5", m)
	}
}
