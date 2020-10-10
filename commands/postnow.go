package commands

import (
	"errors"
	"github.com/zelenin/go-tdlib/client"
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/documentstore/dbwrapper"
	l "github.com/shitpostingio/autopostingbot/localization"
	"github.com/shitpostingio/autopostingbot/posting"
)

// PostNowCommandHandler represents the handler of the /postnow command.
type PostNowCommandHandler struct{}

// Handle handles the /postnow command.
// /postnow forces the posting of a media and resets the posting timers.
func (PostNowCommandHandler) Handle(_ string, message, replyToMessage *client.Message) error {

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
	err = posting.RequestPost(&post)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.COMMANDS_POSTNOW_UNSUCCESSFUL))
	} else {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.COMMANDS_POSTNOW_SUCCESSFUL))
	}

	return err

}
