package command

import (
	"errors"
	"fmt"

	"gitlab.com/shitposting/autoposting-bot/utility"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/shitposting/autoposting-bot/algo"
)

// Handle e` il punto di entrata per il parsing e l'organizzazione dell'azione del bot
// su un messaggio entrante.
func Handle(update tgbotapi.Update, api *tgbotapi.BotAPI, manager *algo.Manager) error {
	if update.Message == nil && update.EditedMessage == nil {
		return errors.New("update Message or EditedMessage body was nil, most likely an error on Telegram side")
	}

	msg := update.Message
	editedMsg := update.EditedMessage

	if editedMsg != nil {
		switch {
		case editedMsg.Video != nil:
			modifyMedia(editedMsg.Video.FileID, editedMsg.Caption, manager, editedMsg.From.ID, editedMsg.MessageID, int(editedMsg.Chat.ID))
		case editedMsg.Photo != nil:
			photos := *editedMsg.Photo
			modifyMedia(photos[len(photos)-1].FileID, editedMsg.Caption, manager, editedMsg.From.ID, editedMsg.MessageID, int(editedMsg.Chat.ID))
		case editedMsg.Text != "":
			switch editedMsg.Command() {
			case "caption":
				editCaption(editedMsg, api, manager, false)
			case "credit":
				editCaption(editedMsg, api, manager, true)
			}
		}

		return nil
	}

	switch {
	case msg.Video != nil:
		saveMedia(msg.Video.FileID, msg.Caption, Video, manager, msg.From.ID, msg.MessageID, int(msg.Chat.ID))
	case msg.Photo != nil:
		photos := *msg.Photo
		saveMedia(photos[len(photos)-1].FileID, msg.Caption, Image, manager, msg.From.ID, msg.MessageID, int(msg.Chat.ID))
	case msg.Text != "":
		switch msg.Command() {
		case "status":
			manager.SendStatusInfo(msg.MessageID, int(msg.Chat.ID))
		case "delete":
			manager.DeleteMedia(msg)
		case "caption":
			editCaption(msg, api, manager, false)
		case "credit":
			editCaption(msg, api, manager, true)
		}

	}

	return nil
}

// editCaption allows the user to edit the caption of a forwarded message or give the credit to the user.
// It is used both by caption and credit command in the bot.
func editCaption(msg *tgbotapi.Message, api *tgbotapi.BotAPI, manager *algo.Manager, isCredit bool) {

	var fileID, newcaption string

	if msg.ReplyToMessage == nil {
		utility.SendTelegramReply(int(msg.Chat.ID), msg.MessageID, api, "Not a reply!")
		return
	}

	if msg.ReplyToMessage.ForwardFrom != nil && isCredit == true {
		newcaption = fmt.Sprintf("%s\n\n[By %s]", msg.CommandArguments(), msg.ReplyToMessage.ForwardFrom.FirstName)
	} else {
		newcaption = msg.CommandArguments()
	}

	switch {
	case msg.ReplyToMessage.Photo != nil:
		photosID := *msg.ReplyToMessage.Photo
		fileID = photosID[len(photosID)-1].FileID
	case msg.ReplyToMessage.Video != nil:
		fileID = msg.ReplyToMessage.Video.FileID
	default:
		utility.SendTelegramReply(int(msg.Chat.ID), msg.MessageID, api, "Not a media!")
		return
	}

	modifyMedia(fileID, newcaption, manager, msg.From.ID, msg.MessageID, int(msg.Chat.ID))
}
