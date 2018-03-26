package utility

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// SendTelegramReply replies with a text to the specified update
func SendTelegramReply(update tgbotapi.Update, bot *tgbotapi.BotAPI, text string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}
