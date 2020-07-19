package updates

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/analysisadapter"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	"gitlab.com/shitposting/autoposting-bot/files"
)

func handlePhoto(message *client.Message) {

	//
	fileInfo, err := api.GetMediaFileInfo(message)
	if err != nil {
		log.Error("handlePhoto: ", err)
		return
	}

	//
	fileInfo, err = files.DownloadFile(fileInfo.Id)
	if err != nil {
		log.Error("handlePhoto: ", err)
		return
	}

	fingerprint, err := analysisadapter.Request(fileInfo.Local.Path, client.TypePhoto, fileInfo.Remote.UniqueId)
	log.Println("Ottenuta risposta da analysis: ", fingerprint, err)
	if err != nil {
		//TODO: RESTITUIRE ERRORE
		return
	}

	post, err := dbwrapper.FindPostByFeatures(fingerprint.Histogram, fingerprint.PHash)
	if err == nil {
		//TODO: SEND DUPLICATE
		log.Println("Match found: ", post)
		_, _ = api.SendMedia(client.TypePhoto, message.ChatId, post.Media.FileID, "DUPLICATE DETECTED", nil)
		return
	}

	log.Println("AGGIUNGO IL POST AL DB")
	avg, sum := entities.GetHistogramAverageAndSum(fingerprint.Histogram)
	media := entities.Media{
		Type:             client.TypePhoto,
		TdlibID:          fileInfo.Id,
		FileUniqueID:     fileInfo.Remote.UniqueId,
		FileID:           fileInfo.Remote.Id,
		Histogram:        fingerprint.Histogram,
		HistogramAverage: avg,
		HistogramSum:     sum,
		PHash:            fingerprint.PHash,
	}

	photoMessage := message.Content.(*client.MessagePhoto)
	err = dbwrapper.AddPost(message.SenderUserId, media, photoMessage.Caption)
	if err != nil {
		log.Error(err)
	}

	_, _ = api.SendPlainText(message.ChatId, "Photo added!")

}
