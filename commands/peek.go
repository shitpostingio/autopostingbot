package commands

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/database/database"
	"gitlab.com/shitposting/autoposting-bot/legacy"
	"gitlab.com/shitposting/autoposting-bot/repository"
)

type PeekCommandHandler struct {
}

func (PeekCommandHandler) Handle(arguments string, message *client.Message) error {

	log.Println("PEEK HANDLER")
	nextPost, err := database.GetNextPost(repository.Db)
	if err != nil {
		log.Error(err)
		return err
	}

	mediaType, found := legacy.NewMediaTypeFromOld(nextPost.TypeID)
	if !found {
		err := fmt.Errorf("new media type not found for old typeid %d", nextPost.TypeID)
		log.Error(err)
		return err
	}
	log.Println("Media type found:", mediaType)

	formattedText, err := legacy.NewFormattedTextFromCaption(nextPost.Caption)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Println("Formatted text found:", formattedText.Text)
	_, err = api.SendMedia(mediaType, message.ChatId, nextPost.FileID, formattedText.Text, formattedText.Entities)
	return err

}
