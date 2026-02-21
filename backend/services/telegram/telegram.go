package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jwnpoh/njcreaderapp/backend/cmd/config"
)

type TelegramService struct {
	botToken string
	chatID   string
}

type TelegramPayload struct {
	Title  string   `json:"title"`
	URL    string   `json:"url"`
	Topics []string `json:"topics"`
}

func NewTelegramService(cfg config.TelegramConfig) *TelegramService {
	return &TelegramService{
		botToken: cfg.BotToken,
		chatID:   cfg.ChatID,
	}
}

func (t *TelegramService) SendDigest(articles []TelegramPayload) error {
	if len(articles) == 0 {
		return fmt.Errorf("telegram: no articles to send")
	}

	message := t.formatMessage(articles)

	return t.sendMessage(message)
}

func (t *TelegramService) formatMessage(articles []TelegramPayload) string {
	var sb strings.Builder

	sb.WriteString("<b>⭐ Today's Must Read(s) ⭐</b>\n\n")

	for i, article := range articles {
		// Title as hyperlink
		sb.WriteString(fmt.Sprintf("<b>%d. <a href=\"%s\">%s</a></b>\n", i+1, article.URL, article.Title))

		// Topics as clickable links to search
		if len(article.Topics) > 0 {
			sb.WriteString("📑 Topics: ")
			topicLinks := make([]string, len(article.Topics))
			for j, topic := range article.Topics {
				searchURL := fmt.Sprintf("https://the-njc-reader.vercel.app/search?term=%s", topic)
				topicLinks[j] = fmt.Sprintf("<a href=\"%s\">%s</a>", searchURL, topic)
			}
			sb.WriteString(strings.Join(topicLinks, ", "))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (t *TelegramService) sendMessage(text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	payload := map[string]interface{}{
		"chat_id":    t.chatID,
		"text":       text,
		"parse_mode": "HTML",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram service: failed to marshal payload - %w", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("telegram service: request failed - %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("telegram service: failed to decode response - %w", err)
	}

	if ok, _ := result["ok"].(bool); !ok {
		return fmt.Errorf("telegram service: sendMessage failed - %v", result)
	}

	return nil
}
