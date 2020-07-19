package updates

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	"gitlab.com/shitposting/autoposting-bot/edition"
	"gitlab.com/shitposting/autoposting-bot/legacy"
	"gitlab.com/shitposting/autoposting-bot/telegram"
	"gitlab.com/shitposting/autoposting-bot/utility"
)

func getDuplicateCaption(duplicatePost *entities.Post) (*client.FormattedText, error) {

	user, err := api.GetUserByID(duplicatePost.AddedBy)
	if err != nil {
		//TODO: pensare a qualcosa
	}

	//TODO: USARE GOTRANS E USARE STRINGHE DI TRADUZIONE
	caption := fmt.Sprintf(
		"🚨 Duplicate detected! 🚨\n\nFirst added by <a href=\"tg://user?id=%d\">%s</a>\non %s",
		duplicatePost.AddedBy, telegram.GetNameFromUser(user), utility.FormatDate(duplicatePost.AddedAt))

	if duplicatePost.MessageID != 0 {
		caption = fmt.Sprintf("%s\nPosted on %s\nLink: t.me/%s/%d", caption, utility.FormatDate(*duplicatePost.PostedAt), edition.ChannelName, duplicatePost.MessageID)
	}

	ft, err := legacy.NewFormattedTextFromCaption(caption)
	return ft, err

}


