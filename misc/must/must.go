package must

func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

func OK(ok bool) {
	if !ok {
		panic("must ok")
	}
}
