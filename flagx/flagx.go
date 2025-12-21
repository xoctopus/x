package flagx

type _Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Flagger[U _Uint] interface {
	Value() U
	Is(U) bool
	With(U) U
	Trim(U) U
}

func NewFlag[U _Uint]() Flagger[U] {
	return &Flag[U]{u: U(0)}
}

type Flag[U _Uint] struct {
	u U
}

func (f Flag[U]) Value() U {
	return f.u
}

func (f Flag[U]) Is(u U) bool {
	return f.u&u == u
}

func (f *Flag[U]) With(u U) U {
	f.u |= u
	return f.u
}

func (f *Flag[U]) Trim(u U) U {
	f.u &= ^u
	return f.u
}
