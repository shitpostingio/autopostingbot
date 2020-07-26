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

func ToHTMLCaption(ft *client.FormattedText) string {

	text := strings.ReplaceAll(ft.Text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")

	if ft.Entities == nil {
		return text
	}

	u16Text := utf16.Encode([]rune(ft.Text))
	b := strings.Builder{}

	for i, v := range u16Text {

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

func ToHTMLCaptionWithCustomStart(ft *client.FormattedText, index int) string {

	text := strings.ReplaceAll(ft.Text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")

	if ft.Entities == nil {
		return text
	}

	u16Text := utf16.Encode([]rune(ft.Text))
	b := strings.Builder{}

	for i, v := range u16Text {

		if i < index {
			continue
		}

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
