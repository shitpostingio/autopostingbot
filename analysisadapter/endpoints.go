package analysisadapter

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
)

// getEndpoint returns the correct Analysis API endpoint
// for the supported media types.
func getEndpoint(mediaType, fileUniqueID string) string {

	switch mediaType {
	case client.TypePhoto:
		return fmt.Sprintf("%s/%s/%s", config.Address, config.ImageEndpoint, fileUniqueID)
	case client.TypeVideo, client.TypeAnimation:
		return fmt.Sprintf("%s/%s/%s", config.Address, config.VideoEndpoint, fileUniqueID)
	}

	return ""

}
