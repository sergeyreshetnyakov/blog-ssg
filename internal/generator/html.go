package generator

import (
	"html/template"
	"io"
	"path/filepath"
)

type Generator struct {
	TemplatesPath string
	Output        io.Writer
}

type Post struct {
	URL         string
	Title       string
	Description string
	Date        string
}

type IndexPage struct {
	Title string
	Posts []Post
}

func New(templatesPath string, output io.Writer) *Generator {
	return &Generator{
		TemplatesPath: templatesPath,
		Output:        output,
	}
}

func (g *Generator) Generate(data any) error {
	tmpl, err := template.ParseFiles(filepath.Join(g.TemplatesPath, "default.html"))
	if err != nil {
		return err
	}

	return tmpl.Execute(g.Output, data)
}

func (g *Generator) GenerateIndex(data any) error {
	tmpl, err := template.ParseFiles(filepath.Join(g.TemplatesPath, "index.html"))
	if err != nil {
		return err
	}

	return tmpl.Execute(g.Output, data)
}
