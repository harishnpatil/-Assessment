package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}

	log.Printf("starting go-service host=%s port=%d log_level=%s enable_metrics=%v",
		cfg.AppHost, cfg.AppPort, cfg.LogLevel, cfg.EnableMetrics)

	
	app := &App{
		Config:  cfg,
		Metrics: &Metrics{},
	}

	
	mux := http.NewServeMux()
	mux.HandleFunc("/health", app.HealthHandler)
	mux.HandleFunc("/config", app.ConfigHandler)
	mux.HandleFunc("/echo", app.EchoHandler)
	mux.HandleFunc("/metrics", app.MetricsHandler)

	addr := fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", addr)
		serverErr <- server.ListenAndServe()
	}()

	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	case sig := <-quit:
		log.Printf("received signal %v – shutting down", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("graceful shutdown failed: %v", err)
		}
		log.Println("server stopped cleanly")
	}
}
