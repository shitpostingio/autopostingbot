package command

import (
	"errors"

	"gitlab.com/shitposting/autoposting-bot/algo"
	"gitlab.com/shitposting/autoposting-bot/database/entities"
	"gitlab.com/shitposting/autoposting-bot/utility"
	"gitlab.com/shitposting/telegram-bot-api"
)

// MediaType is the type of media we're dealing with
type MediaType int

//go:generate stringer -type=MediaType
const (
	Image MediaType = iota
	Video
)

// saveMedia sends the media identified by the fileID to the Manager
func saveMedia(fileID string, caption string, mediaType MediaType, manager *algo.Manager, userID int, messageID int, chatID int) {
	switch mediaType {
	case Image:
		e := entities.Post{
			Media:   fileID,
			Caption: caption,
			UserID:  uint(userID),
			Categories: []entities.Category{
				entities.Category{Name: "image"},
			},
		}

		manager.AddImageChannel <- algo.MediaPayload{
			ChatID:    chatID,
			MessageID: messageID,
			Entity:    e,
		}
	case Video:
		e := entities.Post{
			Media:   fileID,
			Caption: caption,
			UserID:  uint(userID),
			Categories: []entities.Category{
				entities.Category{Name: "video"},
			},
		}

		manager.AddVideoChannel <- algo.MediaPayload{
			ChatID:    chatID,
			MessageID: messageID,
			Entity:    e,
		}
	}
}

// modifyMedia sends the new entity identified by its fileID to the manager, to be
// modified in the database structure
func modifyMedia(fileID string, caption string, manager *algo.Manager, userID int, messageID int, chatID int) {
	e := entities.Post{Media: fileID, Caption: caption, UserID: uint(userID)}

	manager.ModifyMediaChannel <- algo.MediaPayload{
		ChatID:    chatID,
		MessageID: messageID,
		Entity:    e,
	}
}

func checkReplyAndMedia(msg *tgbotapi.Message) (string, error) {

	if msg.ReplyToMessage == nil {
		err := errors.New("not a reply")
		return "", err
	}

	switch {
	case msg.ReplyToMessage.Photo != nil:
		photosID := *msg.ReplyToMessage.Photo
		fileID := photosID[len(photosID)-1].FileID
		return fileID, nil
	case msg.ReplyToMessage.Video != nil:
		fileID := msg.ReplyToMessage.Video.FileID
		return fileID, nil
	default:
		err := errors.New("not a media")
		return "", err
	}
}

func deleteMedia(msg *tgbotapi.Message, api *tgbotapi.BotAPI, manager *algo.Manager) {

	fileID, err := checkReplyAndMedia(msg)

	if err != nil {
		utility.SendTelegramReply(int(msg.Chat.ID), msg.MessageID, api, err.Error())
		return
	}

	e := entities.Post{Media: fileID, UserID: uint(msg.From.ID)}

	manager.DeleteMediaChannel <- algo.MediaPayload{
		ChatID:    int(msg.Chat.ID),
		MessageID: msg.MessageID,
		Entity:    e,
	}

}

func statusSignal(msg *tgbotapi.Message, manager *algo.Manager) {
	e := entities.Post{}

	manager.StatusChannel <- algo.MediaPayload{
		ChatID:    int(msg.Chat.ID),
		MessageID: msg.MessageID,
		Entity:    e,
	}
}
