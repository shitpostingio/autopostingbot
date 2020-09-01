package commands

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	l "gitlab.com/shitposting/autoposting-bot/localization"
)

type PeekCommandHandler struct {}

func (PeekCommandHandler) Handle(_ string, message, _ *client.Message) error {

	//
	nextPost, err := dbwrapper.GetNextPost()
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.COMMANDS_PEEK_NO_POST_FOUND))
		return err
	}

	//
	ft, err := api.GetFormattedText(nextPost.Caption)
	if err != nil {
		ft = &client.FormattedText{
			Text:     nextPost.Caption,
			Entities: nil,
		}
	}

	//
	_, err = api.SendMedia(nextPost.Media.Type, message.ChatId, message.Id, nextPost.Media.FileID, ft.Text, ft.Entities)
	return err

}
