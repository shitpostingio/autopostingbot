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

// CreditCommandHandler represents the handler of the /credit command.
type CreditCommandHandler struct{}

// Handle handles the /credit command.
// /credit allows to credit the original creator of a post. It supports hyperlinking of the source.
// By using /credit without any additional argument on a forwarded, the original sender of the message will be credited.
// By using /credit with a name and a link, the source will be credited and put as a hyperlink.
func (CreditCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

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
	credit, err := getCreditCaption(arguments, message, replyToMessage)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.COMMANDS_CREDIT_UNABLE_TO_CREDIT))
		credit = ""
	}

	//
	err = dbwrapper.UpdatePostCaptionByUniqueID(fi.Remote.UniqueId, credit)

	//
	_ = PreviewCommandHandler{}.Handle("", message, replyToMessage)
	return err

}

// getCreditCaption returns the correct credit caption, checking whether a URL is present or not.
func getCreditCaption(arguments string, message, replyToMessage *client.Message) (string, error) {

	// Different behavior whether we have a URL or not
	urlStart := strings.Index(arguments, "http")

	// We have a URL
	if urlStart != -1 {
		return creditCaptionForURL(urlStart, arguments, message)
	}

	return creditCaptionWithoutURL(arguments, message, replyToMessage)

}

// creditCaptionForURL returns a credit caption with a hyperlinked source.
func creditCaptionForURL(urlStart int, arguments string, message *client.Message) (string, error) {

	/*
	 *	STRUCTURE:
	 *	/credit whoToCredit http://some-url additional comments
	 */

	//
	whoToCredit := strings.TrimSpace(arguments[:urlStart])
	leftoverText := arguments[urlStart:]

	//
	urlEnd := strings.IndexAny(leftoverText, " \t\n")
	if urlEnd == -1 {
		return fmt.Sprintf(l.GetString(l.COMMANDS_CREDIT_CAPTION_WITH_URL), leftoverText, whoToCredit), nil
	}

	// Isolate the URL
	url := leftoverText[:urlEnd]

	// Parse entities and isolate comment
	text := message.Content.(*client.MessageText).Text
	msgLengthDifference := len(text.Text) - len(arguments)
	commentStart := msgLengthDifference + urlStart + urlEnd
	comment := caption.ToHTMLCaptionWithCustomStart(message.Content.(*client.MessageText).Text, commentStart)
	comment = strings.TrimSpace(comment[commentStart:])

	//
	captionEnd := fmt.Sprintf(l.GetString(l.COMMANDS_CREDIT_CAPTION_WITH_URL), url, whoToCredit)
	return fmt.Sprintf("%s\n\n%s", comment, captionEnd), nil

}

// creditCaptionWithoutURL returns the credit caption thanking the user who originally sent the message.
// It will not thank bots or channels.
func creditCaptionWithoutURL(arguments string, message, replyToMessage *client.Message) (string, error) {

	//
	text := message.Content.(*client.MessageText).Text
	msgLengthDifference := len(text.Text) - len(arguments)
	newCaption := caption.ToHTMLCaption(text)
	newCaption = newCaption[msgLengthDifference:]

	// Channel forwards shouldn't be credited
	if replyToMessage.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginChannel {
		return newCaption, nil
	}

	// Hidden users can only be users
	if replyToMessage.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginHiddenUser {
		captionEnd := fmt.Sprintf(l.GetString(l.COMMANDS_CREDIT_CAPTION_WITHOUT_URL), replyToMessage.ForwardInfo.Origin.(*client.MessageForwardOriginHiddenUser).SenderName)
		return fmt.Sprintf("%s\n\n%s", newCaption, captionEnd), nil
	}

	if replyToMessage.ForwardInfo.Origin.MessageForwardOriginType() != client.TypeMessageForwardOriginUser {
		return "", errors.New("unsupported forward type")
	}

	fwd := replyToMessage.ForwardInfo.Origin.(*client.MessageForwardOriginUser)
	user, err := api.GetUserByID(fwd.SenderUserId)
	if err != nil {
		return newCaption, err
	}

	//
	if user.Type.UserTypeType() == client.TypeUserTypeBot {
		return newCaption, nil
	}

	// Use only the first name for normal users
	captionEnd := fmt.Sprintf(l.GetString(l.COMMANDS_CREDIT_CAPTION_WITHOUT_URL), user.FirstName)
	return fmt.Sprintf("%s\n\n%s", newCaption, captionEnd), nil

}
