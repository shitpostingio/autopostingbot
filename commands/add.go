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

func (AddCommandHandler) Handle(arguments string, message *client.Message) error {

	fileInfo, err := api.GetMediaFileInfo(message)
	if err != nil {
		return err
	}

	mediaType := api.GetTypeFromMessageType(message.Content.MessageContentType())

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
	err = dbwrapper.AddPost(message.SenderUserId, media, "")
	return err

}
