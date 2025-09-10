package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func runTelegramBot(token string) {
	log.Println("Starting Telegram bot...")
	offset := 0

	// HTTP client with timeouts
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	for {
		url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=5", token, offset)
		
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		
		if err != nil {
			cancel()
			log.Printf("Error creating request: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		
		resp, err := client.Do(req)
		cancel()
		
		if err != nil {
			log.Printf("Error getting updates: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		
		var result struct {
			Result []Update `json:"result"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()
		
		if err != nil {
			log.Printf("Error decoding response: %v", err)
			continue
		}

		for _, update := range result.Result {
			offset = update.UpdateID + 1
			
			// Input validation
			if update.Message.Chat.ID != 0 && update.Message.Text != "" {
				// Rate limiting - simple check
				if len(result.Result) > 100 {
					log.Printf("Rate limit: too many updates, skipping")
					continue
				}
				
				switch update.Message.Text {
				case "/ping":
					log.Printf("✅ Received /ping from chat %d", update.Message.Chat.ID)
					sendMessage(client, token, update.Message.Chat.ID, "pong")
				case "/version":
					log.Printf("📋 Received /version from chat %d", update.Message.Chat.ID)
					versionMsg := fmt.Sprintf("🔍 Version Info:\n• Version: %s\n• Commit: %s\n• Build: %s\n• GitHub: https://github.com/engagementlabs/safe-code/commit/%s", Version, CommitSHA, BuildTime, CommitSHA)
					sendMessage(client, token, update.Message.Chat.ID, versionMsg)
				case "/attestation":
					log.Printf("🔐 Received /attestation from chat %d", update.Message.Chat.ID)
					attestationMsg := "🔐 Code Attestation:\nFor cryptographic proof, visit:\nhttps://safe-code-p.engagelabs.org/attestation"
					sendMessage(client, token, update.Message.Chat.ID, attestationMsg)
				default:
					// Only send help for commands that start with /
					if strings.HasPrefix(update.Message.Text, "/") {
						log.Printf("❓ Unknown command '%s' from chat %d", update.Message.Text, update.Message.Chat.ID)
						helpMsg := "🤖 Available commands:\n/ping - I'll respond with 'pong'\n/version - Show build information\n/attestation - Code verification info\n\n🔍 Web endpoints:\n• /health - Health check\n• /version - Build info (JSON)\n• /attestation - Code verification (JSON)"
						sendMessage(client, token, update.Message.Chat.ID, helpMsg)
					} else {
						// Just log regular messages, don't respond
						log.Printf("💬 Regular message from chat %d: %s", update.Message.Chat.ID, update.Message.Text)
					}
				}
			}
		}
	}
}

func sendMessage(client *http.Client, token string, chatID int64, text string) {
	// Input validation
	if chatID == 0 || len(text) == 0 || len(text) > 4096 {
		log.Printf("Invalid message parameters")
		return
	}
	
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
		"parse_mode": "HTML", // Prevent injection
	}
	
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error sending message: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		log.Printf("✅ Message sent successfully to chat %d", chatID)
	} else {
		log.Printf("⚠️ Message send failed with status %d for chat %d", resp.StatusCode, chatID)
	}
}