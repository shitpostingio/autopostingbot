package commands

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/caption"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type CaptionCommandHandler struct {
}

func (CaptionCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

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
	newCaption := ""
	if arguments != "" && message.Content.MessageContentType() == client.TypeMessageText {

		//
		text := message.Content.(*client.MessageText).Text
		msgLengthDifference := len(text.Text) - len(arguments)
		log.Debugln("Text: ", text.Text, " len diff: ", msgLengthDifference)

		//
		newCaption = caption.ToHTMLCaptionWithCustomStart(text, msgLengthDifference)
		newCaption = newCaption[msgLengthDifference:]
		log.Debugln("new caption is: ", newCaption)

	}

	// Save new caption to database
	err = dbwrapper.UpdatePostCaptionByUniqueID(fi.Remote.UniqueId, newCaption)

	// Send how the new post looks like
	_ = PreviewCommandHandler{}.Handle("", message, replyToMessage)
	return err

}
