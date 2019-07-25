package utility

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// GetFileIDFromMessage returns fileid of the message if the media type is supported
func GetFileIDFromMessage(msg *tgbotapi.Message) (fileID string, err error) {

	if msg == nil {
		err = errors.New("message is nil")
		return
	}

	switch {
	case msg.Photo != nil:
		fileID = msg.Photo[len(msg.Photo)-1].FileID
	case msg.Animation != nil:
		fileID = msg.Animation.FileID
	case msg.Video != nil:
		fileID = msg.Video.FileID
	}

	if fileID == "" {
		err = errors.New("not a supported media message")
	}

	return
}
