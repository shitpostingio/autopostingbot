package commands

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/posting"
	"gitlab.com/shitposting/autoposting-bot/telegram"
	"strconv"
	"strings"
	"time"
)

type PauseCommandHandler struct {}

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
	posting.RequestPause(duration)

	//
	var whoPaused string
	user, err := api.GetUserByID(message.SenderUserId)
	if err != nil {
		whoPaused = telegram.GetNameFromUser(user)
	} else {
		whoPaused =  strconv.Itoa(int(message.SenderUserId))
	}
	//

	log.Info(fmt.Sprintf("%s paused posting", whoPaused))
	return StatusCommandHandler{}.Handle("", message, replyToMessage)

}
