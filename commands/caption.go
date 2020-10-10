package commands

import (
	"errors"
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/caption"
	"github.com/shitpostingio/autopostingbot/documentstore/dbwrapper"
	l "github.com/shitpostingio/autopostingbot/localization"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
)

// CaptionCommandHandler represents the handler of the /caption command.
type CaptionCommandHandler struct{}

// Handle handles the /caption command.
// /caption allows to set a caption to a forwarded message.
// By using /caption without any additional argument, one can delete the previous caption.
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
