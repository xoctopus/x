package docx

type Doc interface {
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
