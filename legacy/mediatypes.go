package legacy

import (
	"github.com/zelenin/go-tdlib/client"
)

func NewMediaTypeFromOld(typeID uint) (string, bool) {

	//TODO: CONTROLLARE I VALORI
	switch typeID {
	case 0:
		return client.TypePhoto, true
	case 1:
		return client.TypeVideo, true
	case 2:
		return client.TypeAnimation, true
	}

	return "", false

}
