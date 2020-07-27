package api

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/repository"
)

func GetFormattedText(text string) (*client.FormattedText, error) {

	formattedText, err := repository.Tdlib.ParseTextEntities(&client.ParseTextEntitiesRequest{
		Text:      text,
		ParseMode: &client.TextParseModeHTML{},
	})

	return formattedText, err

}
