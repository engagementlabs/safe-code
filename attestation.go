package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func attestationHandler(w http.ResponseWriter, r *http.Request) {
	// Security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Content-Type", "application/json")
	
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// Generate source hash from GitHub
	sourceHash, err := generateSourceHash(CommitSHA)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to verify source"})
		return
	}
	
	// Generate attestation
	attestation := generateAttestation(CommitSHA, BuildTime, sourceHash)
	
	codeAttestation := map[string]interface{}{
		"commit_sha":   CommitSHA,
		"build_time":   BuildTime,
		"source_hash":  sourceHash,
		"build_hash":   generateBuildHash(),
		"github_url":   fmt.Sprintf("https://github.com/engagementlabs/safe-code/commit/%s", CommitSHA),
		"attestation":  attestation,
	}
	
	// Add Nitro Enclave attestation if available
	if nitroDoc := getNitroAttestationDocument(); nitroDoc != "" {
		codeAttestation["nitro_document"] = nitroDoc
		codeAttestation["nitro_enabled"] = true
	} else {
		codeAttestation["nitro_enabled"] = false
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(codeAttestation)
}

func generateSourceHash(commitSHA string) (string, error) {
	// TODO: Replace with actual repository URL
	url := fmt.Sprintf("https://api.github.com/repos/engagementlabs/safe-code/tarball/%s", commitSHA)
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	hasher := sha256.New()
	_, err = io.Copy(hasher, resp.Body)
	if err != nil {
		return "", err
	}
	
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func generateAttestation(commitSHA, buildTime, sourceHash string) string {
	data := fmt.Sprintf("%s:%s:%s", commitSHA, buildTime, sourceHash)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func generateBuildHash() string {
	data := fmt.Sprintf("%s:%s:%s", Version, CommitSHA, BuildTime)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}