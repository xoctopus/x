package ptrx

// Ptr returns pointer ov V
// Deprecated: upgrade to ^go1.26 use new(V) instead
//
//go:fix inline
func Ptr[V any](v V) *V { return new(v) }
