package manager

import (
	"errors"
	"fmt"
	"time"

	"gitlab.com/shitposting/telegram-markdown-processor/telegramCaption"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"

	log "github.com/sirupsen/logrus"

	"gitlab.com/shitposting/autoposting-bot/database/database"
	"gitlab.com/shitposting/autoposting-bot/edition"
	"gitlab.com/shitposting/autoposting-bot/types"
	"gitlab.com/shitposting/autoposting-bot/utility"
)

// post tries to post to the channel
func post() error {

	nextPost, err := database.GetNextPost(manager.db)
	if err != nil {

		/* MAKE SURE TO SCHEDULE NEXT POST EVEN IN CASE OF A FAILED POST */
		CalculateRateAndSchedulePosting(false)
		return err
	}

	return PostToChannel(nextPost, false)
}

// PostToChannel posts the input `Post` to the channel
func PostToChannel(post entities.Post, isPostNow bool) error {

	timeFromLastPost := -time.Until(manager.previousPostTime)

	if timeFromLastPost < 5*time.Minute {

		log.Warn(fmt.Sprintf("Posting too frequently is not allowed. Last post was %s ago", timeFromLastPost))

		if !manager.isTesting {
			return errors.New("the previous post was less than 5 minutes ago")
		}

		log.Info("Proceeding anyway because we're testing")

	}

	sentMessage, err := SendPostToChatID(post, manager.config.ChannelID, "", 0, true)
	if err == nil {

		err = database.MarkPostAsPosted(post, sentMessage.MessageID, manager.db)
		if err != nil {
			log.Error(err.Error())
		}

		manager.previousPostTime = time.Now()

	} else {

		_ = database.MarkPostAsFailed(post, manager.db)
		err = fmt.Errorf("%s on media with ID %s", err, post.FileID)

	}

	CalculateRateAndSchedulePosting(isPostNow)

	return err
}

// SendPostToChatID sends the input `Post` to the chatID, with an option to save the file to the disk
func SendPostToChatID(post entities.Post, chatID int64, caption string, replyToMessageID int, saveToDisk bool) (sentMessage tgbotapi.Message, err error) {

	if caption == "" {
		caption = prepareCaption(post)
	}

	/* SHARE MEDIA */
	switch post.TypeID {
	case types.Image:
		sentMessage, err = sharePhoto(chatID, post.FileID, caption, replyToMessageID, saveToDisk)
	case types.Video:
		sentMessage, err = shareVideo(chatID, post.FileID, caption, replyToMessageID, saveToDisk)
	case types.Animation:
		sentMessage, err = shareAnimation(chatID, post.FileID, caption, replyToMessageID, saveToDisk)
	}

	return
}

// PostToChannelByFileID posts the `Post` with the input fileID to the channel
func PostToChannelByFileID(fileID string, isPostNow bool) error {

	post := database.FindPostByFileID(fileID, manager.db)
	if post.ID == 0 {
		return errors.New("post not found")
	}

	return PostToChannel(post, isPostNow)
}

// SendPostByFileID sends the `Post` with the input fileID to the chatID, with an option to save the file to the disk
func SendPostByFileID(fileID string, chatID int64, caption string, saveToDisk bool, replyToMessageID int) (sentMessage tgbotapi.Message, err error) {

	post := database.FindPostByFileID(fileID, manager.db)
	if post.ID == 0 {
		return tgbotapi.Message{}, errors.New("post not found")
	}

	return SendPostToChatID(post, chatID, caption, replyToMessageID, saveToDisk)
}

// prepareCaption converts the markdown into HTML tags
func prepareCaption(post entities.Post) string {
	return telegramCaption.PrepareCaptionToSend(post.Caption, edition.ChannelName)
}

// sharePhoto tries to share a photo to a chatID
func sharePhoto(chatID int64, fileID, caption string, replyToMessageID int, saveToDisk bool) (sentMessage tgbotapi.Message, err error) {

	/* PREPARE PHOTO SHARE */
	image := tgbotapi.NewPhotoShare(chatID, fileID)
	image.ParseMode = "HTML"
	image.Caption = caption
	if replyToMessageID != 0 {
		image.ReplyToMessageID = replyToMessageID
	}

	/* SEND THE POST */
	sentMessage, err = manager.bot.Send(image)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to share image with fileID %s: %s. Another attempt will be made by uploading it", fileID, err.Error()))
		return sentMessage, err
	}

	/* WE MAY NOT WANT TO SAVE PICTURES TO DISK */
	if !saveToDisk || edition.IsSushiporn() {
		return sentMessage, err
	}

	/* SAVE IMAGES LOCALLY */
	url, err := manager.bot.GetFileDirectURL(fileID)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to get direct URL for image with fileID %s: %s", fileID, err.Error()))
		return sentMessage, err
	}

	/* DOWNLOAD */
	err = utility.DownloadFile(fmt.Sprintf("%s/%s.jpg", manager.config.MemePath, fileID), url)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to download image with fileID %s: %s", fileID, err.Error()))
	}

	return
}

// shareVideo tries to share a video to a chatID
func shareVideo(chatID int64, fileID, caption string, replyToMessageID int, saveToDisk bool) (sentMessage tgbotapi.Message, err error) {

	video := tgbotapi.NewVideoShare(chatID, fileID)
	video.ParseMode = "HTML"
	video.Caption = caption
	if replyToMessageID != 0 {
		video.ReplyToMessageID = replyToMessageID
	}

	sentMessage, err = manager.bot.Send(video)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to send image with fileID %s: %s", fileID, err.Error()))
	}

	/* WE MAY NOT WANT TO VIDEO PICTURES TO DISK */
	if !saveToDisk || edition.IsSushiporn() {
		return sentMessage, err
	}

	/* SAVE VIDEO LOCALLY */
	url, err := manager.bot.GetFileDirectURL(fileID)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to get direct URL for image with fileID %s: %s", fileID, err.Error()))
		return sentMessage, err
	}

	/* DOWNLOAD */
	err = utility.DownloadFile(fmt.Sprintf("%s/%s.mp4", manager.config.MemePath, fileID), url)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to download image with fileID %s: %s", fileID, err.Error()))
	}

	return
}

// shareAnimation tries to share a animation to a chatID
func shareAnimation(chatID int64, fileID, caption string, replyToMessageID int, saveToDisk bool) (sentMessage tgbotapi.Message, err error) {

	animation := tgbotapi.NewAnimationShare(chatID, fileID)
	animation.ParseMode = "HTML"
	animation.Caption = caption
	if replyToMessageID != 0 {
		animation.ReplyToMessageID = replyToMessageID
	}

	sentMessage, err = manager.bot.Send(animation)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to send image with fileID %s: %s", fileID, err.Error()))
	}

	/* WE MAY NOT WANT TO SAVE ANIMATION TO DISK */
	if !saveToDisk || edition.IsSushiporn() {
		return sentMessage, err
	}

	/* SAVE ANIMATION LOCALLY */
	url, err := manager.bot.GetFileDirectURL(fileID)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to get direct URL for image with fileID %s: %s", fileID, err.Error()))
		return sentMessage, err
	}

	/* DOWNLOAD */
	err = utility.DownloadFile(fmt.Sprintf("%s/%s.mp4", manager.config.MemePath, fileID), url)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to download image with fileID %s: %s", fileID, err.Error()))
	}

	return
}
