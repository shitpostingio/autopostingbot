package api

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
)

var (

	sendFunctions = map[string]func(int64, int64, string, string, []*client.TextEntity) (*client.Message, error){
		client.TypeAnimation: SendAnimation,
		client.TypePhoto:     SendPhoto,
		client.TypeVideo:     SendVideo,
	}

	fileInfoFunctions = map[string]func(*client.Message) *client.File{
		client.TypeMessageAnimation: GetAnimationFileInfoFromMessage,
		client.TypeMessagePhoto:     GetPhotoFileInfoFromMessage,
		client.TypeMessageVideo:     GetVideoFileInfoFromMessage,
	}

)

func SendMedia(mediaType string, chatID, replyToMessageID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

	send, found := sendFunctions[mediaType]
	if !found {
		return nil, fmt.Errorf("send function not found for media type %s", mediaType)
	}

	return send(chatID, replyToMessageID, remoteFileID, caption, entities)

}

func GetMediaFileInfo(message *client.Message) (*client.File, error) {

	//
	mediaType := message.Content.MessageContentType()

	//
	getIDs, found := fileInfoFunctions[mediaType]
	if !found {
		return nil, fmt.Errorf("get file info function not found for media type %s", mediaType)
	}

	return getIDs(message), nil

}

func GetMediaFormattedText(message *client.Message) *client.FormattedText {

	switch message.Content.MessageContentType() {
	case client.TypeMessagePhoto:
		return message.Content.(*client.MessagePhoto).Caption
	case client.TypeMessageAnimation:
		return message.Content.(*client.MessageAnimation).Caption
	case client.TypeMessageVideo:
		return message.Content.(*client.MessageVideo).Caption
	default:
		return nil
	}

}
