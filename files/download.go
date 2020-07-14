package files

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/repository"
)

func DownloadFile(fileID int32) (*client.File, error) {
	return repository.Tdlib.DownloadFile(&client.DownloadFileRequest{
		FileId:      fileID,
		Priority:    32,
		Synchronous: true,
	})
}
