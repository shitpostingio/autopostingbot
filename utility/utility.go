package utility

import (
	"io"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// CloseSafely closes an entity and logs in case of errors
func CloseSafely(toClose io.Closer) {
	err := toClose.Close()
	if err != nil {
		log.Println(err)
	}
}

// FormatDate formats a date
func FormatDate(date time.Time) string {
	return date.Format("Mon _2 Jan 2006 15:04:05")
}

// GetMessageEntities returns the correct message entities
func GetMessageEntities(message *tgbotapi.Message) []tgbotapi.MessageEntity {
	if message.Entities != nil {
		return message.Entities
	}

	return message.CaptionEntities
}

// GetHandleOrName returns the handle or the name of a user
func GetHandleOrName(user *tgbotapi.User) string {

	if user.UserName != "" {
		return "@" + user.UserName
	}

	if user.LastName != "" {
		return user.FirstName + " " + user.LastName
	}

	return user.FirstName
}
