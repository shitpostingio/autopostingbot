package legacy

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/types"
)

func NewMediaTypeFromOld(typeID uint) (string, bool) {

	switch typeID {
	case types.Image:
		return client.TypePhoto, true
	case types.Video:
		return client.TypeVideo, true
	case types.Animation:
		return client.TypeAnimation, true
	}

	return "", false

}
