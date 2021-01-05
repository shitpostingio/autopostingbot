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

	// TODO:	hasMarkup is a workaround:
	//			for some reason, emojis are represented by two UTF-16 characters
	// 			and writing them one by one corrupts them. If we don't have any
	//			meaningful entities, just return the full text
	if !hasMarkup(ft) {
		return text
	}

	//
	u16Text := utf16.Encode([]rune(ft.Text))
	b := strings.Builder{}

	//
	for i, v := range u16Text {

		// If it's before the start, ignore all entities
		if i < index {
			b.WriteString(string(utf16.Decode([]uint16{v})))
			continue
		}

		//Start entities
		for j := 0; j < len(ft.Entities); j++ {

			// Check start of entity
			if ft.Entities[j].Offset == int32(i) {

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
			if ft.Entities[j].Offset+ft.Entities[j].Length == int32(i) {

				tag, found := closingTags[ft.Entities[j].Type.TextEntityTypeType()]
				if found {
					b.WriteString(tag)
				}

			}

		}

		b.WriteString(string(utf16.Decode([]uint16{v})))

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
