package main

import (
	"encoding/gob"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
)

func init() {
	gob.Register(&ast.GenDecl{})
	gob.Register(&ast.Field{})
	gob.Register(&ast.FieldList{})
	gob.Register(&ast.Ident{})
	gob.Register(&ast.StructType{})
	gob.Register(&ast.TypeSpec{})
}

// createFunc создает AST узел для функции.
func createFunc(funcName, structName string) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(funcName),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("obj")},
						Type: &ast.StarExpr{
							X: ast.NewIdent(structName),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func main() {
	// Загрузка AST структуры из файла
	fileAst := "struct_ast.gob"
	inputFile, err := os.Open(fileAst)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	var genDecl ast.GenDecl
	decoder := gob.NewDecoder(inputFile)
	err = decoder.Decode(&genDecl)
	if err != nil {
		panic(err)
	}

	structName := "Example" // Имя структуры должно совпадать с тем, что было использовано ранее

	// Создание функций репозитория
	funcs := []*ast.FuncDecl{
		createFunc("Create"+structName, structName),
		createFunc("Read"+structName, structName),
		createFunc("Update"+structName, structName),
		createFunc("Delete"+structName, structName),
	}

	// Создание файла
	file := &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: []ast.Decl{&genDecl},
	}

	// Добавление функций к файлу
	for _, f := range funcs {
		file.Decls = append(file.Decls, f)
	}

	// Сохранение в файл
	fset := token.NewFileSet()
	outputFileName := "generated_repository.go"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = printer.Fprint(outputFile, fset, file)
	if err != nil {
		panic(err)
	}

	println("CRUD репозиторий успешно сгенерирован и сохранен в " + outputFileName)
}
