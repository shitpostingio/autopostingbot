package api

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
)

var (

	uploadFunctions = map[string]func(int64, int64, string, string, []*client.TextEntity) (*client.Message, error){
		client.TypeAnimation: UploadAnimation,
		client.TypePhoto:     UploadPhoto,
		client.TypeVideo:     UploadVideo,
	}

	shareFunctions = map[string]func(int64, int64, string, string, []*client.TextEntity) (*client.Message, error){
		client.TypeAnimation: ShareAnimation,
		client.TypePhoto:     SharePhoto,
		client.TypeVideo:     ShareVideo,
	}

	fileInfoFunctions = map[string]func(*client.Message) *client.File{
		client.TypeMessageAnimation: GetAnimationFileInfoFromMessage,
		client.TypeMessagePhoto:     GetPhotoFileInfoFromMessage,
		client.TypeMessageVideo:     GetVideoFileInfoFromMessage,
	}
)

func SendMedia(mediaType string, chatID, replyToMessageID int64, remoteFileID, localFilePath, caption string, entities []*client.TextEntity) (*client.Message, error) {

	msg, err := ShareMedia(mediaType, chatID, replyToMessageID, remoteFileID, caption, entities)
	if err == nil {
		return msg, err
	}

	return UploadMedia(mediaType, chatID, replyToMessageID, localFilePath, caption, entities)

}

// ShareMedia shares a media file to a certain chat.
// If replyToMessageID is not 0, the media will be in reply to that message id.
// caption and entities can be used to attach a message with markdown.
func ShareMedia(mediaType string, chatID, replyToMessageID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

	send, found := shareFunctions[mediaType]
	if !found {
		return nil, fmt.Errorf("send function not found for media type %s", mediaType)
	}

	return send(chatID, replyToMessageID, remoteFileID, caption, entities)

}

// UploadMedia shares a media file to a certain chat.
// If replyToMessageID is not 0, the media will be in reply to that message id.
// caption and entities can be used to attach a message with markdown.
func UploadMedia(mediaType string, chatID, replyToMessageID int64, localFilePath, caption string, entities []*client.TextEntity) (*client.Message, error) {

	send, found := uploadFunctions[mediaType]
	if !found {
		return nil, fmt.Errorf("send function not found for media type %s", mediaType)
	}

	return send(chatID, replyToMessageID, localFilePath, caption, entities)

}

// GetMediaFileInfo returns the client.File structure for supported
// media types.
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

// GetMediaFormattedText returns the client.FormattedText structure for
// supported media types, nil otherwise.
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
