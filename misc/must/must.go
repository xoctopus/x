package must

func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

func BeTrue(ok bool) {
	if !ok {
		panic("must ok")
	}
}

func NoErrorV[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}
	return v
}

func BeTrueV[V any](v V, ok bool) V {
	if !ok {
		panic("must ok")
	}
	return v
}
