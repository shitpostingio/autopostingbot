package updates

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/telegram"
	"unicode/utf16"
)

func handleText(message *client.Message) {

	//
	messageContent := message.Content.(*client.MessageText)
	utf16Text := utf16.Encode([]rune(messageContent.Text.Text))

	//
	command, isCommand := telegram.GetCommand(utf16Text, messageContent.Text.Entities)
	if !isCommand {
		return
	}

	//
	handler, found := handlers[command]
	if !found {
		log.Error("No handler found for ", command)
		_,_ = api.SendPlainTextMessage(message.ChatId, "Unimplemented")
		return
	}

	err := handler.Handle(message)
	if err != nil {
		log.Error(err)
	}

}

