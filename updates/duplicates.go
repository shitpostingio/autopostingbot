package updates

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	l "gitlab.com/shitposting/autoposting-bot/localization"
	"gitlab.com/shitposting/autoposting-bot/repository"
	"gitlab.com/shitposting/autoposting-bot/telegram"
	"gitlab.com/shitposting/autoposting-bot/utility"
	"strconv"
)

func getDuplicateCaption(duplicatePost *entities.Post) (*client.FormattedText, error) {

	var userName string
	user, err := api.GetUserByID(duplicatePost.AddedBy)
	if err != nil {
		userName = strconv.Itoa(int(duplicatePost.AddedBy))
	} else {
		userName = telegram.GetNameFromUser(user)
	}

	caption := fmt.Sprintf(
		l.GetString(l.UPDATES_DUPLICATES_DUPLICATE_ADDED_BY),
		duplicatePost.AddedBy, userName, utility.FormatDate(duplicatePost.AddedAt))

	if duplicatePost.MessageID != 0 {

		link, err := api.GetMessageLink(repository.Config.Autoposting.ChannelID, duplicatePost.MessageID)
		if err != nil {
			link = l.GetString(l.UPDATES_DUPLICATE_LINK_UNAVAILABLE)
		}

		captionEnd := fmt.Sprintf(l.GetString(l.UPDATES_DUPLICATE_DUPLICATE_ADDED_AT), utility.FormatDate(*duplicatePost.PostedAt), link)
		caption = fmt.Sprintf("%s\n%s", caption, captionEnd)

	}

	ft, err := api.GetFormattedText(caption)
	return ft, err

}
