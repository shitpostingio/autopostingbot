package commands

import (
	"fmt"
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/documentstore/dbwrapper"
	l "github.com/shitpostingio/autopostingbot/localization"
	"github.com/shitpostingio/autopostingbot/posting"
	"github.com/zelenin/go-tdlib/client"
	"time"
)

// StatusCommandHandler represents the handler of the /status command.
type StatusCommandHandler struct{}

// Handle handles the /status command.
// /status returns information about the posts enqueued, the posting rate
// and the time until the next post.
func (StatusCommandHandler) Handle(_ string, message, _ *client.Message) error {

	//
	nextPost := posting.GetNextPostTime()
	queueLength := dbwrapper.GetQueueLength()
	postingRate := posting.GetPostingRate().String()
	minutesUntilNextPost := time.Until(nextPost).Truncate(time.Minute)

	//
	text := fmt.Sprintf(l.GetString(l.COMMANDS_STATUS_POSTS_ENQUEUED),
		queueLength,
		postingRate,
		minutesUntilNextPost,
		nextPost.Format("15:04"))

	//
	_, err := api.SendPlainText(message.ChatId, text)
	return err

}
