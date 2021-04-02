package caption

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
	"strings"
	"unicode/utf16"
)

var (
	openingTags = map[string]string{
		client.TypeTextEntityTypeBold:          "<b>",
		client.TypeTextEntityTypeItalic:        "<i>",
		client.TypeTextEntityTypeUnderline:     "<u>",
		client.TypeTextEntityTypeStrikethrough: "<s>",
		client.TypeTextEntityTypeCode:          "<code>",
		client.TypeTextEntityTypePre:           "<pre>",
		client.TypeTextEntityTypeTextUrl:       `<a href="%s">`,
	}

	closingTags = map[string]string{
		client.TypeTextEntityTypeBold:          "</b>",
		client.TypeTextEntityTypeItalic:        "</i>",
		client.TypeTextEntityTypeUnderline:     "</u>",
		client.TypeTextEntityTypeStrikethrough: "</s>",
		client.TypeTextEntityTypeCode:          "</code>",
		client.TypeTextEntityTypePre:           "</pre>",
		client.TypeTextEntityTypeTextUrl:       "</a>",
	}
)

// ToHTMLCaptionWithCustomStart converts a client.FormattedText to a plaintext
// string with HTML markup.
// By setting index, one can specify where the parsing should start.
func ToHTMLCaptionWithCustomStart(ft *client.FormattedText, index int) string {

	//
	text := strings.ReplaceAll(ft.Text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")

	//
	u16Text := utf16.Encode([]rune(ft.Text))
	b := strings.Builder{}

	// Change the offsets and the length from UTF-16 to UTF-8
	utf8Index := len(utf16.Decode(u16Text[:index]))
	utf8Entities := make([]client.TextEntity, len(ft.Entities))

	for i, e := range ft.Entities {

		//
		entityStart := len(string(utf16.Decode(u16Text[:e.Offset])))
		entityLength := len(string(utf16.Decode(u16Text[:e.Offset + e.Length]))) - entityStart

		// Only save the numerical values, we will use actual
		// entities for the rest
		entity := client.TextEntity{
			Offset: int32(entityStart),
			Length: int32(entityLength),
		}

		utf8Entities[i] = entity

	}

	// Iterate over the UTF-8 text runes
	for i, r := range ft.Text {

		// If it's before the start, ignore all entities
		if i < utf8Index {
			b.WriteRune(r)
			continue
		}

		//Start entities
		for j := 0; j < len(ft.Entities); j++ {

			// Check start of entity
			if utf8Entities[j].Offset == int32(i) {

				// We need to add this entity
				// Special check for links
				if ft.Entities[j].Type.TextEntityTypeType() == client.TypeTextEntityTypeTextUrl {

					e := ft.Entities[j].Type.(*client.TextEntityTypeTextUrl)
					b.WriteString(fmt.Sprintf(openingTags[client.TypeTextEntityTypeTextUrl], e.Url))
					continue

				}

				tag, found := openingTags[ft.Entities[j].Type.TextEntityTypeType()]
				if found {
					b.WriteString(tag)
				}

			}

		}

		// End entities: in reverse order so we follow the order
		// in which they started
		for j := len(ft.Entities) - 1; j >= 0; j-- {

			// Check end of entity
			if utf8Entities[j].Offset + utf8Entities[j].Length == int32(i) {

				tag, found := closingTags[ft.Entities[j].Type.TextEntityTypeType()]
				if found {
					b.WriteString(tag)
				}

			}

		}

		b.WriteRune(r)

	}

	// check for entities that end after the text
	// as the for loop above won't get there
	for j := len(ft.Entities) - 1; j >= 0; j-- {

		// Check end of entity
		if ft.Entities[j].Offset+ft.Entities[j].Length == int32(len(u16Text)) {

			tag, found := closingTags[ft.Entities[j].Type.TextEntityTypeType()]
			if found {
				b.WriteString(tag)
			}

		}

	}

	return b.String()

}

// ToHTMLCaption converts a client.FormattedText to a plaintext
// string with HTML markup.
func ToHTMLCaption(ft *client.FormattedText) string {
	return ToHTMLCaptionWithCustomStart(ft, 0)
}

// hasMarkup returns true if the message contains
// entities from the supported list
func hasMarkup(ft *client.FormattedText) bool {

	if ft.Entities == nil {
		return false
	}

	for _, e := range ft.Entities {
		_, found := openingTags[e.Type.TextEntityTypeType()]
		if found {
			return true
		}
	}

	return false

}
