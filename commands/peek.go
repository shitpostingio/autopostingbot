package commands

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type PeekCommandHandler struct {}

func (PeekCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

	//
	nextPost, err := dbwrapper.GetNextPost()
	if err != nil {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, "Unable to find the next post. Is the queue empty?")
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
