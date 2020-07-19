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
		log.Println("Ricevuto messaggio da utente non autorizzato: ", message.SenderUserId)
		return
	}

	log.Printf("Message: %#v", message.Content)

	switch message.Content.MessageContentType() {
	case client.TypeMessageText:
		handleText(message)
	case client.TypeMessageAnimation:

	case client.TypeMessagePhoto:
		handlePhoto(message)
	case client.TypeMessageVideo:

	}

}
