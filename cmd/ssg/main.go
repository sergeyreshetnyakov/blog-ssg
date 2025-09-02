package main

import (
	"flag"
	"html/template"
	"net/http"
	"os"

	"github.com/sergeyreshetnyakov/blog-ssg/internal/generator"
	"github.com/sergeyreshetnyakov/blog-ssg/internal/parser"
)

func main() {
	contentPath := flag.String("content", "content/first-post.md", "path to content")
	templatesDir := flag.String("templates", "templates", "path to templates")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(*contentPath)
		if err != nil {
			panic(err)
		}

		meta, content, err := parser.ParseFrontmatter(data)
		if err != nil {
			panic(err)
		}

		parsed := parser.ParseMarkdown(content)

		g := generator.New(*templatesDir, w)
		g.Generate(map[string]interface{}{
			"Title":   meta.Title,
			"Content": template.HTML(parsed.HTML),
		})
	})

	http.ListenAndServe(":8080", nil)
}
