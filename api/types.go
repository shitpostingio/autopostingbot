package api

import "github.com/zelenin/go-tdlib/client"

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
