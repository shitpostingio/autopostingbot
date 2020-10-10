package commands

import (
	"errors"
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/documentstore/dbwrapper"
	l "github.com/shitpostingio/autopostingbot/localization"
	"github.com/zelenin/go-tdlib/client"
)

// DeleteCommandHandler represents the handler of the /delete command.
type DeleteCommandHandler struct{}

// Handle handles the /delete command.
// /delete deletes a post (that has yet to be posted) from the database.
func (DeleteCommandHandler) Handle(_ string, message, replyToMessage *client.Message) error {

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
	err = dbwrapper.DeletePostByUniqueID(fi.Remote.UniqueId)
	if err != nil {
		_, _ = api.SendPlainReplyText(replyToMessage.ChatId, replyToMessage.Id, l.GetString(l.COMMANDS_DELETE_FAILURE))
	} else {
		_, _ = api.SendPlainReplyText(replyToMessage.ChatId, replyToMessage.Id, l.GetString(l.COMMANDS_DELETE_SUCCESS))
	}

	return err

}
