package updates

import (
	"github.com/shitpostingio/autopostingbot/analysisadapter"
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/caption"
	"github.com/shitpostingio/autopostingbot/documentstore/dbwrapper"
	"github.com/shitpostingio/autopostingbot/documentstore/entities"
	l "github.com/shitpostingio/autopostingbot/localization"
	"github.com/shitpostingio/autopostingbot/posting"
	"github.com/shitpostingio/autopostingbot/repository"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
)

// handleMedia handles incoming media messages.
// It checks for duplicates and adds them to the database if they are unique.
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

		// Only in case we found a duplicate
		if err == nil {
			log.Debugln("Match found: ", post)
			_ = sendDuplicate(&post, message)
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

func sendDuplicate(post *entities.Post, message *client.Message) error {

	fileID := post.Media.FileID
	path := ""

	// If the post has already been sent on the channel, we can get the message
	// and download it. By doing this, we can send it more reliably.
	if post.MessageID != 0 {

		oldMessage, err := api.GetMessage(repository.Config.Autoposting.ChannelID, post.MessageID)
		if err == nil {

			//
			oldFileInfo, err := api.GetMediaFileInfo(oldMessage)
			if err != nil {
				log.Error("sendDuplicate: ", err)
				return err
			}

			//
			oldFileInfo, err = api.DownloadFile(oldFileInfo.Id)
			if err != nil {
				log.Error("sendDuplicate: ", err)
				return err
			}

			fileID = oldFileInfo.Remote.Id
			path = oldFileInfo.Local.Path

		}

	}

	formattedText, err := getDuplicateCaption(post)
	if err != nil {
		formattedText = &client.FormattedText{
			Text:     l.GetString(l.UPDATES_MEDIA_UNABLE_TO_GET_DUPLICATE_CAPTION),
			Entities: nil,
		}
	}

	_, err = api.SendMedia(post.Media.Type, message.ChatId, message.Id, fileID, path, formattedText.Text, formattedText.Entities)
	if err != nil {
		_, err = api.SendText(message.ChatId, message.Id, formattedText.Text, formattedText.Entities)
	}

	return err

}
