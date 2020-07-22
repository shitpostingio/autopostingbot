package messages

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"gitlab.com/shitposting/autoposting-bot/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"

	"gitlab.com/shitposting/autoposting-bot/utility"
)

const (
	statusText = "📋 Posts enqueued: %d\n🕜 Post rate: %s\n\n🔮 Next post in: %s (%s)"
)

// handleCommands handles the commands sent to the bot
func handleCommands(msg *tgbotapi.Message, user *entities.User, repo *repository.Repository) {

	command := strings.ToLower(msg.Command())
	if msg.ReplyToMessage == nil {
		handleNonReplyCommands(command, msg, repo)
		return
	}

	handleReplyCommands(command, msg, user, repo)
}

// handleNonReplyCommands handles the commands that don't require a reply
func handleNonReplyCommands(command string, msg *tgbotapi.Message, repo *repository.Repository) {

	var err error
	switch command {
	case "status":
		//_, err = handleStatusCommand(msg.Chat.ID, repo)
	case "pause":
		//_, err = handlePauseCommand(msg, repo)
	case "peek":
		//_, err = handlePeekCommand(msg, repo)
	default:
		err = errors.New("unsupported")
	}

	if err != nil {
		_, _ = repo.Bot.Send(tgbotapi.NewMessage(msg.Chat.ID, err.Error()))
	}
}

// handleReplyCommands handles the commands that require a reply
func handleReplyCommands(command string, msg *tgbotapi.Message, user *entities.User, repo *repository.Repository) {

	fileID, err := utility.GetFileIDFromMessage(msg.ReplyToMessage)
	if err != nil {
		_, _ = repo.Bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Unable to get fileID from message"))
		return
	}

	var replyText string

	switch command {
	case "delete":
		//replyText, err = handleDeleteCommand(fileID, user, repo)
	case "info":
		//replyText, err = handleInfoCommand(fileID, repo)
	case "preview":
		replyText, err = handlePreviewCommand(fileID, msg)
	case "postnow":
		replyText, err = handlePostNowCommand(msg, fileID, repo)
	case "caption":
		replyText, err = handleCaptionCommand(msg, repo)
	case "thanks", "thank":
		replyText, err = handleThankCommand(msg, repo)
	case "credit":
		replyText, err = handleCreditCommand(msg, repo)
	case "add":
		replyText, err = handleAddCommand(msg.ReplyToMessage, user, repo)
	default:
		err = errors.New("unsupported")
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "")
	reply.ReplyToMessageID = msg.MessageID
	if err != nil {
		reply.Text = err.Error()
		_, err = repo.Bot.Send(reply)
		if err != nil {
			log.Error(fmt.Sprintf("Error when sending error in reply to a command: %s", err.Error()))
		}
		return
	}

	reply.ParseMode = "HTML"
	if replyText != "" {
		reply.Text = replyText
		_, err = repo.Bot.Send(reply)
		if err != nil {
			log.Error(fmt.Sprintf("Error when sending reply to command: %s", err.Error()))
		}
	}
}
