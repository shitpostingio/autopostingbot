package utility

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// SendTelegramReply replies with a text to the specified update
func SendTelegramReply(chatID int, messageID int, bot *tgbotapi.BotAPI, text string) {
	msg := tgbotapi.NewMessage(int64(chatID), text)
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = messageID

	bot.Send(msg)
}
