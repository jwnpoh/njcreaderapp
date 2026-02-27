package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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

	titles := []string{
		"📚 What you need to know today",
		"🎯 Don't miss these stories",
		"⚡ Today's essential reads",
		"🔥 Hot takes for smart students",
		"💡 Insights you can't skip",
		"🌟 Your daily knowledge boost",
		"🎓 Sharp minds read these",
		"⭐ Stories shaping the conversation",
		"🚀 Fuel your brain today",
		"📰 The smart student's briefing",
		"🧠 Think deeper with these",
		"🎯 Required reading for today",
		"💎 Premium picks for you",
		"🔍 What everyone's talking about",
		"📊 Data-driven stories today",
		"🌍 Global issues, local impact",
		"⚡ Quick reads, big ideas",
		"🎨 Perspectives that matter",
		"🔥 Trending in current affairs",
		"💭 Food for thought today",
		"📚 Expand your worldview",
		"🎯 Sharp takes for sharp minds",
		"⭐ Today's conversation starters",
		"🧭 Navigate today's headlines",
		"💡 Illuminate your understanding",
		"🎓 Level up your awareness",
		"🔔 Stories you actually need",
		"🌟 Today's standout journalism",
		"⚡ Brief but brilliant reads",
		"🎯 Your competitive edge today",
		"🔥 What's worth your time",
		"📈 Ideas on the rise",
		"💡 Thought-provoking reads",
		"🎯 The informed student's pick",
		"🌟 Quality over quantity today",
		"🧠 Challenge your thinking",
		"📚 Beyond the classroom",
		"⚡ Stay ahead of the curve",
		"🎓 Smart reads for today",
		"💎 Curated for curious minds",
		"🔍 Deep dives worth taking",
		"🌍 Stories that connect",
		"📊 Numbers tell stories too",
		"🎯 Signal in the noise",
		"💡 Fresh perspectives today",
		"🔥 Can't-miss coverage",
		"⭐ Your strategic advantage",
		"🧭 Navigate complexity here",
		"📰 Context you're missing",
		"🚀 Accelerate your learning",
	}

	// Rotate based on day of year for consistency within a day
	dayOfYear := time.Now().YearDay()
	selectedTitle := titles[dayOfYear%len(titles)]

	sb.WriteString(fmt.Sprintf("<b>%s</b>\n\n", selectedTitle))

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
