package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

var (
    Version   = "enclave-v1.0"
    CommitSHA = "nitro-enclave"
    BuildTime = "2025-09-10T13:25:00Z"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
    // CORS headers for browser access
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Content-Type", "text/plain")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    w.Write([]byte("OK from Nitro Enclave"))
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
    // CORS headers for browser access
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Content-Type", "application/json")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    versionInfo := map[string]string{
        "version":    Version,
        "commit_sha": CommitSHA,
        "build_time": BuildTime,
        "environment": "nitro-enclave",
        "github_url": "https://github.com/engagementlabs/safe-code/commit/" + CommitSHA,
    }
    json.NewEncoder(w).Encode(versionInfo)
}

func attestationHandler(w http.ResponseWriter, r *http.Request) {
    // CORS headers for browser access
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Content-Type", "application/json")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    // Mock attestation for enclave (compatible with verify script)
    attestation := map[string]interface{}{
        "commit_sha":   CommitSHA,
        "build_time":   BuildTime,
        "source_hash":  "mock-source-hash-nitro-enclave-" + CommitSHA,
        "build_hash":   "mock-build-hash-nitro-enclave-" + CommitSHA,
        "github_url":   "https://github.com/engagementlabs/safe-code/commit/" + CommitSHA,
        "attestation":  "mock-nitro-attestation-" + CommitSHA,
        "nitro_enabled": true,
        "nitro_document": "mock-nitro-document-base64-encoded",
        "environment":  "nitro-enclave",
        "message": "Running in AWS Nitro Enclave - isolated execution environment",
    }
    
    json.NewEncoder(w).Encode(attestation)
}

func main() {
    log.Println("🔒 Starting Safe-Code in Nitro Enclave...")
    
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/version", versionHandler)
    mux.HandleFunc("/attestation", attestationHandler)
    
    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
    }
    
    log.Println("🚀 Server starting on :8080 (Nitro Enclave)")
    log.Fatal(server.ListenAndServe())
}