package commands

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type PeekCommandHandler struct {
}

//func (PeekCommandHandler) Handle(arguments string, message *client.Message) error {
//
//	log.Println("PEEK HANDLER")
//	nextPost, err := database.GetNextPost(repository.Db)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//
//	mediaType, found := legacy.NewMediaTypeFromOld(nextPost.TypeID)
//	if !found {
//		err := fmt.Errorf("new media type not found for old typeid %d", nextPost.TypeID)
//		log.Error(err)
//		return err
//	}
//	log.Println("Media type found:", mediaType)
//
//	formattedText, err := legacy.NewFormattedTextFromCaption(nextPost.Caption)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//
//	log.Println("Formatted text found:", formattedText.Text)
//	_, err = api.SendMedia(mediaType, message.ChatId, message.Id, nextPost.FileID, formattedText.Text, formattedText.Entities)
//	return err
//
//}

func (PeekCommandHandler) Handle(arguments string, message *client.Message) error {

	log.Println("PEEK HANDLER")
	nextPost, err := dbwrapper.GetNextPost()
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = api.SendMedia(nextPost.Media.Type, message.ChatId, message.Id, nextPost.Media.FileID, nextPost.Caption.Text, nextPost.Caption.Entities)
	return err

}
