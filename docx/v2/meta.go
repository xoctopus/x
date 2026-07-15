package docx

import (
	"go/ast"
	"strings"
)

var keywords = map[string]struct{}{
	"XXX":      {},
	"ERROR":    {},
	"BUG":      {},
	"HACK":     {},
	"TODO":     {},
	"PERF":     {},
	"OPTIMIZE": {},
	"REFACTOR": {},
	"NOTE":     {},
	"INFO":     {},
	"NOTICE":   {},
	"WARNING":  {},
	"REVIEW":   {},
}

func ParseDocument(lines []string) *Meta {
	d := &Meta{
		annotations: Annotations{},
		directives:  []string{},
	}
	if len(lines) == 0 {
		return d
	}

	i := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if parts := strings.Split(line, " "); len(parts) > 0 {
			if _, ok := keywords[parts[0]]; ok {
				continue
			}
		}
		if after, ok := strings.CutPrefix(line, "@"); ok {
			idx := strings.Index(after, " ")
			if idx == 0 || idx == -1 {
				continue
			}
			d.annotations.add(after[0:idx], strings.TrimSpace(after[idx+1:]))
			continue
		}
		if i == 0 {
			d.title = line
			i++
			continue
		}
		if after, ok := strings.CutPrefix(line, "+"); ok {
			if strings.HasPrefix(after, " ") {
				continue
			}
			d.directives = append(d.directives, line)
			continue
		}
		d.description.append(line)
	}
	if len(d.description.lines) > 0 {
		d.description.text = strings.Join(d.description.lines, "\n")
	}
	return d
}

func ParseDocumentFromComments(groups ...*ast.CommentGroup) *Meta {
	lines := make([]string, 0)
	for _, g := range groups {
		if g != nil {
			for _, c := range g.List {
				if strings.HasPrefix(c.Text, "/*") {
					line := c.Text
					line = strings.TrimPrefix(line, "/*")
					line = strings.TrimSuffix(line, "*/")
					lines = append(lines, strings.Split(line, "\n")...)
				} else {
					line := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
					lines = append(lines, line)
				}
			}
		}
	}
	return ParseDocument(lines)
}

type Meta struct {
	title       string
	description Description
	annotations Annotations
	directives  []string
}

func (m *Meta) Lines() []string {
	return append([]string{m.title}, m.description.Lines()...)
}

func (m *Meta) Title(prefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(m.title, prefix+" "))
}

func (m *Meta) Annotations() Annotations {
	return m.annotations
}

func (m *Meta) Description() *Description {
	return &m.description
}

func (m *Meta) Directives() []string {
	return m.directives
}
