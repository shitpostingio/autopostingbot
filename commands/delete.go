package commands

import (
	"errors"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type DeleteCommandHandler struct {}

func (DeleteCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

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
	err = dbwrapper.DeletePostByUniqueID(fi.Remote.UniqueId)
	if err != nil {
		_, _ = api.SendPlainReplyText(replyToMessage.ChatId, replyToMessage.Id, "Unable to delete the post")
	} else {
		_, _ = api.SendPlainReplyText(replyToMessage.ChatId, replyToMessage.Id, "Post deleted correctly")
	}

	return err

}
