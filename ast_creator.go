package main

import (
	"encoding/gob"
	"go/ast"
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

// createStruct создает AST узел для структуры данных с полями, переданными снаружи.
func createStruct(structName string, fields map[string]string) *ast.GenDecl {
	var fieldList []*ast.Field

	for name, typ := range fields {
		fieldList = append(fieldList, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(name)},
			Type:  ast.NewIdent(typ),
		})
	}

	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(structName),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: fieldList,
					},
				},
			},
		},
	}
}

func main() {
	structName := "Example"
	fields := map[string]string{
		"ID":   "int",
		"Name": "string",
		"Age":  "int",
	}

	astNode := createStruct(structName, fields)

	// Создание файла для сохранения AST
	fileAst := "struct_ast.gob"
	outputFile, err := os.Create(fileAst)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	encoder := gob.NewEncoder(outputFile)
	err = encoder.Encode(astNode)
	if err != nil {
		panic(err)
	}

	println("AST структуры успешно сохранено в " + fileAst)
}
