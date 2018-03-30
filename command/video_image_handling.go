package command

import (
	"gitlab.com/shitposting/autoposting-bot/algo"
	"gitlab.com/shitposting/autoposting-bot/database/entities"
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
