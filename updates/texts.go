package updates

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/commands"
	"gitlab.com/shitposting/autoposting-bot/telegram"
	"unicode/utf16"
)

var (
	handlers = map[string]commands.Handler{
		"status": commands.StatusCommandHandler{},
		//"peek":   commands.PeekCommandHandler{},
		"pause": commands.PauseCommandHandler{},
	}
)

func handleText(message *client.Message) {

	log.Println("HANDLETEXT")

	//
	messageContent := message.Content.(*client.MessageText)
	utf16Text := utf16.Encode([]rune(messageContent.Text.Text))

	//
	command, arguments, isCommand := telegram.GetCommand(utf16Text, messageContent.Text.Entities)
	log.Println("Command:", command, " IsCommand", isCommand)
	if !isCommand {
		return
	}

	//
	handler, found := handlers[command]
	if !found {
		log.Error("No handler found for ", command)
		_, _ = api.SendPlainText(message.ChatId, "Unimplemented")
		return
	}

	err := handler.Handle(arguments, message)
	if err != nil {
		log.Error(err)
	}

}
