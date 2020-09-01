package updates

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	"gitlab.com/shitposting/autoposting-bot/repository"
)

func handleNewMessage(message *client.Message) {

	// Tdlib delivers updates from self
	if message.SenderUserId == repository.Me.Id {
		log.Debugln("Message from self")
		return
	}

	// Users need to be authorized to talk to the bot
	if !dbwrapper.UserIsAuthorized(message.SenderUserId) {
		log.Println("Received message from unauthorized user with id ", message.SenderUserId)
		return
	}

	//
	log.Debugf("Message: %#v", message.Content)

	//
	switch message.Content.MessageContentType() {
	case client.TypeMessageText:
		handleText(message)
	case client.TypeMessageAnimation:
		handleMedia(message, client.TypeAnimation, false)
	case client.TypeMessagePhoto:
		handleMedia(message, client.TypePhoto, false)
	case client.TypeMessageVideo:
		handleMedia(message, client.TypeVideo, false)
	}

}
