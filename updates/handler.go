package updates

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	"gitlab.com/shitposting/autoposting-bot/repository"
)

const(
	tdlibdeletemessageidfuckery = 1048575
)

func HandleUpdates(listener *client.Listener) {

	for update := range listener.Updates {

		if update.GetClass() == client.ClassUpdate {
			switch update.GetType() {
			case client.TypeUpdateNewMessage:
				handleNewMessage(update.(*client.UpdateNewMessage).Message)
			case client.TypeUpdateDeleteMessages:
				handleNewDeletion(update.(*client.UpdateDeleteMessages))

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

func handleNewDeletion(messages *client.UpdateDeleteMessages) {

	// We care only about permanent deletions in the channel
	if messages.ChatId != repository.Config.ChannelID || !messages.IsPermanent {
		return
	}

	fmt.Println("permanent deletions: ", messages.MessageIds)

	for _, id := range messages.MessageIds {

		id -= tdlibdeletemessageidfuckery

		//TODO: vedere se usare id tdlib o botapi/client normali
		err := dbwrapper.MarkPostAsDeletedByMessageID(id)
		if err != nil {
			log.Error(err)
		}

	}

}
