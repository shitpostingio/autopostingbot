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

type CreditCommandHandler struct {}

func (CreditCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

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
	credit, err := getCreditCaption(arguments, message, replyToMessage)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "Unable to credit correctly")
		credit = ""
	}

	//
	err = dbwrapper.UpdatePostCaptionByUniqueID(fi.Remote.UniqueId, credit)

	//
	_ = PreviewCommandHandler{}.Handle("", message, replyToMessage)
	return err

}

func getCreditCaption(arguments string, message, replyToMessage *client.Message) (string, error) {

	// Different behavior whether we have a URL or not
	urlStart := strings.Index(arguments, "http")

	// We have a URL
	if urlStart != -1 {
		return creditCaptionForURL(urlStart, arguments, message, replyToMessage)
	}

	return creditCaptionWithoutURL(arguments, message, replyToMessage)

}

func creditCaptionForURL(urlStart int, arguments string, message, replyToMessage *client.Message) (string, error) {

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
		return fmt.Sprintf("[By <a href=\"%s\">%s</a>]", leftoverText, whoToCredit), nil
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
	return fmt.Sprintf("%s\n\n[By <a href=\"%s\">%s</a>]", comment, url, whoToCredit), nil

}

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
		return fmt.Sprintf("%s\n\n[By %s]", newCaption, replyToMessage.ForwardInfo.Origin.(*client.MessageForwardOriginHiddenUser).SenderName), nil
	}

	// Get the user and check it isn't a bot
	user, err := api.GetUserByID(replyToMessage.SenderUserId)
	if err != nil {
		return newCaption, err
	}

	//
	if user.Type.UserTypeType() == client.TypeUserTypeBot {
		return newCaption, nil
	}

	// Use only the first name for normal users
	return fmt.Sprintf("%s\n\n[By %s]", newCaption, user.FirstName), nil

}

