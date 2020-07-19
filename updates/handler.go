package updates

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

func HandleUpdates(listener *client.Listener) {

	for update := range listener.Updates {

		if update.GetClass() == client.ClassUpdate {
			switch update.GetType() {
			case client.TypeUpdateNewMessage:
				handleNewMessage(update.(*client.UpdateNewMessage).Message)

			default:
				log.Printf("Type: %s, Value: %#v", update.GetType(), update)
			}

		}
	}

}

func handleNewMessage(message *client.Message) {

	//
	if !dbwrapper.UserIsAuthorized(message.SenderUserId) {
		log.Debugln("Ricevuto messaggio da utente non autorizzato: ", message.SenderUserId)
		return
	}

	log.Printf("Message: %#v", message.Content)

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
