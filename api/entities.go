package api

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/repository"
)

// GetFormattedText returns the client.FormattedText structure of the
// input text, parsing it as HTML markup.
func GetFormattedText(text string) (*client.FormattedText, error) {

	formattedText, err := repository.Tdlib.ParseTextEntities(&client.ParseTextEntitiesRequest{
		Text:      text,
		ParseMode: &client.TextParseModeHTML{},
	})

	return formattedText, err

}
