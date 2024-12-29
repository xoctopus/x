package pkgx_test

import (
	"strconv"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/misc/must"
	. "github.com/xoctopus/x/pkgx"
)

func TestImportPathAndExpose(t *testing.T) {
	cases := []struct {
		imported string
		expose   string
		s        string
	}{
		{"", "B", "B"},
		{"testing", "B", "testing.B"},
		{"a.b.c.d/c", "B", "a.b.c.d/c.B"},
		{"e", "B", "a.b.c.d/vendor/e.B"},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			imported, expose := ImportPathAndTypeID(c.s)
			NewWithT(t).Expect(imported).To(Equal(c.imported))
			NewWithT(t).Expect(expose).To(Equal(c.expose))
		})
	}
}

var pkg *Pkg

func init() {
	pkg, _ = LoadFrom("./testdata")
	must.NotNilWrap(pkg, "load package info failed")
}

func TestPkg_Imports(t *testing.T) {
	pkgid := "github.com/xoctopus/x/pkgx/testdata/sub"
	imported := false
	for _, p := range pkg.Imports() {
		if p.ID == pkgid {
			imported = true
			break
		}
	}
	NewWithT(t).Expect(imported).To(BeTrue())
}

func TestPkg_Const(t *testing.T) {
	cases := []struct {
		name     string
		varName  string
		isNil    bool
		comments []string
	}{
		{"ConstA", "A", false, []string{"A comment", "A inline comment"}},
		{"ConstB", "B", false, []string{"B comment", "B inline comment"}},
		{"ConstC", "C", false, []string{"C comment", "C inline comment"}},
		{"Undefined", "Undefined", true, nil},
		{"Uname", "_", false, []string{"placeholder"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := pkg.Const(c.varName)
			if !c.isNil {
				NewWithT(t).Expect(v).NotTo(BeNil())
			} else {
				NewWithT(t).Expect(v).To(BeNil())
			}
			comments := pkg.CommentsOf(pkg.IdentOf(v))
			NewWithT(t).Expect(comments).To(Equal(c.comments))
		})
	}
}

func TestPkg_TypeName(t *testing.T) {
	cases := []struct {
		name     string
		comments []string
	}{
		{"Date", []string{"Date defines corresponding time.Time"}},
		{"Test", []string{"Test struct"}},
		{"Test2", []string{"Test2 struct"}},
		{"Func", nil},
		{"String", []string{}},
		{"Bar", []string{"Bar interface"}},
		{"Foo", []string{"Foo struct"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			comments := pkg.CommentsOf(pkg.IdentOf(pkg.TypeName(c.name)))
			NewWithT(t).Expect(comments).To(Equal(c.comments))
		})
	}
}

func TestPkg_Var(t *testing.T) {
	cases := []struct {
		name     string
		comments []string
	}{
		{"test", []string{"var"}},
		{"test2", []string{"test2"}},
		{"test3", []string{"test3"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			comments := pkg.CommentsOf(pkg.IdentOf(pkg.Var(c.name)))
			NewWithT(t).Expect(comments).To(Equal(c.comments))
		})
	}
}

func TestPkg_Func(t *testing.T) {
	cases := []struct {
		name     string
		comments []string
	}{
		{"v", nil},
		{"call", nil},
		{"CurryCall", nil},
		{"main", []string{"main"}},
		{"Print", []string{"Print function"}},
		{"fn", []string{"func fn"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			comments := pkg.CommentsOf(pkg.IdentOf(pkg.Func(c.name)))
			NewWithT(t).Expect(comments).To(Equal(c.comments))
		})
	}
}

func TestPkg_PkgByPath(t *testing.T) {
	cases := []struct {
		name    string
		exists  bool
		path    string
		pkgname string
	}{
		{"Sub", true, "github.com/xoctopus/x/pkgx/testdata/sub", "sub"},
		{"UnImported", false, "", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := pkg.PkgByPath(c.path)
			if c.exists {
				NewWithT(t).Expect(v.ID).To(Equal(c.path))
				NewWithT(t).Expect(v.PkgPath).To(Equal(c.path))
				NewWithT(t).Expect(v.Name).To(Equal(c.pkgname))
				return
			}
			NewWithT(t).Expect(v).To(BeNil())
		})
	}
}

func TestPkg_PkgByPos_PkgOf_PkgInfoOf_FileOf(t *testing.T) {
	cases := []struct {
		name  string
		pos   Pos
		pkgid string
	}{
		{"TypeName", pkg.TypeName("Date"), "github.com/xoctopus/x/pkgx/testdata"},
		{"Func", pkg.Func("Print"), "github.com/xoctopus/x/pkgx/testdata"},
		{"Var", pkg.Var("test"), "github.com/xoctopus/x/pkgx/testdata"},
		{"Const", pkg.Const("A"), "github.com/xoctopus/x/pkgx/testdata"},
		{"Invalid", nil, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pp := pkg.PkgByPos(c.pos)
			pt := pkg.PkgOf(c.pos)
			ti := pkg.PkgInfoOf(c.pos)
			af := pkg.FileOf(c.pos)
			if c.pkgid != "" {
				NewWithT(t).Expect(pp.ID).To(Equal(c.pkgid))
				NewWithT(t).Expect(pt.Path()).To(Equal(c.pkgid))
				NewWithT(t).Expect(ti).NotTo(BeNil())
				NewWithT(t).Expect(af.Name.Name).To(Equal("main"))
				NewWithT(t).Expect(af.Name.String()).To(Equal("main"))
			} else {
				NewWithT(t).Expect(pp).To(BeNil())
				NewWithT(t).Expect(pt).To(BeNil())
				NewWithT(t).Expect(ti).To(BeNil())
				NewWithT(t).Expect(af).To(BeNil())
			}
		})
	}
}

