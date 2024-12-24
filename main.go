package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	names     = []string{"Andrea", "Nazar"}
	indexFile = "/data/index.txt"
)

func main() {
	// Hack: workaround for fly.io, since they only allow hourly and daily schedules
	currentHour := time.Now().UTC().Hour()

	// UTC is -1 in winter and -2 in summer (thanks, daylight savings)
	if currentHour != 18 {
		log.Println("Time is not 18. Exiting.")
		os.Exit(0)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
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

	sendMessage(bot, chatID)
}

// Load the last index from a file
func loadIndex() int {
	data, err := os.ReadFile(indexFile)
	if err != nil {
		// Default to 0 if file doesn't exist
		return 0
	}
	index, err := strconv.Atoi(string(data))
	if err != nil {
		return 0
	}
	return index
}

// Save the current index to a file
func saveIndex(index int) {
	err := os.WriteFile(indexFile, []byte(strconv.Itoa(index)), 0644)
	if err != nil {
		log.Fatalf("Failed to save index: %v", err)
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64) {
	// Load the last index from a file or set default
	index := loadIndex()
	message := fmt.Sprintf("Time to send a song today, %s! 🙃", names[index])

	// Prepare and send message
	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	// Update index and save for the next run
	index = (index + 1) % len(names)
	saveIndex(index)

	// Shut down the machine gracefully
	log.Println("Message sent successfully. Exiting.")
	os.Exit(0)
}
