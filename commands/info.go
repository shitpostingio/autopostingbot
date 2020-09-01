package commands

import (
	"errors"
	"fmt"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	l "gitlab.com/shitposting/autoposting-bot/localization"
	"gitlab.com/shitposting/autoposting-bot/posting"
	"gitlab.com/shitposting/autoposting-bot/telegram"
	"gitlab.com/shitposting/autoposting-bot/utility"
	"strconv"
	"time"
)

type InfoCommandHandler struct{}

func (InfoCommandHandler) Handle(_ string, message, replyToMessage *client.Message) error {

	//
	if replyToMessage == nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.COMMANDS_REPLY_TO_MEDIA_FILE))
		return errors.New("reply to message nil")
	}

	//
	fi, err := api.GetMediaFileInfo(replyToMessage)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.COMMANDS_REPLY_TO_MEDIA_FILE))
		return err
	}

	//
	post, err := dbwrapper.FindPostByUniqueID(fi.Remote.UniqueId)
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.DATABASE_UNABLE_TO_FIND_POST))
		return err
	}

	//
	log.Debugln("Found post: ", post)

	//
	var name string
	user, err := api.GetUserByID(post.AddedBy)
	if err != nil {
		name = strconv.Itoa(int(post.AddedBy))
	} else {
		name = telegram.GetNameFromUser(user)
	}

	//
	var reply string
	if post.PostedAt != nil {

		//
		reply = fmt.Sprintf(l.GetString(l.COMMANDS_INFO_ALREADY_POSTED),
			post.AddedBy, name, utility.FormatDate(post.AddedAt), utility.FormatDate(*post.PostedAt), posting.GetPostingManager().GetEditionName(), post.MessageID)

		//
		ft, err := api.GetFormattedText(reply)
		if err != nil {
			ft = &client.FormattedText{
				Text:     reply,
				Entities: nil,
			}
		}

		_, _ = api.SendText(replyToMessage.ChatId, replyToMessage.Id, ft.Text, ft.Entities)
		return err

	}

	//
	position := dbwrapper.GetQueuePositionByAddTime(post.AddedAt)
	timeToPost := posting.GetNextPostTime().Add(posting.GetPostingManager().EstimatePostTime(position - 1))
	durationUntilPost := durafmt.Parse(time.Until(timeToPost).Truncate(time.Minute))
	reply = fmt.Sprintf(l.GetString(l.COMMANDS_INFO_NOT_YET_POSTED),
		position, post.AddedBy, name, utility.FormatDate(post.AddedAt), durationUntilPost.String(), utility.FormatDate(timeToPost))

	//
	ft, err := api.GetFormattedText(reply)
	if err != nil {
		ft = &client.FormattedText{
			Text:     reply,
			Entities: nil,
		}
	}

	//
	_, _ = api.SendText(replyToMessage.ChatId, replyToMessage.Id, ft.Text, ft.Entities)
	return err

}
