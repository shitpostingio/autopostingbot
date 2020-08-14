package commands

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/caption"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	l "gitlab.com/shitposting/autoposting-bot/localization"
)

type CaptionCommandHandler struct {}

func (CaptionCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

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
