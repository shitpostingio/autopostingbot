package updates

import (
	"github.com/shitpostingio/autopostingbot/documentstore/dbwrapper"
	"github.com/shitpostingio/autopostingbot/repository"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
)

const (
	tdlibDeletedMessageConversionFactor = 1048575
)

// handleNewDeletion handles deletion notifications and marks
// deleted channel posts as deleted in the database.
func handleNewDeletion(messages *client.UpdateDeleteMessages) {

	// We care only about permanent deletions in the channel
	if messages.ChatId != repository.Config.Autoposting.ChannelID || !messages.IsPermanent {
		return
	}

	log.Debugln("permanent deletions: ", messages.MessageIds)

	for _, id := range messages.MessageIds {

		id -= tdlibDeletedMessageConversionFactor
		err := dbwrapper.MarkPostAsDeletedByMessageID(id)
		if err != nil {
			log.Error(err)
		}

	}

}
