/*
Safe-Code Telegram Bot
Copyright (c) 2024 Engagement Labs. All rights reserved.

This code is provided for VERIFICATION and TRANSPARENCY purposes only.
You may NOT use, copy, modify, distribute, or incorporate this code.
See LICENSE file for full terms.
*/

package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Build-time variables
var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildTime = "unknown"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}
	
	// Validate token format (basic)
	if len(token) < 40 || !strings.Contains(token, ":") {
		log.Fatal("Invalid TELEGRAM_BOT_TOKEN format")
	}
	
	// Start health check server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Create a custom mux for better routing control
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/version", versionHandler)
	mux.HandleFunc("/attestation", attestationHandler)
	mux.HandleFunc("/", notFoundHandler) // Catch-all for 404s
	
	// Apply logging middleware
	server.Handler = loggingMiddleware(mux)
	
	go func() {
		log.Printf("Starting HTTP server on :8080")
		log.Printf("Available endpoints: /health, /version, /attestation")
		log.Fatal(server.ListenAndServe())
	}()
	
	// Start Telegram bot
	runTelegramBot(token)
}

