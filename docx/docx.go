package docx

import "github.com/xoctopus/x/contextx"

type Doc interface {
	DocOf(names ...string) ([]string, bool)
}

type Provider interface {
	DocOf(names ...string) ([]string, bool)
}

// nolint:unused,staticcheck
func Of(v any, prefix string, names ...string) ([]string, bool) {
	if x, ok := v.(Doc); ok {
		doc, ok := x.DocOf(names...)
		if ok {
			if prefix != "" && len(doc) > 0 {
				doc[0] = prefix + doc[0]
				return doc, true
			}
			return doc, true
		}
	}
	return []string{}, false
}

type tCtxProvider struct{}

var (
	ProviderFrom  = contextx.From[tCtxProvider, Provider]
	WithProvider  = contextx.With[tCtxProvider, Provider]
	MustProvider  = contextx.Must[tCtxProvider, Provider]
	CarryProvider = contextx.Carry[tCtxProvider, Provider]
)
