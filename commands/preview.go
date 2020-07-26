package commands

import (
	"errors"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type PreviewCommandHandler struct {

}

func (PreviewCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

	if replyToMessage == nil {
		return errors.New("reply message nil")
	}

	fi, err := api.GetMediaFileInfo(replyToMessage)
	if err != nil {
		return err
	}

	post, err := dbwrapper.FindPostByUniqueID(fi.Remote.UniqueId)
	if err != nil {
		return err
	}

	entities, _ := api.GetFormattedText(post.Caption)
	_, err = api.SendMedia(post.Media.Type, message.ChatId, message.Id, post.Media.FileID, entities.Text, entities.Entities)
	return err

}
