package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sergeyreshetnyakov/blog-ssg/internal/generator"
	"github.com/sergeyreshetnyakov/blog-ssg/internal/parser"
)

func processFiles(contentPath, templatesPath, outputPath string) ([]generator.Post, error) {
	var posts []generator.Post
	return posts, filepath.Walk(filepath.Join(contentPath, "posts"), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return err
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		meta, content, err := parser.ParseFrontmatter(data)
		if err != nil {
			return err
		}

		parsed := parser.ParseMarkdown(content)
		posts = append(posts, generator.Post{
			URL:         "pages/" + strings.Replace(info.Name(), ".md", ".html", 1),
			Title:       meta.Title,
			Description: meta.Description,
			Date:        meta.Date,
		})
		out, err := os.Create(filepath.Join(outputPath, "pages", strings.Replace(info.Name(), ".md", ".html", 1)))
		if err != nil {
			return err
		}
		defer out.Close()

		g := generator.New(templatesPath, out)
		return g.Generate(map[string]interface{}{
			"Title":   meta.Title,
			"Content": template.HTML(parsed.HTML),
			"Date":    meta.Date,
		})
	})
}

func processIndex(contentPath, templatesPath, outputPath string, posts []generator.Post) error {
	data, err := os.ReadFile(filepath.Join(contentPath, "index.md"))
	if err != nil {
		return err
	}

	meta, content, err := parser.ParseFrontmatter(data)
	if err != nil {
		return err
	}

	parsed := parser.ParseMarkdown(content)

	out, err := os.Create(filepath.Join(outputPath, "index.html"))
	if err != nil {
		return err
	}
	defer out.Close()

	g := generator.New(templatesPath, out)
	return g.GenerateIndex(map[string]interface{}{
		"Title":   meta.Title,
		"Content": template.HTML(parsed.HTML),
		"Posts":   posts,
	})
}

func parseCSS(templatesPath, outputPath string) error {
	f, err := os.ReadFile(filepath.Join(templatesPath, "style.css"))
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outputPath, "style.css"), f, 0755)
}

func main() {
	contentPath := flag.String("content", "content", "path to content")
	templatesPath := flag.String("templates", "templates", "path to templates")
	outputPath := flag.String("output", "output", "path to folder with generated files")

	posts, err := processFiles(*contentPath, *templatesPath, *outputPath)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = processIndex(*contentPath, *templatesPath, *outputPath, posts)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	parseCSS(*templatesPath, *outputPath)
	fs := http.FileServer(http.Dir(*outputPath))
	http.Handle("/", fs)
	fmt.Println("The server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
