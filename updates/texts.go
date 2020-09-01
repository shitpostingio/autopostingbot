package updates

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/commands"
	l "gitlab.com/shitposting/autoposting-bot/localization"
	"gitlab.com/shitposting/autoposting-bot/telegram"
	"unicode/utf16"
)

var (
	handlers = map[string]commands.Handler{
		"status":  commands.StatusCommandHandler{},
		"peek":    commands.PeekCommandHandler{},
		"pause":   commands.PauseCommandHandler{},
		"delete":  commands.DeleteCommandHandler{},
		"info":    commands.InfoCommandHandler{},
		"postnow": commands.PostNowCommandHandler{},
		"add":     commands.AddCommandHandler{},
		"caption": commands.CaptionCommandHandler{},
		"thanks":  commands.ThanksCommandHandler{},
		"preview": commands.PreviewCommandHandler{},
		"credit":  commands.CreditCommandHandler{},
	}
)

func handleText(message *client.Message) {

	//
	messageContent := message.Content.(*client.MessageText)
	utf16Text := utf16.Encode([]rune(messageContent.Text.Text))

	//
	command, arguments, isCommand := telegram.GetCommand(utf16Text, messageContent.Text.Entities)
	log.Debugln("Command:", command, " IsCommand", isCommand)
	if !isCommand {
		return
	}

	//
	handler, found := handlers[command]
	if !found {
		log.Error("No handler found for ", command)
		_, _ = api.SendPlainText(message.ChatId, l.GetString(l.UPDATES_TEXTS_COMMAND_UNIMPLEMENTED))
		return
	}

	//
	var err error
	var replyMessage *client.Message
	if message.ReplyToMessageId != 0 {

		replyMessage, err = api.GetMessage(message.ChatId, message.ReplyToMessageId)
		if err != nil {
			log.Error("Unable to get reply to message: ", err)
			_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.UPDATES_TEXTS_UNABLE_TO_GET_REPLY_MESSAGE))
			return
		}

		err = handler.Handle(arguments, message, replyMessage)

	} else {
		err = handler.Handle(arguments, message, nil)
	}

	if err != nil {
		log.Error(err)
	}

}
