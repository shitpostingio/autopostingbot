package commands

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"github.com/shitpostingio/autopostingbot/api"
	l "github.com/shitpostingio/autopostingbot/localization"
	"github.com/shitpostingio/autopostingbot/posting"
	"github.com/shitpostingio/autopostingbot/telegram"
	"strconv"
	"strings"
	"time"
)

// PauseCommandHandler represents the handler of the /pause command.
type PauseCommandHandler struct{}

// Handle handles the /pause command.
// /pause pauses the posting for a certain amount of hours, 1 by default.
// The posting manager will not allow pauses too close between each other,
// as to prevent accidental pauses.
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
