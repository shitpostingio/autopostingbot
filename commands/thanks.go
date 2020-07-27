package commands

import (
	"errors"
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/caption"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type ThanksCommandHandler struct {}

func (ThanksCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

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
	var newCaption string
	if arguments != "" {

		text := message.Content.(*client.MessageText).Text
		msgLengthDifference := len(text.Text) - len(arguments)
		newCaption = caption.ToHTMLCaption(text)
		newCaption = newCaption[msgLengthDifference:]

	}

	//
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
