package generator

import (
	"html/template"
	"io"
	"path/filepath"
)

type Generator struct {
	TemplateDir string
	Output      io.Writer
}

func New(templateDir string, output io.Writer) *Generator {
	return &Generator{
		TemplateDir: templateDir,
		Output:      output,
	}
}

func (g *Generator) Generate(data interface{}) error {
	tmpl, err := template.ParseFiles(filepath.Join(g.TemplateDir, "default.html"))
	if err != nil {
		return err
	}

	return tmpl.Execute(g.Output, data)
}
