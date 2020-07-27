package commands

import (
	"errors"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/analysisadapter"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/caption"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
)

type AddCommandHandler struct {}

func (AddCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

	//
	if replyToMessage == nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "This command needs to be used in reply to a media file")
		return errors.New("reply to message nil")
	}

	//
	fileInfo, err := api.GetMediaFileInfo(replyToMessage)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "This command needs to be used in reply to a media file")
		return err
	}

	//
	mediaType := api.GetTypeFromMessageType(replyToMessage.Content.MessageContentType())
	fingerprint, err := analysisadapter.Request(fileInfo.Local.Path, mediaType, fileInfo.Remote.UniqueId)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "Unable to get media fingerprint")
		return err
	}

	//
	avg, sum := entities.GetHistogramAverageAndSum(fingerprint.Histogram)

	//
	media := entities.Media{
		Type:             mediaType,
		TdlibID:          fileInfo.Id,
		FileUniqueID:     fileInfo.Remote.UniqueId,
		FileID:           fileInfo.Remote.Id,
		Histogram:        fingerprint.Histogram,
		HistogramAverage: avg,
		HistogramSum:     sum,
		PHash:            fingerprint.PHash,
	}

	// If the message is a forward, remove the caption
	postCaption := ""

	// if the message is not a forward, get the caption
	if message.ForwardInfo == nil {
		ft := api.GetMediaFormattedText(replyToMessage)
		postCaption = caption.ToHTMLCaption(ft)
	}

	err = dbwrapper.AddPost(replyToMessage.SenderUserId, media, postCaption)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "Error while trying to add the post")
	} else {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "Media added correctly!")
	}

	return err

}
