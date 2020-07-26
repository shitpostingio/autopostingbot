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

type CreditCommandHandler struct {

}

func (CreditCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

	if replyToMessage == nil {
		return errors.New("no reply message")
	}

	fi, err := api.GetMediaFileInfo(replyToMessage)
	if err != nil {
		return err
	}

	var newCaption string
	//if arguments != "" {
	//
	//	text := message.Content.(*client.MessageText).Text
	//	msgLengthDifference := len(text.Text) - len(arguments)
	//	fmt.Println("Text: ", text.Text, " len diff: ", msgLengthDifference)
	//
	//	newCaption = caption.ToHTMLCaption(text)
	//	fmt.Println("NewCaption: ", newCaption)
	//	newCaption = newCaption[msgLengthDifference:]
	//	fmt.Println("NC should be: ", newCaption)
	//
	//}

	thanks, err := getCreditCaption(arguments, message, replyToMessage)
	if err == nil {
		newCaption = fmt.Sprintf("%s\n\n%s", newCaption, thanks)
	}

	return dbwrapper.UpdatePostCaptionByUniqueID(fi.Remote.UniqueId, newCaption)

}

//TODO: MIGLIORARE *TANTO*
// Veramente putrido
func getCreditCaption(arguments string, message, replyToMessage *client.Message) (string, error) {

	fmt.Println("Forward type: ", replyToMessage.ForwardInfo.Origin.MessageForwardOriginType())
	urlStart := strings.Index(arguments, "http")

	// We have a URL
	if urlStart != -1 {

		/* WHO TO CREDIT: SUBSTRING UNTIL THE FIRST URL STARTS */
		whoToCredit := arguments[:urlStart]

		/* ISOLATE THE LEFTOVER TEXT */
		leftoverText := arguments[urlStart:]

		/* FIND THE END OF THE URL */
		var url string
		var comment string

		urlEnd := strings.IndexAny(leftoverText, " \t\n")
		if urlEnd != -1 {

			/* URL: SUBSTRING OF leftoverText UNTIL urlEnd */
			url = leftoverText[:urlEnd]

			/* COMMENT: SUBSTRING OF leftoverText FROM urlEnd */
			//comment = strings.TrimSpace(leftoverText[urlEnd:])

			text := message.Content.(*client.MessageText).Text
			msgLengthDifference := len(text.Text) - len(arguments)
			startParsing := msgLengthDifference + urlStart + urlEnd
			comment = caption.ToHTMLCaptionWithCustomStart(message.Content.(*client.MessageText).Text, startParsing)
			fmt.Println("startParsing: ", startParsing, " comment: ", comment)
			comment = strings.TrimSpace(comment)

		} else {
			url = leftoverText
		}

		/* CREATE NEW CAPTION */
		return fmt.Sprintf("%s\n\n[By <a href=\"%s\">%s</a>]", strings.TrimSpace(comment), url, strings.TrimSpace(whoToCredit)), nil

	}

		text := message.Content.(*client.MessageText).Text
		msgLengthDifference := len(text.Text) - len(arguments)
		fmt.Println("Text: ", text.Text, " len diff: ", msgLengthDifference)

		newCaption := ""
		newCaption = caption.ToHTMLCaption(text)
		fmt.Println("NewCaption: ", newCaption)
		newCaption = newCaption[msgLengthDifference:]
		fmt.Println("NC should be: ", newCaption)

		if replyToMessage.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginChannel {
			return newCaption, nil
		}

		if replyToMessage.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginHiddenUser {
			return fmt.Sprintf("%s\n\n[By %s]",newCaption, replyToMessage.ForwardInfo.Origin.(*client.MessageForwardOriginHiddenUser).SenderName), nil
		}

		user, err := api.GetUserByID(replyToMessage.SenderUserId)
		if err != nil {
			return newCaption, err
		}

		if user.Type.UserTypeType() == client.TypeUserTypeBot {
			return newCaption, nil
		}

		return fmt.Sprintf("%s\n\n[By %s]", newCaption, user.FirstName), nil



	//words := strings.Fields(arguments)
	//fmt.Println("len words: ", len(words))
	//if len(words) < 2 {
	//
	//	if replyToMessage.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginChannel {
	//		return "", nil
	//	}
	//
	//	if replyToMessage.ForwardInfo.Origin.MessageForwardOriginType() == client.TypeMessageForwardOriginHiddenUser {
	//		return fmt.Sprintf("[By %s]", replyToMessage.ForwardInfo.Origin.(*client.MessageForwardOriginHiddenUser).SenderName), nil
	//	}
	//
	//	user, err := api.GetUserByID(replyToMessage.SenderUserId)
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	if user.Type.UserTypeType() == client.TypeUserTypeBot {
	//		return "", nil
	//	}
	//
	//	return fmt.Sprintf("[By %s]", user.FirstName), nil
	//
	//}
	//
	//return fmt.Sprintf(`[By <a href="%s">%s</a>]`, words[1], words[0]), nil




}

