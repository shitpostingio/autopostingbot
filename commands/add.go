package commands

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/analysisadapter"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
)

type AddCommandHandler struct {
}

func (AddCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

	fileInfo, err := api.GetMediaFileInfo(replyToMessage)
	if err != nil {
		return err
	}

	mediaType := api.GetTypeFromMessageType(replyToMessage.Content.MessageContentType())

	fingerprint, err := analysisadapter.Request(fileInfo.Local.Path, mediaType, fileInfo.Remote.UniqueId)
	if err != nil {
		return err
	}

	avg, sum := entities.GetHistogramAverageAndSum(fingerprint.Histogram)

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

	//TODO: VEDERE CHE FARE CON LA CAPTION
	err = dbwrapper.AddPost(replyToMessage.SenderUserId, media, "")
	return err

}
