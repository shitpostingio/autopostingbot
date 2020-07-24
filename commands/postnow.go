package commands

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
)

type PostNowCommandHandler struct {

}

func (PostNowCommandHandler) Handle(arguments string, message *client.Message) error {

	fileInfo, err := api.GetMediaFileInfo(message)
	if err != nil {
		return err
	}

	post, err := dbwrapper.FindPostByUniqueID(fileInfo.Remote.UniqueId)
	if err != nil {
		return err
	}

	//POST
	fmt.Println(post)
	return err

}
