package telegramCaption

import (
	"fmt"

	"github.com/russross/blackfriday"
)

var (
	renderer telegramMarkdownRenderer
)

func init() {
	renderer = telegramMarkdownRenderer{
		bfHTMLRenderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.HTMLFlagsNone,
		})}
}

// PrepareCaptionToSend prepares the caption for posting
// by running it through the markdown converter
func PrepareCaptionToSend(inputCaption, channelName string) string {

	if channelName == "" {
		return prepareCaptionNoChannel(inputCaption)
	}

	return prepareCaptionWithChannel(inputCaption, channelName)
}

func prepareCaptionNoChannel(inputCaption string) string {

	if inputCaption == "" {
		return ""
	}

	return string(blackfriday.Run([]byte(inputCaption), blackfriday.WithExtensions(blackfriday.Strikethrough), blackfriday.WithRenderer(&renderer)))
}

func prepareCaptionWithChannel(inputCaption, channelName string) string {

	if inputCaption == "" {
		return "@" + channelName
	}

	bfInput := fmt.Sprintf("%s\n\n@%s", inputCaption, channelName)
	return string(blackfriday.Run([]byte(bfInput), blackfriday.WithExtensions(blackfriday.Strikethrough), blackfriday.WithRenderer(&renderer)))
}
