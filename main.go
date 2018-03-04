package main

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	_, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Panic(err)
	}

}
