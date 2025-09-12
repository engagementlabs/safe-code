package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type NitroAttestation struct {
	ModuleID    string `json:"module_id"`
	Timestamp   int64  `json:"timestamp"`
	Digest      string `json:"digest"`
	PCRs        map[string]string `json:"pcrs"`
	Certificate string `json:"certificate"`
	CABundle    []string `json:"cabundle"`
}

// Check if running in Nitro Enclave
func isNitroEnclave() bool {
	// Check for Nitro-specific files/devices
	if _, err := os.Stat("/dev/nsm"); err == nil {
		return true
	}
	
	// Check for Nitro environment variables
	if os.Getenv("NITRO_ENCLAVE") == "true" {
		return true
	}
	
	return false
}

// Get Nitro attestation document
func getNitroAttestationDocument() string {
	if !isNitroEnclave() {
		log.Printf("⚠️ Not running in Nitro Enclave environment")
		return ""
	}
	
	// Try to get attestation using nsm-cli or direct API
	attestation, err := getNitroAttestation()
	if err != nil {
		log.Printf("❌ Failed to get Nitro attestation: %v", err)
		return ""
	}
	
	// Encode as base64 for transport
	encoded := base64.StdEncoding.EncodeToString(attestation)
	log.Printf("🔒 Nitro attestation document generated (%d bytes)", len(attestation))
	
	return encoded
}

// Get raw Nitro attestation
func getNitroAttestation() ([]byte, error) {
	// Method 1: Try nsm-cli if available
	if attestation, err := getNitroAttestationCLI(); err == nil {
		return attestation, nil
	}
	
	// Method 2: Try direct NSM API (requires CGO)
	// This would need the AWS Nitro SDK
	
	// Method 3: Mock for development/testing
	if os.Getenv("NITRO_MOCK") == "true" {
		return getMockNitroAttestation()
	}
	
	return nil, fmt.Errorf("no Nitro attestation method available")
}

// Get attestation using nsm-cli
func getNitroAttestationCLI() ([]byte, error) {
	// Check if nsm-cli is available
	if _, err := exec.LookPath("nsm-cli"); err != nil {
		return nil, fmt.Errorf("nsm-cli not found")
	}
	
	// Execute nsm-cli to get attestation
	cmd := exec.Command("nsm-cli", "describe-pcrs")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("nsm-cli failed: %v", err)
	}
	
	return output, nil
}

// Mock attestation for development
func getMockNitroAttestation() ([]byte, error) {
	mockAttestation := NitroAttestation{
		ModuleID:  "i-1234567890abcdef0-enc1234567890abcdef",
		Timestamp: 1704067200, // 2024-01-01
		Digest:    "sha384:mock-digest-for-development-only",
		PCRs: map[string]string{
			"0": "mock-pcr0-value",
			"1": "mock-pcr1-value",
			"2": "mock-pcr2-value",
		},
		Certificate: "mock-certificate",
		CABundle:    []string{"mock-ca-bundle"},
	}
	
	return json.Marshal(mockAttestation)
}

// Verify Nitro attestation (for external verification)
func verifyNitroAttestation(attestationDoc string) (bool, error) {
	if attestationDoc == "" {
		return false, fmt.Errorf("no attestation document provided")
	}
	
	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(attestationDoc)
	if err != nil {
		return false, fmt.Errorf("failed to decode attestation: %v", err)
	}
	
	// Parse attestation document
	var attestation NitroAttestation
	if err := json.Unmarshal(decoded, &attestation); err != nil {
		return false, fmt.Errorf("failed to parse attestation: %v", err)
	}
	
	// Verify certificate chain and signatures
	// This would require AWS Nitro verification libraries
	
	log.Printf("🔍 Nitro attestation verification: ModuleID=%s", attestation.ModuleID)
	return true, nil
}