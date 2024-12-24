package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	names = []string{"Andrea", "Nazar"}
	index int
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	webhookURL := os.Getenv("TELEGRAM_WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("TELEGRAM_WEBHOOK_URL environment variable not set")
	}

	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if chatIDStr == "" {
		log.Fatal("TELEGRAM_CHAT_ID environment variable not set")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing TELEGRAM_CHAT_ID into an int: %v\n", err)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to initialize bot api: %v", err)
	}

	// Set up webhook URL
	wh, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = bot.Request(wh); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/webhook", webhookHandler)
	go scheduleDailyMessages(bot, chatID)

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("Error decoding update: %v", err)
		http.Error(w, "Can't process update", http.StatusBadRequest)
		return
	}

	if update.Message == nil {
		return
	}

	fmt.Println("Received message:", update.Message.Text)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func scheduleDailyMessages(bot *tgbotapi.BotAPI, chatID int64) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Calculate duration until 9:00 AM
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), 17, 36, 0, 0, now.Location())
	if now.After(next) {
		next = next.Add(24 * time.Hour)
	}
	time.Sleep(time.Until(next))

	for {
		// Send message
		message := fmt.Sprintf("Time to send a song today, %s! 🙃", names[index])
		index = (index + 1) % len(names)

		msg := tgbotapi.NewMessage(chatID, message)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}

		<-ticker.C
	}
}
