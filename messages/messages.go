package messages

//import (
//	"fmt"
//
//	log "github.com/sirupsen/logrus"
//
//	"gitlab.com/shitposting/autoposting-bot/edition"
//	"gitlab.com/shitposting/autoposting-bot/manager"
//	"gitlab.com/shitposting/autoposting-bot/types"
//
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
//	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"
//	"gitlab.com/shitposting/telegram-markdown-processor/telegramCaption"
//
//	"gitlab.com/shitposting/autoposting-bot/database/database"
//	"gitlab.com/shitposting/autoposting-bot/media"
//	"gitlab.com/shitposting/autoposting-bot/repository"
//	"gitlab.com/shitposting/autoposting-bot/utility"
//)
//
//// HandleNew handles new `Message`
//func HandleNew(msg *tgbotapi.Message, repo *repository.Repository) {
//
//	user := database.GetUserByTelegramID(msg.From.ID, repo.Db)
//	if !database.UserIsAuthorized(user) {
//		return
//	}
//
//	var replyText string
//	var duplicatePost entities.Post
//	switch {
//	case msg.Photo != nil:
//		//replyText, duplicatePost = media.HandleNewPhoto(msg, &user, repo, true)
//	case msg.Video != nil:
//		//replyText, duplicatePost = media.HandleNewVideo(msg, &user, repo, true)
//	case msg.Animation != nil:
//		//replyText, duplicatePost = media.HandleNewAnimation(msg, &user, repo, true)
//	case msg.Text != "" && msg.IsCommand():
//		handleCommands(msg, &user, repo)
//		return
//	}
//
//	if duplicatePost.ID != 0 {
//		sendDuplicateAlert(&duplicatePost, msg.Chat.ID, msg.MessageID, repo)
//		return
//	}
//
//	queueLength := database.GetQueueLength(repo.Db)
//	if queueLength == 1 {
//		manager.CalculateRateAndSchedulePosting(true)
//	}
//
//	if replyText == "" {
//		return
//	}
//
//	reply := tgbotapi.NewMessage(msg.Chat.ID, replyText)
//	reply.ReplyToMessageID = msg.MessageID
//	_, err := repo.Bot.Send(reply)
//	if err != nil {
//		log.Error(fmt.Sprintf("Unable to send reply to new message: %s", err.Error()))
//	}
//}
//
//// HandleEdited handles edited `Message`
//func HandleEdited(msg *tgbotapi.Message, repo *repository.Repository) {
//
//	user := database.GetUserByTelegramID(msg.From.ID, repo.Db)
//	if !database.UserIsAuthorized(user) {
//		return
//	}
//
//	var replyText string
//	switch {
//	case msg.Photo != nil:
//		replyText = media.HandleEditedPhoto(msg, repo)
//	case msg.Video != nil:
//		replyText = media.HandleEditedVideo(msg, repo)
//	case msg.Animation != nil:
//		replyText = media.HandleEditedAnimation(msg, repo)
//	case msg.Text != "" && msg.IsCommand():
//		handleCommands(msg, &user, repo)
//	}
//
//	if replyText == "" {
//		return
//	}
//
//	reply := tgbotapi.NewMessage(msg.Chat.ID, replyText)
//	reply.ReplyToMessageID = msg.MessageID
//	_, err := repo.Bot.Send(reply)
//	if err != nil {
//		log.Error(fmt.Sprintf("Unable to send reply to edited message: %s", err.Error()))
//	}
//}
//
//// sendDuplicateAlert replies to a message with the duplicate media
//func sendDuplicateAlert(duplicatePost *entities.Post, chatID int64, replyID int, repo *repository.Repository) {
//
//	caption := fmt.Sprintf("🚨 Duplicate detected! 🚨\n\nFirst added by @%s on %s", duplicatePost.User.Handle, utility.FormatDate(duplicatePost.CreatedAt))
//
//	if duplicatePost.MessageID != 0 {
//		caption = fmt.Sprintf("%s\nPosted on %s\nLink: t.me/%s/%d", caption, utility.FormatDate(*duplicatePost.PostedAt), edition.ChannelName, duplicatePost.MessageID)
//	}
//
//	_, err := manager.SendPostToChatID(*duplicatePost, chatID, caption, replyID, false)
//	if err == nil {
//		return
//	}
//
//	err = tryUploading(duplicatePost, chatID, replyID, caption, repo)
//	if err != nil {
//		log.Error(fmt.Sprintf("Unable to send duplicate message: %s", err.Error()))
//	}
//}
//
//// tryUploading tries to upload a photo that can't be shared
//func tryUploading(duplicatePost *entities.Post, chatID int64, replyID int, caption string, repo *repository.Repository) error {
//
//	if duplicatePost.TypeID == types.Image {
//		photoConfig := tgbotapi.NewPhotoUpload(chatID, fmt.Sprintf("%s/%s.jpg", repo.Config.MemePath, duplicatePost.FileID))
//		photoConfig.ReplyToMessageID = replyID
//		photoConfig.Caption = telegramCaption.PrepareCaptionToSend(caption, edition.ChannelName)
//		photoConfig.ParseMode = "HTML"
//		_, err := repo.Bot.Send(photoConfig)
//		if err == nil {
//			return nil
//		}
//	}
//
//	_, err := repo.Bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("---Unable to send duplicate---\n\n%s", caption)))
//	return err
//}
