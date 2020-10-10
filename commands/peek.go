package commands

import (
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/documentstore/dbwrapper"
	l "github.com/shitpostingio/autopostingbot/localization"
	"github.com/zelenin/go-tdlib/client"
)

// PeekCommandHandler represents the handler of the /peek command.
type PeekCommandHandler struct{}

// Handle handles the /peek command.
// /peek returns the first post in the queue, along with its caption.
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
