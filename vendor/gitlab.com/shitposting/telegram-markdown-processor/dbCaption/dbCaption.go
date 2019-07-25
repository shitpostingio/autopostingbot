package dbCaption

import (
	"fmt"
	"strings"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// PrepareCaptionForDB checks the entities to add markdown where is needed
func PrepareCaptionForDB(caption, channelName string, entities []tgbotapi.MessageEntity, surplus int) string {

	if entities == nil {
		return caption
	}

	builder := strings.Builder{}
	cRunesUTF16 := utf16.Encode([]rune(caption))
	previousIndex := 0
	normalizedIndex := 0

	for _, entity := range entities {

		normalizedIndex = entity.Offset - surplus

		switch entity.Type {
		case "italic":
			previousIndex = addMarkdownRune("_", previousIndex, normalizedIndex, entity.Length, cRunesUTF16, &builder)
		case "bold":
			previousIndex = addMarkdownRune("**", previousIndex, normalizedIndex, entity.Length, cRunesUTF16, &builder)
		case "code":
			previousIndex = addMarkdownRune("`", previousIndex, normalizedIndex, entity.Length, cRunesUTF16, &builder)
		case "pre":
			previousIndex = addMarkdownRune("```", previousIndex, normalizedIndex, entity.Length, cRunesUTF16, &builder)
		case "text_link":
			previousIndex = addMarkdownLink(entity.URL, previousIndex, normalizedIndex, entity.Length, cRunesUTF16, &builder)
		}

	}

	builder.WriteString(string(utf16.Decode(cRunesUTF16[previousIndex:])))
	output := strings.ReplaceAll(builder.String(), channelName, "")
	return strings.TrimSpace(output)
}

func addMarkdownRune(mdChars string, previousIndex, currentIndex, length int, cRunesUTF16 []uint16, builder *strings.Builder) int {

	captionBeginning := string(utf16.Decode(cRunesUTF16[previousIndex:currentIndex]))
	captionMiddle := string(utf16.Decode(cRunesUTF16[currentIndex:currentIndex+length]))
	trailing := ""

	/* WE WANT TO MOVE NEW LINES OUT OF MARKDOWN OR IT MAY CAUSE PROBLEMS */
	if strings.HasSuffix(captionMiddle, "\n") {
		captionMiddle = strings.TrimRight(captionMiddle, "\n")
		trailing = "\n"
	}

	builder.WriteString(fmt.Sprintf("%s%s%s%s%s", captionBeginning, mdChars, captionMiddle, mdChars, trailing))
	return currentIndex + length
}

func addMarkdownLink(url string, previousIndex, currentIndex, length int, cRunesUTF16 []uint16, builder *strings.Builder) int {

	captionBeginning := string(utf16.Decode(cRunesUTF16[previousIndex:currentIndex]))
	captionMiddle := string(utf16.Decode(cRunesUTF16[currentIndex:currentIndex+length]))
	trailing := ""

	/* WE WANT TO MOVE NEW LINES OUT OF MARKDOWN OR IT MAY CAUSE PROBLEMS */
	if strings.HasSuffix(captionMiddle, "\n") {
		captionMiddle = strings.TrimRight(captionMiddle, "\n")
		trailing = "\n"
	}

	builder.WriteString(fmt.Sprintf("%s[%s](%s)%s", captionBeginning, captionMiddle, url, trailing))
	return currentIndex + length
}
