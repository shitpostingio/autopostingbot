package main

import (
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	_, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(ReadConfigFile("./config_example.toml"))

}
