package api

import (
	"github.com/shitpostingio/autopostingbot/repository"
	"github.com/zelenin/go-tdlib/client"
)

// DownloadFile downloads synchronously a file with maximum priority.
func DownloadFile(fileID int32) (*client.File, error) {

	file, err := repository.Tdlib.DownloadFile(&client.DownloadFileRequest{
		FileId:      fileID,
		Priority:    32,
		Synchronous: true,
	})

	return file, err

}
