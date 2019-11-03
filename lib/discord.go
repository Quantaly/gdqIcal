package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var webhookURL, hasWebhookURL = os.LookupEnv("DISCORD_WEBHOOK_URL")

// LogRequest logs the request to a Discord webhook, if present
func LogRequest(addr, path string) {
	if !hasWebhookURL {
		return
	}

	msg := map[string]interface{}{"content": fmt.Sprintf("Request for %s from %s", path, addr)}
	body, err := json.Marshal(msg)
	if err != nil {
		log.Println("Failed to marshal JSON for webhook:", err)
		return
	}

	http.Post(webhookURL, "application/json", bytes.NewReader(body))
}
