package resultx

func ResultsOf(rs ...any) Results {
	return Results(rs)
}

type Results []any

func (rs Results) At(i int) any {
	return rs[i]
}

// func (rs Results) At[T any](i int) T {
// 	return rs[i].(T)
// }

func At[T any](rs Results, i int) T {
	return rs[i].(T)
}
