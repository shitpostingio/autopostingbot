package updates

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
)

// HandleUpdates handles incoming updates, dispatching them
// to the appropriate sub-handlers.
func HandleUpdates(listener *client.Listener) {

	for update := range listener.Updates {

		if update.GetClass() == client.ClassUpdate {

			switch update.GetType() {
			case client.TypeUpdateNewMessage:
				handleNewMessage(update.(*client.UpdateNewMessage).Message)
			case client.TypeUpdateDeleteMessages:
				handleNewDeletion(update.(*client.UpdateDeleteMessages))
			default:
				log.Debugf("Type: %s, Value: %#v", update.GetType(), update)
			}

		}

	}

}
