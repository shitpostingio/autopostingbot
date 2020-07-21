package updates

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/analysisadapter"
	"gitlab.com/shitposting/autoposting-bot/api"
	caption "gitlab.com/shitposting/autoposting-bot/caption"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	"gitlab.com/shitposting/autoposting-bot/files"
)

func handleMedia(message *client.Message, mediatype string, skipDuplicateChecks bool) {

	//
	fileInfo, err := api.GetMediaFileInfo(message)
	if err != nil {
		log.Error("handleMedia: ", err)
		return
	}

	//
	fileInfo, err = files.DownloadFile(fileInfo.Id)
	if err != nil {
		log.Error("handleMedia: ", err)
		return
	}

	fingerprint, err := analysisadapter.Request(fileInfo.Local.Path, mediatype, fileInfo.Remote.UniqueId)
	log.Println("Ottenuta risposta da analysis: ", fingerprint, err)
	if err != nil {
		//TODO: RESTITUIRE ERRORE
		return
	}

	if !skipDuplicateChecks {

		post, err := dbwrapper.FindPostByFeatures(fingerprint.Histogram, fingerprint.PHash)
		if err == nil {
			log.Println("Match found: ", post)
			//TODO: VEDERE L'ERRORE
			formattedText, _ := getDuplicateCaption(&post)
			_, _ = api.SendMedia(mediatype, message.ChatId, message.Id, post.Media.FileID, formattedText.Text, formattedText.Entities)
			return
		}

	}

	log.Println("AGGIUNGO IL POST AL DB")
	avg, sum := entities.GetHistogramAverageAndSum(fingerprint.Histogram)
	media := entities.Media{
		Type:             mediatype,
		TdlibID:          fileInfo.Id,
		FileUniqueID:     fileInfo.Remote.UniqueId,
		FileID:           fileInfo.Remote.Id,
		Histogram:        fingerprint.Histogram,
		HistogramAverage: avg,
		HistogramSum:     sum,
		PHash:            fingerprint.PHash,
	}

	c := caption.ToHTMLCaption(api.GetMediaFormattedText(message))
	err = dbwrapper.AddPost(message.SenderUserId, media, c)
	if err != nil {
		log.Error(err)
	}

	_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "Media added!")

}