// func ExamplePrintDefs() {
// 	defines := make([]types.Object, 0)
// 	for _, def := range pkg.TypesInfo.Defs {
// 		if def == nil {
// 			continue
// 		}
// 		defines = append(defines, def)
// 	}
// 	sort.Slice(defines, func(i, j int) bool {
// 		return defines[i].String() < defines[j].String()
// 	})
// 	fmt.Println("Types:")
// 	for _, def := range defines {
// 		if d, ok := def.(*types.TypeName); ok {
// 			fmt.Printf("%10s: %s\n", d.Name(), d.Type().String())
// 		}
// 	}
// 	fmt.Println("Constants:")
// 	for _, def := range defines {
// 		if d, ok := def.(*types.Const); ok {
// 			fmt.Printf("%10s: %s\n", d.Name(), d.Type().String())
// 		}
// 	}
// 	fmt.Println("Vars:")
// 	for _, def := range defines {
// 		if d, ok := def.(*types.Var); ok {
// 			fmt.Printf("%10s: %s\n", d.Name(), d.Type().String())
// 		}
// 	}
//
// 	// Output:
// }
//
// func ExampleTypesExpr() {
// 	exprs := map[string]struct{}{}
// 	for expr, tv := range pkg.TypesInfo.Types {
// 		exprs[reflect.TypeOf(expr).String()] = struct{}{}
// 		_ = tv
// 	}
//
// 	typenames := maps.Keys(exprs)
// 	sort.Slice(typenames, func(i, j int) bool {
// 		return typenames[i] < typenames[j]
// 	})
//
// 	for _, name := range typenames {
// 		fmt.Println(name)
// 	}
// 	// Output:
// }
//
// func TestPkgFuncReturns(t *testing.T) {
// 	var pkgid string
//
// 	{
// 		_, current, _, _ := runtime.Caller(0)
// 		dir := filepath.Join(filepath.Dir(current), "./__tests__")
// 		pkgid = must.NoErrorV(PkgIdByPath(dir))
// 	}
//
// 	var cases = []struct {
// 		FuncName string
// 		Results  [][]string
// 	}{
// 		{
// 			"FuncSingleReturn",
// 			[][]string{{"untyped int(2)"}},
// 		},
// 		{
// 			"FuncSelectExprReturn",
// 			[][]string{{"string"}},
// 		},
// 		{
// 			"FuncWillCall",
// 			[][]string{
// 				{"interface{}"},
// 				{strings.Join([]string{pkgid, "String"}, ".")},
// 			},
// 		},
// 		{
// 			"FuncReturnWithCallDirectly",
// 			[][]string{
// 				{"interface{}"},
// 				{strings.Join([]string{pkgid, "String"}, ".")},
// 			},
// 		},
// 		{
// 			"FuncWithNamedReturn",
// 			[][]string{
// 				{"interface{}"},
// 				{strings.Join([]string{pkgid, "String"}, ".")},
// 			},
// 		},
// 		{
// 			"FuncSingleNamedReturnByAssign",
// 			[][]string{
// 				{`untyped string("1")`},
// 				{strings.Join([]string{pkgid, `String("2")`}, ".")},
// 			},
// 		},
// 		{
// 			"FuncWithSwitch",
// 			[][]string{
// 				{
// 					`untyped string("a1")`,
// 					`untyped string("a2")`,
// 					`untyped string("a3")`,
// 				},
// 				{
// 					strings.Join([]string{pkgid, `String("b1")`}, "."),
// 					strings.Join([]string{pkgid, `String("b2")`}, "."),
// 					strings.Join([]string{pkgid, `String("b3")`}, "."),
// 				},
// 			},
// 		},
// 	}
// 	for _, c := range cases {
// 		t.Run(c.FuncName, func(t *testing.T) {
// 			results, n := pkg.FuncResults(pkg.Func(c.FuncName))
// 			NewWithT(t).Expect(results).To(HaveLen(n))
// 			NewWithT(t).Expect(c.Results).To(Equal(PrintValues(pkg.Fset, results)))
// 		})
// 	}
// }
//
// func PrintValues(fs *token.FileSet, res map[int][]TypeAndValueExpr) [][]string {
// 	if res == nil {
// 		return [][]string{}
// 	}
// 	ret := make([][]string, len(res))
// 	for i := range ret {
// 		tve := res[i]
// 		ret[i] = make([]string, len(tve))
// 		for j, v := range tve {
// 			fmt.Println(v.Type, v.Value)
// 			if v.Value == nil {
// 				ret[i][j] = v.Type.String()
// 			} else {
// 				ret[i][j] = fmt.Sprintf("%s(%s)", v.Type, v.Value)
// 			}
// 		}
// 	}
// 	return ret
// }

func PrintAstInfo(t *testing.T) {
	// fset := token.NewFileSet()
	// fpth, err := filepath.Abs(path.Join(root, "ast.go"))
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// fast, err := parser.ParseFile(fset, fpth, nil, parser.AllErrors)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// _ = ast.Print(fset, fast)
}
