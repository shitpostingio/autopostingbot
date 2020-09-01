package commands

import (
	"errors"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	l "gitlab.com/shitposting/autoposting-bot/localization"
)

type PreviewCommandHandler struct {}

func (PreviewCommandHandler) Handle(_ string, message, replyToMessage *client.Message) error {

	//
	if replyToMessage == nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.COMMANDS_REPLY_TO_MEDIA_FILE))
		return errors.New("reply to message nil")
	}

	//
	fi, err := api.GetMediaFileInfo(replyToMessage)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.COMMANDS_REPLY_TO_MEDIA_FILE))
		return err
	}

	//
	post, err := dbwrapper.FindPostByUniqueID(fi.Remote.UniqueId)
	if err != nil {
		_, _ = api.SendPlainReplyText(replyToMessage.ChatId, replyToMessage.Id, l.GetString(l.DATABASE_UNABLE_TO_FIND_POST))
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
