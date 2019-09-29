package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// LogRequest posts information about a request to a Discord webhook
func LogRequest(address string) {
	message := map[string]string{"content": fmt.Sprintf("Calendar request from %s", address)}
	json, err := json.Marshal(message)
	if err != nil {
		log.Println("LogRequest:", err)
		return
	}
	log.Println(string(json))
	resp, err := http.Post(os.Getenv("DISCORD_WEBHOOK_URL"), "application/json", strings.NewReader(string(json)))
	if err != nil {
		log.Println("LogRequest:", err)
		return
	}
	defer resp.Body.Close()
}
