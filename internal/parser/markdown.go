package parser

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type MarkdownContent struct {
	Content string
	HTML    string
}

func ParseMarkdown(content []byte) *MarkdownContent {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	html := markdown.ToHTML(content, p, nil)
	return &MarkdownContent{
		Content: string(content),
		HTML:    string(html),
	}
}
