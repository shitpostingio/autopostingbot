package messages

//import (
//	"errors"
//	"fmt"
//
//	"gitlab.com/shitposting/autoposting-bot/edition"
//
//	"strings"
//
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
//
//	"gitlab.com/shitposting/autoposting-bot/database/database"
//	"gitlab.com/shitposting/autoposting-bot/repository"
//	"gitlab.com/shitposting/autoposting-bot/utility"
//
//	"gitlab.com/shitposting/telegram-markdown-processor/dbCaption"
//)
//
//const (
//	thanksCaption = "%s\n\n\\[Thanks to %s]"
//	creditCaption = "%s\n\n\\[By %s]"
//)
//
//// editCaption allows the user to edit the caption of a forwarded message or give the credit to the user.
//// It is used both by caption and credit command in the bot.
//func editCaption(msg *tgbotapi.Message, repo *repository.Repository, surplusCharacters int) (string, error) {
//
//	fileID, err := utility.GetFileIDFromMessage(msg.ReplyToMessage)
//
//	if err != nil {
//		return "", err
//	}
//
//	caption := dbCaption.PrepareCaptionForDB(msg.CommandArguments(), edition.ChannelName, utility.GetMessageEntities(msg), surplusCharacters)
//	database.UpdatePostCaptionByFileID(fileID, caption, repo.Db)
//	return "OK!\n" + caption, nil
//}
//
//// thankUser thanks the user from which the post is forwarded
//func thankUser(msg *tgbotapi.Message, repo *repository.Repository, surplusCharacters int) (string, error) {
//
//	// Checks if reply and Media
//	fileID, err := utility.GetFileIDFromMessage(msg.ReplyToMessage)
//	if err != nil {
//		return "", err
//	}
//
//	/* CREATE NEW CAPTION */
//	caption := createThanksCaption(msg)
//	if caption == "" {
//		return "", errors.New("nothing to thank")
//	}
//
//	caption = dbCaption.PrepareCaptionForDB(caption, edition.ChannelName, utility.GetMessageEntities(msg), surplusCharacters)
//	database.UpdatePostCaptionByFileID(fileID, caption, repo.Db)
//	return "OK!\n" + caption, nil
//}
//
//// creditCreator adds a link to the original content of the author
//// or credits the user from which the post is forwarded
//func creditCreator(msg *tgbotapi.Message, repo *repository.Repository, surplusCharacters int) (string, error) {
//
//	// Checks if reply and Media
//	fileID, err := utility.GetFileIDFromMessage(msg.ReplyToMessage)
//	if err != nil {
//		return "", err
//	}
//
//	var caption string
//	urlStart := strings.Index(strings.ToLower(msg.CommandArguments()), "http")
//
//	if urlStart != -1 {
//
//		/* WHO TO CREDIT: SUBSTRING UNTIL THE FIRST URL STARTS */
//		whoToCredit := msg.CommandArguments()[:urlStart]
//
//		/* ISOLATE THE LEFTOVER TEXT */
//		leftoverText := msg.CommandArguments()[urlStart:]
//
//		/* FIND THE END OF THE URL */
//		var url string
//		var comment string
//
//		urlEnd := strings.IndexAny(leftoverText, " \t\n")
//		if urlEnd != -1 {
//
//			/* URL: SUBSTRING OF leftoverText UNTIL urlEnd */
//			url = leftoverText[:urlEnd]
//
//			/* COMMENT: SUBSTRING OF leftoverText FROM urlEnd */
//			comment = leftoverText[urlEnd:]
//		} else {
//			url = leftoverText
//		}
//
//		/* CREATE NEW CAPTION */
//		caption = fmt.Sprintf("%s\n\n\\[By [%s](%s)]", strings.TrimSpace(comment), strings.TrimSpace(whoToCredit), url)
//		if urlEnd != -1 {
//			surplusCharacters += urlEnd + len(whoToCredit) + 1
//			caption = dbCaption.PrepareCaptionForDB(caption, edition.ChannelName, utility.GetMessageEntities(msg), surplusCharacters)
//		}
//
//	} else {
//
//		if msg.ReplyToMessage.ForwardFrom != nil && msg.ReplyToMessage.ForwardFrom.IsBot {
//			return "", errors.New("the message was forwarded from a bot and so it won't be credited")
//		}
//
//		/* CREATE NEW CAPTION */
//		caption = createCreditCaption(msg)
//		caption = dbCaption.PrepareCaptionForDB(caption, edition.ChannelName, utility.GetMessageEntities(msg), surplusCharacters)
//
//	}
//
//	if caption == "" {
//		return "", errors.New("nothing to credit")
//	}
//
//	database.UpdatePostCaptionByFileID(fileID, caption, repo.Db)
//	return "OK!\n" + caption, nil
//}
//
//// createThanksCaption creates the appropriate caption to thank the user
//func createThanksCaption(msg *tgbotapi.Message) (caption string) {
//
//	if msg.ReplyToMessage.ForwardFrom != nil {
//		caption = fmt.Sprintf(thanksCaption, msg.CommandArguments(), msg.ReplyToMessage.ForwardFrom.FirstName)
//	} else if msg.ReplyToMessage.ForwardSenderName != "" {
//		caption = fmt.Sprintf(thanksCaption, msg.CommandArguments(), msg.ReplyToMessage.ForwardSenderName)
//	} else {
//		caption = msg.CommandArguments()
//	}
//
//	return
//}
//
//// createCreditCaption creates the appropriate caption to credit the user
//func createCreditCaption(msg *tgbotapi.Message) (caption string) {
//
//	if msg.ReplyToMessage.ForwardFrom != nil {
//		caption = fmt.Sprintf(creditCaption, msg.CommandArguments(), msg.ReplyToMessage.ForwardFrom.FirstName)
//	} else if msg.ReplyToMessage.ForwardSenderName != "" {
//		caption = fmt.Sprintf(creditCaption, msg.CommandArguments(), msg.ReplyToMessage.ForwardSenderName)
//	} else {
//		caption = msg.CommandArguments()
//	}
//
//	return
//}
