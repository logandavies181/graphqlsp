package state

import "github.com/vektah/gqlparser/v2/ast"

type symbol interface {
	*ast.Definition | *ast.ArgumentDefinition | *ast.FieldDefinition | *ast.Type
}

func lift[S symbol](in any) S {
	return in.(S)
}
