package commands

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type DeleteCommandHandler struct {
}

func (DeleteCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

	fileInfo, err := api.GetMediaFileInfo(replyToMessage)
	if err != nil {
		return err
	}

	err = dbwrapper.DeletePostByUniqueID(fileInfo.Remote.UniqueId)
	if err != nil {
		_, _ = api.SendPlainReplyText(replyToMessage.ChatId, replyToMessage.Id, fmt.Sprintf("Unable to delete message: %s", err.Error()))
	} else {
		_, _ = api.SendPlainReplyText(replyToMessage.ChatId, replyToMessage.Id, "Message deleted correctly")
	}

	return err

}
