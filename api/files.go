package api

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/repository"
)

func DownloadFile(fileID int32) (*client.File, error) {

	file, err := repository.Tdlib.DownloadFile(&client.DownloadFileRequest{
		FileId:      fileID,
		Priority:    32,
		Synchronous: true,
	})

	return file, err

}
