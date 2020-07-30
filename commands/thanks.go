package commands

import (
	"errors"
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/caption"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	"strings"
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
	newCaption, err := getThanksCaption(arguments, message, replyToMessage)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, fmt.Sprintf("Error while creating thank caption: %v", err))
	}

	//
	err = dbwrapper.UpdatePostCaptionByUniqueID(fi.Remote.UniqueId, newCaption)
	_ = PreviewCommandHandler{}.Handle("", message, replyToMessage)
	return err

}

func getThanksCaption(arguments string, message, replyToMessage *client.Message) (string, error) {

	// In case there's an error in the thanks part, fall back to
	// just adding the comment, if available
	thanks, err := getThanks(replyToMessage)
	if err != nil {

		if arguments == "" {
			return "", err
		} else {
			thanks = ""
		}

	}

	//
	comment := ""
	if arguments != "" {
		comment = getComment(arguments, message)
	}

	return strings.TrimSpace(fmt.Sprintf("%s\n\n%s", comment, thanks)), err

}

func getComment(arguments string, message *client.Message) string {
	text := message.Content.(*client.MessageText).Text
	msgLengthDifference := len(text.Text) - len(arguments)
	return caption.ToHTMLCaption(text)[msgLengthDifference:]
}

func getThanks(message *client.Message) (string, error) {

	if message.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginChannel {
		return "", errors.New("can't thank channels")
	}

	if message.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginHiddenUser {
		return fmt.Sprintf("[Thanks to %s]", message.ForwardInfo.Origin.(*client.MessageForwardOriginHiddenUser).SenderName), nil
	}

	if message.ForwardInfo.Origin.MessageForwardOriginType() != client.TypeMessageForwardOriginUser {
		return "", errors.New("unsupported forward type")
	}

	fwd := message.ForwardInfo.Origin.(*client.MessageForwardOriginUser)
	user, err := api.GetUserByID(fwd.SenderUserId)
	if err != nil {
		return "", err
	}

	if user.Type.UserTypeType() == client.TypeUserTypeBot {
		return "", errors.New("can't thank bots")
	}

	return fmt.Sprintf("[Thanks to %s]", user.FirstName), nil

}
