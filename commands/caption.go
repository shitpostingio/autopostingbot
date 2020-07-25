package commands

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/caption"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type CaptionCommandHandler struct {
}

func (CaptionCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

	fmt.Println("Arguments: ", arguments)

	//
	newCaption := ""

	//
	if arguments != "" {

		text := message.Content.(*client.MessageText).Text
		msgLengthDifference := len(text.Text) - len(arguments)
		fmt.Println("Text: ", text.Text, " len diff: ", msgLengthDifference)

		newCaption = caption.ToHTMLCaption(text)
		fmt.Println("NewCaption: ", newCaption)
		newCaption = newCaption[msgLengthDifference:]
		fmt.Println("NC should be: ", newCaption)

	}

	fi, err := api.GetMediaFileInfo(replyToMessage)
	if err != nil {
		return err
	}

	fmt.Println("FI: ", fi.Remote.UniqueId)
	return dbwrapper.UpdatePostCaptionByUniqueID(fi.Remote.UniqueId, newCaption)

}
