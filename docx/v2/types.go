package docx

import (
	"github.com/xoctopus/x/contextx"
)

type Provider interface {
	Title(names ...string) (string, bool)
	Description(names ...string) *Description
	Annotations(names ...string) Annotations
	Directives() Directives
}

type Description struct {
	lines []string
	text  string
}

func (d *Description) Lines() []string {
	return d.lines
}

func (d *Description) append(line string) {
	d.lines = append(d.lines, line)
}

func (d *Description) String() (text string) {
	return d.text
}

type Annotation struct {
	name string
	text string
}

func (a *Annotation) Name() string {
	return a.name
}

func (a *Annotation) Text() string {
	return a.text
}

type Annotations map[string][]Annotation

func (as Annotations) add(name, text string) {
	as[name] = append(as[name], Annotation{name: name, text: text})
}

type Directive struct {
	name   string
	suffix string
}

type Directives map[string]Directive

type tCtxProvider struct{}

var (
	ProviderFrom  = contextx.From[tCtxProvider, Provider]
	WithProvider  = contextx.With[tCtxProvider, Provider]
	MustProvider  = contextx.Must[tCtxProvider, Provider]
	CarryProvider = contextx.Carry[tCtxProvider, Provider]
)
