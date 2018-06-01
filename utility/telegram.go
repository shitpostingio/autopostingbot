package utility

import (
	"github.com/empetrone/telegram-bot-api"
)

// SendTelegramReply replies with a text to the specified update
func SendTelegramReply(chatID int, messageID int, bot *tgbotapi.BotAPI, text string) {
	msg := tgbotapi.NewMessage(int64(chatID), text)
	msg.ReplyToMessageID = messageID

	bot.Send(msg)
}
