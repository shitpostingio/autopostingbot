package messages

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.com/shitposting/telegram-markdown-processor/telegramCaption"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"
	log "github.com/sirupsen/logrus"

	"gitlab.com/shitposting/autoposting-bot/algo"
	"gitlab.com/shitposting/autoposting-bot/database/database"
	"gitlab.com/shitposting/autoposting-bot/edition"
	"gitlab.com/shitposting/autoposting-bot/manager"
	"gitlab.com/shitposting/autoposting-bot/media"
	"gitlab.com/shitposting/autoposting-bot/repository"
	"gitlab.com/shitposting/autoposting-bot/utility"

	"github.com/hako/durafmt"
)

func handleStatusCommand(chatID int64, repo *repository.Repository) (reply string, err error) {
	//err = sendStatus(chatID, repo)
	//if err != nil {
	//	log.Error(fmt.Sprintf("Unable to send status message: %s", err.Error()))
	//}

	return
}

func handlePauseCommand(msg *tgbotapi.Message, repo *repository.Repository) (reply string, err error) {

	var toParse string
	if strings.HasSuffix(msg.CommandArguments(), "h") {
		toParse = msg.CommandArguments()
	} else {
		toParse = msg.CommandArguments() + "h"
	}

	duration, _ := time.ParseDuration(toParse)
	err = manager.PausePosting(duration)
	if err != nil {
		return
	}

	log.Info(fmt.Sprintf("%s paused posting", utility.GetHandleOrName(msg.From)))
	return handleStatusCommand(msg.Chat.ID, repo)
}

func handlePeekCommand(msg *tgbotapi.Message, repo *repository.Repository) (reply string, err error) {
	//nextPost, err := database.GetNextPost(repo.Db)
	//_, err = manager.SendPostToChatID(nextPost, msg.Chat.ID, "", msg.MessageID, false)
	return
}

func handleDeleteCommand(fileID string, user *entities.User, repo *repository.Repository) (reply string, err error) {

	post := database.FindPostByFileID(fileID, repo.Db)
	if post.ID == 0 {
		err = errors.New("post not found")
		return
	}

	if post.UserID != user.ID {
		err = fmt.Errorf("you can't delete someone else's post: it was sent by @%s", post.User.Handle)
		return
	}

	if post.PostedAt != nil {
		err = fmt.Errorf("the post was already posted on %s", utility.FormatDate(*post.PostedAt))
		return
	}

	err = database.DeletePostByFileID(fileID, repo.Db)
	if err == nil {
		reply = "Deleted! 🚮"
	}

	return
}

func handleInfoCommand(fileID string, repo *repository.Repository) (reply string, err error) {

	post := database.FindPostByFileID(fileID, repo.Db)
	if post.ID == 0 {
		err = errors.New("post not found")
		return
	}

	if post.PostedAt != nil {
		reply = fmt.Sprintf("Post added by @%s on %s\nPosted on %s\nLink: t.me/%s/%d",
			post.User.Handle, utility.FormatDate(post.CreatedAt), utility.FormatDate(*post.PostedAt),
			edition.ChannelName, post.MessageID)
		return
	}

	position := database.GetQueuePositionByDatabaseID(post.ID, repo.Db)
	timeToPost := manager.GetNextPostTime().Add(algo.EstimatePostTime(position - 1))
	durationUntilPost := durafmt.Parse(time.Until(timeToPost).Truncate(time.Minute))

	reply = fmt.Sprintf("📋 The post is number %d in the queue\n👤 Added by @%s on %s\n\n🕜 It should be posted roughly in %s\n📅 On %s",
		position, post.User.Handle, utility.FormatDate(post.CreatedAt), durationUntilPost.String(), utility.FormatDate(timeToPost))
	return
}

func handlePreviewCommand(fileID string, msg *tgbotapi.Message) (reply string, err error) {

	_, err = manager.SendPostByFileID(fileID, msg.Chat.ID, "", false, 0)
	return
}

func handlePostNowCommand(msg *tgbotapi.Message, fileID string, repo *repository.Repository) (reply string, err error) {

	post := database.FindPostByFileID(fileID, repo.Db)
	if post.PostedAt != nil {
		err = fmt.Errorf("the post has been already posted")
		return
	}

	err = manager.PostToChannelByFileID(fileID, true)
	if err != nil {
		return
	}

	log.Info(fmt.Sprintf("%s used PostNow", utility.GetHandleOrName(msg.From)))
	return "Posted!", err
}

func handleCaptionCommand(msg *tgbotapi.Message, repo *repository.Repository) (reply string, err error) {

	reply, err = editCaption(msg, repo, len(msg.Command())+2)
	reply = telegramCaption.PrepareCaptionToSend(reply, edition.ChannelName)
	return
}

func handleThankCommand(msg *tgbotapi.Message, repo *repository.Repository) (reply string, err error) {

	/* FORWARDS NOT HIDDEN */
	if msg.ReplyToMessage.ForwardFrom != nil {

		if !msg.ReplyToMessage.ForwardFrom.IsBot {
			reply, err = thankUser(msg, repo, len(msg.Command())+2)
			reply = telegramCaption.PrepareCaptionToSend(reply, edition.ChannelName)
		} else {
			reply = "The message was forwarded from a bot and so it won't be thanked"
		}

		return
	}

	reply, err = thankUser(msg, repo, len(msg.Command())+2)
	reply = telegramCaption.PrepareCaptionToSend(reply, edition.ChannelName)
	return
}

func handleCreditCommand(msg *tgbotapi.Message, repo *repository.Repository) (reply string, err error) {

	/* FORWARDS NOT HIDDEN */
	if msg.ReplyToMessage.ForwardFrom != nil {
		reply, err = creditCreator(msg, repo, len(msg.Command())+2)
		reply = telegramCaption.PrepareCaptionToSend(reply, edition.ChannelName)
		return
	}

	reply, err = creditCreator(msg, repo, len(msg.Command())+2)
	reply = telegramCaption.PrepareCaptionToSend(reply, edition.ChannelName)
	return
}

func handleAddCommand(msg *tgbotapi.Message, user *entities.User, repo *repository.Repository) (reply string, err error) {

	if msg == nil {
		err = errors.New("no reply message")
		return
	}

	switch {
	case msg.Photo != nil:
		reply, _ = media.HandleNewPhoto(msg, user, repo, false)
	case msg.Video != nil:
		reply, _ = media.HandleNewVideo(msg, user, repo, false)
	case msg.Animation != nil:
		reply, _ = media.HandleNewAnimation(msg, user, repo, false)
	}
	return
}

func sendStatus(chatID int64, repo *repository.Repository) error {

	nextPost := manager.GetNextPostTime()
	message := fmt.Sprintf(statusText, database.GetQueueLength(repo.Db), manager.GetPostingRate().String(),
		time.Until(nextPost).Truncate(time.Minute), nextPost.Format("15:04"))

	_, err := repo.Bot.Send(tgbotapi.NewMessage(chatID, message))
	return err
}
