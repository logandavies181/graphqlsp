package main

import (
	"fmt"
	"os"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	err := mainE()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func walk(def *ast.Definition) {
	switch def.Kind {
	case ast.Scalar:
		printPos(def)
	case ast.Object:
		walkObj(def)
	}
}

func walkObj(def *ast.Definition) {
	if def.Name[0:2] == "__" {
		// no doing stuff with meta types which get populated into the AST
		return
	}

	printPos(def)
	for _, v := range def.Fields {
		if v != nil {
			walkField(v)
		}
	}
}

func walkField(def *ast.FieldDefinition) {
	printPosField(def)
	walkFieldArgs(def.Arguments)
	printPosFieldType(def.Type)
	if def.Type.Elem != nil {
		walkArray(def.Type.Elem)
	}
}

func walkArray(ty *ast.Type) {
	if ty.Elem != nil {
		walkArray(ty.Elem)
	}
}

func walkFieldArgs(args ast.ArgumentDefinitionList) {
	for _, v := range args {
		walkFieldArg(v)
	}
}

func walkFieldArg(arg *ast.ArgumentDefinition) {
	printPosArg(arg)
	printPosArgType(arg.Type)
	if arg.Type.Elem != nil {
		walkArray(arg.Type.Elem)
	}
}

func printPos(def *ast.Definition) {
	pos := def.Position
	fmt.Printf("top level type %s: %s on line %d, column %d\n", def.Name, def.Kind, pos.Line, pos.Column)
}

func printPosField(def *ast.FieldDefinition) {
	pos := def.Position
	if pos == nil {
		// meta types not actually present in the file are here
		return
	}
	fmt.Printf("field: %s: on line %d, column %d\n", def.Name, pos.Line, pos.Column)
}

func printPosArg(def *ast.ArgumentDefinition) {
	pos := def.Position
	if pos == nil {
		// meta types not actually present in the file are here
		return
	}
	fmt.Printf("arg: %s: on line %d, column %d\n", def.Name, pos.Line, pos.Column)
}

func printPosFieldType(ty *ast.Type) {
	if ty.Position != nil && ty.NamedType != "" {
		fmt.Printf("field type: %s: on line %d, column %d\n", ty.NamedType, ty.Position.Line, ty.Position.Column)
	}
}

func printPosArgType(ty *ast.Type) {
	if ty.Position != nil && ty.NamedType != "" {
		fmt.Printf("arg type: %s: on line %d, column %d\n", ty.NamedType, ty.Position.Line, ty.Position.Column)
	}
}

func printArray(ty *ast.Type) {
	if ty.NamedType != "" {
		// arrays don't get names
		fmt.Printf("array of %s: on line %d, column %d \n", ty.NamedType, ty.Position.Line, ty.Position.Column)
	}
}

func mainE() error {
	dat, err := os.ReadFile("test/schema.graphql")
	if err != nil {
		return err
	}

	source := ast.Source{
		Name: "schema.graphql",
		Input: string(dat),
	}
	schema, err := gqlparser.LoadSchema(&source)
	if err != nil {
		return err
	}

	walkObj(schema.Query)
	walkObj(schema.Mutation)
	for _, v := range schema.Types {
		walkObj(v)
	}

	return nil
}
