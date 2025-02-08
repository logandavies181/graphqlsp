package state

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
)


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
