package updates

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/analysisadapter"
	"gitlab.com/shitposting/autoposting-bot/api"
	caption "gitlab.com/shitposting/autoposting-bot/caption"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	l "gitlab.com/shitposting/autoposting-bot/localization"
	"gitlab.com/shitposting/autoposting-bot/posting"
)

func handleMedia(message *client.Message, mediatype string, skipDuplicateChecks bool) {

	//
	fileInfo, err := api.GetMediaFileInfo(message)
	if err != nil {
		log.Error("handleMedia: ", err)
		return
	}

	//
	fileInfo, err = api.DownloadFile(fileInfo.Id)
	if err != nil {
		log.Error("handleMedia: ", err)
		return
	}

	//
	fingerprint, err := analysisadapter.Request(fileInfo.Local.Path, mediatype, fileInfo.Remote.UniqueId)
	log.Debugln("Analysis response: ", fingerprint, err)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.ANALYSIS_NO_MEDIA_FINGERPRINT))
		return
	}

	//
	if !skipDuplicateChecks {

		post, err := dbwrapper.FindPostByUniqueID(fileInfo.Remote.UniqueId)
		if err != nil {
			post, err = dbwrapper.FindPostByFeatures(fingerprint.Histogram, fingerprint.PHash)
		}

		if err == nil {
			log.Debugln("Match found: ", post)
			formattedText, err := getDuplicateCaption(&post)
			if err != nil {
				_, _ = api.SendMedia(mediatype, message.ChatId, message.Id, post.Media.FileID, l.GetString(l.UPDATES_MEDIA_UNABLE_TO_GET_DUPLICATE_CAPTION), nil)
			} else {
				_, _ = api.SendMedia(mediatype, message.ChatId, message.Id, post.Media.FileID, formattedText.Text, formattedText.Entities)
			}

			return
		}

	}

	log.Debugln("Adding the post to the database")

	// Remove caption from forwarded posts
	var c string
	if message.ForwardInfo == nil {
		c = caption.ToHTMLCaption(api.GetMediaFormattedText(message))
	}

	//
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

	err = dbwrapper.AddPost(message.SenderUserId, media, c)
	if err != nil {
		log.Error(err)
	}

	//
	_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "Media added!")

	//
	if dbwrapper.GetQueueLength() == 1 {
		posting.ForcePostScheduling()
	}

}
