package commands

import (
	"errors"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type PreviewCommandHandler struct {}

func (PreviewCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

	//
	if replyToMessage == nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "This command needs to be used in reply to a media file")
		return errors.New("reply to message nil")
	}

	//
	fi, err := api.GetMediaFileInfo(replyToMessage)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "This command needs to be used in reply to a media file")
		return err
	}

	//
	post, err := dbwrapper.FindPostByUniqueID(fi.Remote.UniqueId)
	if err != nil {
		_, _ = api.SendPlainReplyText(replyToMessage.ChatId, replyToMessage.Id, "Unable to find the post")
		return err
	}

	//
	ft, err := api.GetFormattedText(post.Caption)
	if err != nil {
		ft = &client.FormattedText{
			Text:     post.Caption,
			Entities: nil,
		}
	}

	//
	_, err = api.SendMedia(post.Media.Type, message.ChatId, message.Id, post.Media.FileID, ft.Text, ft.Entities)
	return err

}
