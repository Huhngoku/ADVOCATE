package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a directory path")
		return
	}

	dirPath := os.Args[1]
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			instrumentFile(path)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func instrumentFile(filePath string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	spawnFunction := "runtime.AdvocateSpawn"
	header := `if true {
		runtime.InitAdvocate(0)
		defer advocate.CreateTrace("advocateTrace.log")
	} else {
		trace := advocate.ReadTrace("advocateTrace.log")
		runtime.EnableReplay(trace)
		defer runtime.WaitForReplayFinish()
	}`

	addRuntimeImport := false
	addAdvocateImport := false

	ast.Inspect(node, func(n ast.Node) bool {
		// instrument spawn statement
		goStmt, ok := n.(*ast.GoStmt)
		if ok {
			addRuntimeImport = true
			switch callExpr := goStmt.Call.Fun.(type) {
			case *ast.Ident:
				params := ""
				for i, arg := range goStmt.Call.Args {
					if i > 0 {
						params += ", "
					}
					switch arg := arg.(type) {
					case *ast.Ident:
						params += arg.Name
					case *ast.BasicLit:
						params += arg.Value
					}
				}
				callExpr.Name = "func() {" + spawnFunction + "(); " + callExpr.Name + "(" + params + ")}"
				goStmt.Call.Args = nil
			case *ast.FuncLit:
				callExpr.Body.List = append([]ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{Fun: ast.NewIdent(spawnFunction)}}}, callExpr.Body.List...)
			}
			return true
		}
		// add header in main function
		funcDecl, ok := n.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == "main" {
			addAdvocateImport = true
			newStmts := []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.Ident{
						Name: header,
					},
				},
				// Add more statements here...
			}

			// Prepend the new statements to the body of the main function
			funcDecl.Body.List = append(newStmts, funcDecl.Body.List...)
		}
		return true
	})

	// add runtime import if needed
	if addRuntimeImport {
		// Check if runtime package is already imported
		runtimeImported := false
		for _, imp := range node.Imports {
			if imp.Path.Value == "\"runtime\"" {
				runtimeImported = true
				break
			}
		}

		// If not, add it
		if !runtimeImported {
			node.Imports = append(node.Imports, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"runtime\"",
				},
			})
		}
	}
	if addAdvocateImport {
		// Check if runtime package is already imported
		runtimeImported := false
		for _, imp := range node.Imports {
			if imp.Path.Value == "\"advocate\"" {
				runtimeImported = true
				break
			}
		}

		// If not, add it
		if !runtimeImported {
			node.Imports = append(node.Imports, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"advocate\"",
				},
			})
		}
	}

	err = os.WriteFile(filePath, []byte{}, 0644)
	if err != nil {
		fmt.Println("Error clearing file:", err)
		return
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	err = printer.Fprint(file, fset, node)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
