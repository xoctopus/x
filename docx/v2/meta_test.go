package docx_test

import (
	"fmt"
	"go/ast"
	"maps"
	"slices"

	"golang.org/x/tools/go/packages"

	"github.com/xoctopus/x/docx/v2"
	"github.com/xoctopus/x/iterx"
	"github.com/xoctopus/x/misc/must"
)

func ExampleParseDocument() {
	m := docx.ParseDocument([]string{
		" Typename this line is title ",
		" this is line 1 for descriptions ",
		" this is line 2 for descriptions ",
		" +genx:directive directive parameters ... ",
		" + genx:... ", // invalid directive will be skipped
		" @attr k1=v1,attr parameters...",
		" TODO prefixed with keywords ",
		" todo keyword is case-sensitive",
		" this is line 3 for descriptions ",
		"", // empty line will be skipped
		" @def annotation content ",
		" @attr k2=v2,attr parameters...",
		" @def ",  // invalid annotation will be skipped
		" @ def ", // invalid annotation will be skipped
	})

	fmt.Println("Title")
	fmt.Println(m.Title("Typename"))
	fmt.Println()

	fmt.Println("Description Lines")
	for _, line := range m.Description().Lines() {
		fmt.Println(line)
	}
	fmt.Println()

	fmt.Println("Description String")
	fmt.Println(m.Description().String())
	fmt.Println()

	fmt.Println("Annotations")
	keys := slices.Collect(maps.Keys(m.Annotations()))
	slices.Sort(keys)
	for _, key := range keys {
		for _, anno := range m.Annotations()[key] {
			fmt.Println(anno.Name() + " " + anno.Text())
		}
	}
	fmt.Println()

	fmt.Println("Directives")
	for _, line := range m.Directives() {
		fmt.Println(line)
	}
	fmt.Println()

	fmt.Println("Empty documents")
	fmt.Println("BEGIN")
	fmt.Println(docx.ParseDocument(nil).Title(""))
	fmt.Println("END")

	// Output:
	// Title
	// this line is title
	//
	// Description Lines
	// this is line 1 for descriptions
	// this is line 2 for descriptions
	// todo keyword is case-sensitive
	// this is line 3 for descriptions
	//
	// Description String
	// this is line 1 for descriptions
	// this is line 2 for descriptions
	// todo keyword is case-sensitive
	// this is line 3 for descriptions
	//
	// Annotations
	// attr k1=v1,attr parameters...
	// attr k2=v2,attr parameters...
	// def annotation content
	//
	// Directives
	// +genx:directive directive parameters ...
	//
	// Empty documents
	// BEGIN
	//
	// END
}

func ExampleParseDocumentFromComments() {
	pkgpath := "github.com/xoctopus/x/docx/v2/testdata"
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.LoadAllSyntax,
	}, pkgpath)
	if err != nil {
		panic(err)
	}
	var p *packages.Package
	for _, p = range pkgs {
		if p.PkgPath == pkgpath {
			break
		}
	}

	must.BeTrue(p != nil)

	var (
		packageDocuments []*ast.CommentGroup
		typeDocuments    = make(map[string][]*ast.CommentGroup)
		valueDocuments   = make(map[string][]*ast.CommentGroup)
	)

	for _, f := range p.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.File:
				packageDocuments = append(packageDocuments, x.Doc)
			case *ast.GenDecl:
				for _, spec := range x.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						name := s.Name.String()
						typeDocuments[name] = append(typeDocuments[name], x.Doc, s.Doc)
					case *ast.ValueSpec:
						for _, ident := range s.Names {
							name := ident.String()
							valueDocuments[name] = append(valueDocuments[name], x.Doc, s.Doc)
						}
					}
				}
			case *ast.StructType:
				if x.Fields != nil {
					for _, field := range x.Fields.List {
						if field.Doc != nil {
							_ = field.Doc.List
						}
					}
				}
			default:
			}

			return true
		})
	}

	fmt.Println("Package documents")
	for _, line := range docx.ParseDocumentFromComments(packageDocuments...).Lines() {
		fmt.Println(line)
	}
	fmt.Println()

	fmt.Println("Type documents")
	for name, cs := range iterx.OrderedMapOf(typeDocuments) {
		fmt.Println(name)
		for _, line := range docx.ParseDocumentFromComments(cs...).Lines() {
			fmt.Println(line)
		}
	}
	fmt.Println()

	fmt.Println("Value documents")
	for name, cs := range iterx.OrderedMapOf(valueDocuments) {
		fmt.Println(name)
		for _, line := range docx.ParseDocumentFromComments(cs...).Lines() {
			fmt.Println(line)
		}
	}
	fmt.Println()

	// Output:
	// Package documents
	// Package testdata package level document
	// comments for testdata package
	//
	// Type documents
	// A
	// GenDecl for type list
	// A some type A
	// B
	// GenDecl for type list
	// B some type B
	// C
	// GenDecl for type list
	// C interface type
	// Fields
	// Fields for field document
	// this is a quoted multiline comments
	// Structure
	// Structure is a struct type for testing
	// line1
	// line2
	//
	// Value documents
	// ConstA
	// GenDecl for const list
	// ConstA int
	// ConstB
	// GenDecl for const list
	// ConstB float
	// VarA
	// GenDecl for var list
	// VarA int
	// VarB
	// GenDecl for var list
	// VarB float
}
