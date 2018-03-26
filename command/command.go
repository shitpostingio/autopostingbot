package command

import (
	"errors"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/shitposting/autoposting-bot/algo"
	"gitlab.com/shitposting/autoposting-bot/utility"
)

// Handle e` il punto di entrata per il parsing e l'organizzazione dell'azione del bot
// su un messaggio entrante.
func Handle(update tgbotapi.Update, api *tgbotapi.BotAPI, manager *algo.Manager) error {
	if update.Message == nil {
		return errors.New("update Message body was nil, most likely an error on Telegram side")
	}

	msg := update.Message

	switch {
	case msg.Video != nil:
		saveMedia(msg.Video.FileID, msg.Caption, Video, manager, msg.From.ID)
		utility.SendTelegramReply(update, api, "Added!")
	case msg.Photo != nil:
		photos := *msg.Photo
		saveMedia(photos[len(photos)-1].FileID, msg.Caption, Photo, manager, msg.From.ID)
		utility.SendTelegramReply(update, api, "Added!")
	}

	return nil
}
