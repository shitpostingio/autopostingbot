package commands

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type DeleteCommandHandler struct {
}

func (DeleteCommandHandler) Handle(arguments string, message *client.Message) error {

	fileInfo, err := api.GetMediaFileInfo(message)
	if err != nil {
		return err
	}

	err = dbwrapper.DeletePostByUniqueID(fileInfo.Remote.UniqueId)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, fmt.Sprintf("Unable to delete message: %s", err.Error()))
	} else {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "Message deleted correctly")
	}

	return err

}
