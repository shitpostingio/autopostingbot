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
func saveMedia(fileID string, caption string, mediaType MediaType, manager *algo.Manager, userID int) {
	switch mediaType {
	case Image:
		manager.AddImageChannel <- entities.Post{
			Media:   fileID,
			Caption: caption,
			UserID:  uint(userID),
			Categories: []entities.Category{
				entities.Category{Name: "image"},
			},
		}
	case Video:
		manager.AddVideoChannel <- entities.Post{
			Media:   fileID,
			Caption: caption,
			UserID:  uint(userID),
			Categories: []entities.Category{
				entities.Category{Name: "video"},
			},
		}
	}
}
