package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN environment variable is not set")
	}

	chatID := os.Getenv("CHAT_ID")
	if chatID == "" {
		log.Fatal("CHAT_ID environment variable is not set")
	}

	// Test token with getMe
	fmt.Println("Testing bot token...")
	getMe(token)

	// Send a message to the channel
	fmt.Println("Sending message to channel...")
	sendMessage(token, chatID, "Hello from testbot!")
}

func getMe(token string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", token)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	result := parseResponse(resp)

	pretty, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(pretty))

	if ok, _ := result["ok"].(bool); !ok {
		log.Fatal("authentication failed: token is invalid")
	}

	fmt.Println("\nBot token is valid!\n")
}

func sendMessage(token, chatID, text string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id": {chatID},
		"text":    {text},
	})
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	result := parseResponse(resp)

	pretty, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(pretty))

	if ok, _ := result["ok"].(bool); !ok {
		log.Fatalf("sendMessage failed: %s", pretty)
	}

	fmt.Println("\nMessage sent successfully!")
}

func parseResponse(resp *http.Response) map[string]interface{} {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("failed to parse response: %v", err)
	}

	return result
}
