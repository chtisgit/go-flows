package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
)

type function struct {
	Name   string
	Return string
	Args   string
	Code   string
}

/*
func poorMansImporter(imports map[string]*ast.Object, path string) (*ast.Object, error) {
	pkg := imports[path]
	if pkg == nil {
		// note that strings.LastIndex returns -1 if there is no "/"
		pkg = ast.NewObj(ast.Pkg, path[strings.LastIndex(path, "/")+1:])
		pkg.Data = ast.NewScope(nil) // required by ast.NewPackage for dot-import
		imports[path] = pkg
	}
	return pkg, nil
}

type funcall struct {
	result    string
	arguments []string
}

var libs = make(map[string]map[string]funcall)

func parseLib(lib string, fun string) funcall {
	if _, ok := libs[lib]; ok {
		return libs[lib][fun]
	}
	libs[lib] = make(map[string]funcall)
	pkg, err := build.Default.Import(lib, ".", 0)
	if err != nil {
		log.Fatal(err)
	}
	fs := token.NewFileSet()
	pkgfiles := append(pkg.GoFiles, pkg.CgoFiles...)
	asts := make(map[string]*ast.File)
	for _, name := range pkgfiles {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		parsedFile, err := parser.ParseFile(fs, path.Join(pkg.Dir, name), nil, 0)
		if err != nil {
			log.Fatal(err)
		}
		asts[name] = parsedFile
	}
	p, err := ast.NewPackage(fs, asts, poorMansImporter, nil)
	for _, object := range p.Scope.Objects {
		if object.Kind != ast.Fun {
			continue
		}
		decl := object.Decl.(*ast.FuncDecl)
		params := make([]string, len(decl.Type.Params.List))
		for i, p := range decl.Type.Params.List {
			params[i] = fmt.Sprint(p.Type)
		}
		libs[lib][object.Name] = funcall{
			fmt.Sprint(decl.Type.Results.List[0].Type),
			params,
		}
	}
	return libs[lib][fun]
} */

func main() {
	f, err := os.Create("number_interfaces.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fs := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fs, "number.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	cmap := ast.NewCommentMap(fs, parsedFile, parsedFile.Comments)

	var types []string
	imports := make(map[string]struct{})
	var functions []function

	fileTemplate := template.New("main")
	fileTemplate.Funcs(map[string]interface{}{
		"CallTemplate": func(name string, data interface{}) (ret string, err error) {
			buf := bytes.NewBuffer([]byte{})
			err = fileTemplate.ExecuteTemplate(buf, name, data)
			ret = buf.String()
			return
		},
	})
	fileTemplate.Parse(`// GENERATED BY number_generate.go; DO NOT EDIT!

package flows

{{range $key, $value := .Imports }}import "{{ $key }}"
{{ end }}
{{ $functions := .Functions }}{{range .Types }}{{ $t := . }}{{range $functions }}
func (a {{ $t }}) {{ .Name }}({{ .Args }}) {{ .Return }} {
	return {{CallTemplate .Name $t}};
}{{end}}{{end}}

`)

	for node, comment := range cmap {
		switch node.(type) {
		case *ast.GenDecl:
			node := node.(*ast.GenDecl)
			if node.Tok == token.TYPE {
				types = append(types, node.Specs[0].(*ast.TypeSpec).Name.String())
			}
		case *ast.Field:
			node := node.(*ast.Field)
			comment := strings.SplitN(comment[0].Text(), ":", 2)
			var f function
			f.Name = node.Names[0].String()
			funcDef := node.Type.(*ast.FuncType)
			f.Return = fmt.Sprint(funcDef.Results.List[0].Type)

			if comment[0] != "oper" {
				continue
			}

			expr, err := parser.ParseExpr(comment[1])
			if err != nil {
				log.Fatal(err)
			}
			var highest byte
			ast.Inspect(expr, func(n ast.Node) bool {
				switch n.(type) {
				case *ast.Ident:
					name := &n.(*ast.Ident).Name
					if len(*name) == 1 && *name != "a" {
						letter := *name
						if letter[0]-'a' > highest {
							highest = letter[0] - 'a'
						}
						*name += ".({{ . }})"
					}
				case *ast.CallExpr:
					if n, ok := n.(*ast.CallExpr); ok {
						if fun, ok := n.Fun.(*ast.SelectorExpr); ok {
							imports[fun.X.(*ast.Ident).String()] = struct{}{}
						}
					}
				}
				return true
			})
			var buf bytes.Buffer
			printer.Fprint(&buf, fs, expr)
			f.Code = buf.String()
			fileTemplate.New(f.Name).Parse(f.Code)
			arguments := make([]string, highest)
			for i := range arguments {
				arguments[i] = fmt.Sprintf("%c Number", 'b'+i)
			}
			f.Args = strings.Join(arguments, ", ")
			functions = append(functions, f)
		}
	}

	fileTemplate.Execute(f, struct {
		Imports   map[string]struct{}
		Types     []string
		Functions []function
	}{
		imports,
		types,
		functions,
	})
}
