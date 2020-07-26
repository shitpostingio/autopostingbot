package commands

import (
	"errors"
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/caption"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type ThanksCommandHandler struct {
}

func (ThanksCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

	if replyToMessage == nil {
		return errors.New("no reply message")
	}

	fi, err := api.GetMediaFileInfo(replyToMessage)
	if err != nil {
		return err
	}

	var newCaption string
	if arguments != "" {

		text := message.Content.(*client.MessageText).Text
		msgLengthDifference := len(text.Text) - len(arguments)
		fmt.Println("Text: ", text.Text, " len diff: ", msgLengthDifference)

		newCaption = caption.ToHTMLCaption(text)
		fmt.Println("NewCaption: ", newCaption)
		newCaption = newCaption[msgLengthDifference:]
		fmt.Println("NC should be: ", newCaption)

	}

	thanks, err := getThanksCaption(replyToMessage)
	if err == nil {
		newCaption = fmt.Sprintf("%s\n\n%s", newCaption, thanks)
	}

	return dbwrapper.UpdatePostCaptionByUniqueID(fi.Remote.UniqueId, newCaption)

}

//TODO: MIGLIORARE *TANTO*
func getThanksCaption(message *client.Message) (string, error) {

	fmt.Println("Forward type: ", message.ForwardInfo.Origin.MessageForwardOriginType())

	if message.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginChannel {
		return "", nil
	}

	if message.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginHiddenUser {
		return fmt.Sprintf("[Thanks to %s]", message.ForwardInfo.Origin.(*client.MessageForwardOriginHiddenUser).SenderName), nil
	}

	user, err := api.GetUserByID(message.SenderUserId)
	if err != nil {
		return "", err
	}

	if user.Type.UserTypeType() == client.TypeUserTypeBot {
		return "", nil
	}

	return fmt.Sprintf("[Thanks to %s]", user.FirstName), nil

}
