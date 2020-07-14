package updates

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/files"
)

func handlePhoto(message *client.Message) {

	//
	fileInfo, err := api.GetMediaFileInfo(message)
	if err != nil {
		log.Error("handlePhoto: ", err)
		return
	}

	//
	fileInfo, err = files.DownloadFile(fileInfo.Id)
	if err != nil {
		log.Error("handlePhoto: ", err)
		return
	}





	//utf16Text := utf16.Encode([]rune(messageContent.Caption.Text))



}

