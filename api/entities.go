package api

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/repository"
)

func GetFormattedText(caption string) (*client.FormattedText, error) {

	formattedText, err := repository.Tdlib.ParseTextEntities(&client.ParseTextEntitiesRequest{
		Text:      caption,
		ParseMode: &client.TextParseModeHTML{},
	})

	fmt.Println("ERRORE? ", err)
	if formattedText.Entities != nil && len(formattedText.Entities) > 0 {
		fmt.Println(formattedText.Entities[0].Type.TextEntityTypeType())
	}

	return formattedText, err

}
