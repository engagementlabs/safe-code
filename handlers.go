package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Content-Type", "text/plain")
	
	// Only allow GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	// Security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Content-Type", "application/json")
	
	// Only allow GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	versionInfo := map[string]string{
		"version":    Version,
		"commit_sha": CommitSHA,
		"build_time": BuildTime,
		"github_url": fmt.Sprintf("https://github.com/engagementlabs/safe-code/commit/%s", CommitSHA),
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(versionInfo)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Log the 404 attempt
	log.Printf("404 - Path not found: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	
	// Security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": "Path not found",
		"path": r.URL.Path,
		"available_endpoints": []string{"/health", "/version", "/attestation"},
		"message": "This is a Telegram bot. Available endpoints are listed above.",
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s - %v", r.Method, r.URL.Path, time.Since(start))
	})
}