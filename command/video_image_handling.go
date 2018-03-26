package command

import (
	"gitlab.com/shitposting/autoposting-bot/algo"
	"gitlab.com/shitposting/autoposting-bot/database/entities"
)

// MediaType is the type of media we're dealing with
type MediaType int

//go:generate stringer -type=MediaType
const (
	Photo MediaType = iota
	Video
)

// saveMedia sends the media identified by the fileID to the Manager
func saveMedia(fileID string, caption string, mediaType MediaType, manager *algo.Manager, userID int) {
	manager.AddChannel <- entities.Post{Media: fileID, Caption: caption, UserID: uint(userID)}
}
