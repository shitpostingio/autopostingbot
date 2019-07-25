package telegramCaption

import (
	"io"

	"github.com/russross/blackfriday"
)

type telegramMarkdownRenderer struct {
	bfHTMLRenderer *blackfriday.HTMLRenderer
}

func (r *telegramMarkdownRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {

	switch node.Type {
	case blackfriday.Paragraph:
		if !entering {
			_, _ = w.Write([]byte("\n\n"))
		}
	case blackfriday.Item:
		if entering {
			_, _ = w.Write([]byte("-"))
		}
	case blackfriday.BlockQuote:
		if entering {
			_, _ = w.Write([]byte(">"))
		}
	case blackfriday.HorizontalRule:
		_, _ = w.Write([]byte("***\n\n"))
	case blackfriday.Heading:
		if entering {
			_, _ = w.Write([]byte("#"))
		} else {
			_, _ = w.Write([]byte("\n\n"))
		}
	case blackfriday.Code:
		fallthrough
	case blackfriday.Text:
		fallthrough
	case blackfriday.Emph:
		fallthrough
	case blackfriday.Strong:
		fallthrough
	case blackfriday.Link:
		r.bfHTMLRenderer.RenderNode(w, node, entering)

	}

	return blackfriday.GoToNext
}

func (r *telegramMarkdownRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	r.bfHTMLRenderer.RenderHeader(w, ast)
}

func (r *telegramMarkdownRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	r.bfHTMLRenderer.RenderFooter(w, ast)
}
