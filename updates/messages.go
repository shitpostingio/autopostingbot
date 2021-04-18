package updates

import (
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/caption"
	"github.com/shitpostingio/autopostingbot/documentstore/dbwrapper"
	"github.com/shitpostingio/autopostingbot/repository"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
)

// handleNewMessage handles incoming messages.
func handleNewMessage(message *client.Message) {

	senderUserID, err := api.GetSenderUserID(message)
	if err != nil {
		log.Debugln("Message not from a user")
		return
	}

	// Tdlib delivers updates from self
	if senderUserID == repository.Me.Id {
		log.Debugln("Message from self")
		return
	}

	// Users need to be authorized to talk to the bot
	if !dbwrapper.UserIsAuthorized(senderUserID) {
		log.Println("Received message from unauthorized user with id ", senderUserID)
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

// handleUpdatedMessage handles updated messages.
func handleUpdatedMessage(umc *client.UpdateMessageContent) {

	// Since we are just given the updated content,
	// we need to get the full message
	message, err := api.GetMessage(umc.ChatId, umc.MessageId)
	if err != nil {
		log.Debugln("Unable to get message data")
		return
	}

	senderUserID, err := api.GetSenderUserID(message)
	if err != nil {
		log.Debugln("Message not from a user")
		return
	}

	// Tdlib delivers updates from self
	if senderUserID == repository.Me.Id {
		log.Debugln("Message from self")
		return
	}

	// Users need to be authorized to talk to the bot
	if !message.IsChannelPost && !dbwrapper.UserIsAuthorized(senderUserID) {
		log.Println("Received message from unauthorized user with id ", senderUserID)
		return
	}

	//
	log.Debugf("Message: %#v", message.Content)

	// Updates to textual messages can be handled normally, without any specific worry
	if umc.NewContent.MessageContentType() == client.TypeMessageText {

		if !message.IsChannelPost {
			handleText(message)
		}

		return
	}

	//
	fileInfo, err := api.GetMediaFileInfo(message)
	if err != nil {
		log.Error("handleUpdatedMessage: ", err)
		return
	}

	c := caption.ToHTMLCaption(api.GetMediaFormattedText(message))
	err = dbwrapper.UpdatePostCaptionByUniqueID(fileInfo.Remote.UniqueId, c)

	if err != nil {
		log.Debugln("handleUpdatedMessage:", err)
	}

}
