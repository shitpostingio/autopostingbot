package commands

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	"gitlab.com/shitposting/autoposting-bot/posting"
	"time"
)

const (
	//TODO: ASSOLUTAMENTE NON QUA, POSSIBILMENTE STRINGA TRADOTTA
	statusText = "📋 Posts enqueued: %d\n🕜 Post rate: %s\n\n🔮 Next post in: %s (%s)"
)

type StatusCommandHandler struct {
}

func (StatusCommandHandler) Handle(arguments string, message *client.Message) error {

	nextPost := posting.GetNextPostTime()
	text := fmt.Sprintf(statusText,
		dbwrapper.GetQueueLength(),
		posting.GetPostingRate().String(),
		time.Until(nextPost).Truncate(time.Minute),
		nextPost.Format("15:04"))

	_, err := api.SendPlainText(message.ChatId, text)
	return err

}
