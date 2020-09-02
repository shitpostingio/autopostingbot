package api

import "github.com/zelenin/go-tdlib/client"

// GetTypeFromMessageType returns the message type of
// supported input media types.
func GetTypeFromMessageType(messageType string) string {

	switch messageType {
	case client.TypeMessagePhoto:
		return client.TypePhoto
	case client.TypeMessageAnimation:
		return client.TypeAnimation
	case client.TypeMessageVideo:
		return client.TypeVideo
	default:
		return ""
	}

}
