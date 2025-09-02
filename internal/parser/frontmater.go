package parser

import (
	"bytes"

	"gopkg.in/yaml.v2"
)

type Frontmatter struct {
	Title       string
	Date        string
	Tags        []string
	Description string
}

func ParseFrontmatter(content []byte) (*Frontmatter, []byte, error) {
	parts := bytes.Split(content, []byte("---"))
	if len(parts) < 3 {
		return nil, content, nil
	}

	var meta Frontmatter
	if err := yaml.Unmarshal(parts[1], &meta); err != nil {
		return nil, content, err
	}

	return &meta, bytes.Join(parts[2:], []byte("---")), nil
}
