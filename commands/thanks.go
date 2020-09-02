package commands

import (
	"errors"
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/caption"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	l "gitlab.com/shitposting/autoposting-bot/localization"
	"strings"
)

// ThanksCommandHandler represents the handler of the /thanks command.
type ThanksCommandHandler struct{}

// Handle handles the /thanks command.
// /thanks thanks the user that originally sent the media.
func (ThanksCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

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
	newCaption, err := getThanksCaption(arguments, message, replyToMessage)
	if err != nil {
		thankError := fmt.Sprintf(l.GetString(l.COMMANDS_THANK_UNABLE_TO_THANK), err)
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, thankError)
	}

	//
	err = dbwrapper.UpdatePostCaptionByUniqueID(fi.Remote.UniqueId, newCaption)
	_ = PreviewCommandHandler{}.Handle("", message, replyToMessage)
	return err

}

// getThanksCaption returns the thank user caption.
func getThanksCaption(arguments string, message, replyToMessage *client.Message) (string, error) {

	// In case there's an error in the thanks part, fall back to
	// just adding the comment, if available
	thanks, err := getThanks(replyToMessage)
	if err != nil {

		if arguments == "" {
			return "", err
		}

		thanks = ""

	}

	//
	comment := ""
	if arguments != "" {
		comment = getComment(arguments, message)
	}

	return strings.TrimSpace(fmt.Sprintf("%s\n\n%s", comment, thanks)), err

}

// getComment returns any additional comment that the user wants to add to the thank.
func getComment(arguments string, message *client.Message) string {
	text := message.Content.(*client.MessageText).Text
	msgLengthDifference := len(text.Text) - len(arguments)
	return caption.ToHTMLCaption(text)[msgLengthDifference:]
}

// getThanks thanks the user that originally sent the media.
// It will not thank bots or channels.
func getThanks(message *client.Message) (string, error) {

	if message.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginChannel {
		return "", errors.New(l.GetString(l.COMMANDS_THANK_CANT_THANK_CHANNELS))
	}

	if message.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginHiddenUser {
		return fmt.Sprintf(l.GetString(l.COMMANDS_THANK_THANK_CAPTION), message.ForwardInfo.Origin.(*client.MessageForwardOriginHiddenUser).SenderName), nil
	}

	if message.ForwardInfo.Origin.MessageForwardOriginType() != client.TypeMessageForwardOriginUser {
		return "", errors.New(l.GetString(l.COMMANDS_THANK_UNSUPPORTED_FORWARD_TYPE))
	}

	fwd := message.ForwardInfo.Origin.(*client.MessageForwardOriginUser)
	user, err := api.GetUserByID(fwd.SenderUserId)
	if err != nil {
		return "", err
	}

	if user.Type.UserTypeType() == client.TypeUserTypeBot {
		return "", errors.New(l.GetString(l.COMMANDS_THANK_CANT_THANK_BOTS))
	}

	return fmt.Sprintf(l.GetString(l.COMMANDS_THANK_THANK_CAPTION), user.FirstName), nil

}
