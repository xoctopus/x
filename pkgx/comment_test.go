package pkgx_test

import (
	"go/types"
	"testing"

	. "github.com/onsi/gomega"
)

// func TestCommentScanner(t *testing.T) {
// 	fset := token.NewFileSet()
//
// 	fpth, err := filepath.Abs(path.Join(root, "comments.go"))
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
//
// 	fast, err := parser.ParseFile(fset, fpth, nil, parser.ParseComments)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
//
// 	ast.Inspect(fast, func(node ast.Node) bool {
// 		comments := strings.Split(NewCommentScanner(fset, fast).CommentsOf(node), "\n")
// 		NewWithT(t).Expect(3 >= len(comments)).To(BeTrue())
// 		return true
// 	})
// }

func TestPkgComments(t *testing.T) {
	for _, v := range []struct {
		name   string
		object types.Object // identifier object
		expect []string     // expect identifier's comment
	}{
		{"Date", pkg.TypeName("Date"), []string{"Date defines corresponding time.Time"}},
		{"Test", pkg.Var("test"), []string{"var"}},
		{"A", pkg.Const("A"), []string{"A comment", "A inline comment"}},
		{"Print", pkg.Func("Print"), []string{"Print function"}},
	} {
		t.Run(v.name, func(t *testing.T) {
			NewWithT(t).Expect(pkg.CommentsOf(pkg.IdentOf(v.object))).To(Equal(v.expect))
		})
	}
}
