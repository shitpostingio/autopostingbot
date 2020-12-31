package commands

import (
	"fmt"
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/documentstore/dbwrapper"
	l "github.com/shitpostingio/autopostingbot/localization"
	"github.com/shitpostingio/autopostingbot/posting"
	"github.com/zelenin/go-tdlib/client"
	"strings"
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
	caption := nextPost.Caption
	if posting.GetChannelHandle() != "" && !strings.Contains(nextPost.Caption, "@"+posting.GetChannelHandle()) {
		caption = fmt.Sprintf("%s\n\n@%s", caption, posting.GetChannelHandle())
	}

	ft, err := api.GetFormattedText(caption)
	if err != nil {
		ft = &client.FormattedText{
			Text:     caption,
			Entities: nil,
		}
	}

	//
	_, err = api.ShareMedia(nextPost.Media.Type, message.ChatId, message.Id, nextPost.Media.FileID, ft.Text, ft.Entities)
	return err

}
