package commands

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	l "gitlab.com/shitposting/autoposting-bot/localization"
	"gitlab.com/shitposting/autoposting-bot/posting"
	"gitlab.com/shitposting/autoposting-bot/telegram"
	"strconv"
	"strings"
	"time"
)

type PauseCommandHandler struct{}

func (PauseCommandHandler) Handle(arguments string, message, replyToMessage *client.Message) error {

	//
	var toParse string
	if arguments == "" {
		toParse = "1h"
	} else if strings.HasSuffix(arguments, "h") {
		toParse = arguments
	} else {
		toParse = arguments + "h"
	}

	//
	duration, err := time.ParseDuration(toParse)
	if err != nil {
		duration = 1 * time.Hour
	}

	//
	err = posting.RequestPause(duration)
	if err != nil {
		reply := fmt.Sprintf(l.GetString(l.COMMANDS_PAUSE_UNSUCCESSFUL), err)
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, reply)
		return err
	}

	//
	var whoPaused string
	user, err := api.GetUserByID(message.SenderUserId)
	if err != nil {
		whoPaused = telegram.GetNameFromUser(user)
	} else {
		whoPaused = strconv.Itoa(int(message.SenderUserId))
	}
	//

	log.Info(fmt.Sprintf("%s paused posting", whoPaused))
	return StatusCommandHandler{}.Handle("", message, replyToMessage)

}
