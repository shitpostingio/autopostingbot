package commands

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/manager"
	"strings"
	"time"
)

type PauseCommandHandler struct {
}

func (PauseCommandHandler) Handle(arguments string, message *client.Message) error {

	var toParse string
	if strings.HasSuffix(arguments, "h") {
		toParse = arguments
	} else {
		toParse = arguments + "h"
	}

	duration, _ := time.ParseDuration(toParse)
	err := manager.PausePosting(duration)
	if err != nil {
		return err
	}

	// TODO: SISTEMARE PRINT
	log.Info(fmt.Sprintf("%d paused posting", message.SenderUserId))
	return StatusCommandHandler{}.Handle("", message)

}
