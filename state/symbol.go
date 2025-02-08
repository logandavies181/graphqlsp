package state

import "github.com/vektah/gqlparser/v2/ast"

type symbol interface {
	*ast.Definition | *ast.ArgumentDefinition | *ast.FieldDefinition | *ast.Type
}

func lift[S symbol](in any) S {
	s, ok := in.(S)
	if !ok {
		return nil
	}
	return s
}
