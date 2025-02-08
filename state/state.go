package state

import (
	"os"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type state struct {
	schema *ast.Schema
	locator locator

	typeFunc func(*ast.Type)
	fieldFunc func(*ast.FieldDefinition)
	defFunc func(*ast.Definition)
	argFunc func(*ast.ArgumentDefinition)
}

func newFromFile(fname string) (*state, error) {
	dat, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	source := ast.Source{
		Name: fname,
		Input: string(dat),
	}
	schema, err := gqlparser.LoadSchema(&source)
	if err != nil {
		return nil, err
	}

	return &state{
		schema: schema,
		typeFunc: handleType,
		fieldFunc: handleField,
		defFunc: handleDef,
		argFunc: handleArg,
	}, nil
}

func handleType(ty *ast.Type) {

}

func handleField(ty *ast.FieldDefinition) {

}

func handleDef(ty *ast.Definition) {

}

func handleArg(ty *ast.ArgumentDefinition) {

}

func (s *state) walk(def *ast.Definition) {
	switch def.Kind {
	case ast.Scalar:
		s.defFunc(def)
	case ast.Object:
		s.walkObj(def)
	}
}

func (s *state) walkObj(def *ast.Definition) {
	if def.Name[0:2] == "__" {
		// no doing stuff with meta types which get populated into the AST
		return
	}

	s.defFunc(def)
	for _, v := range def.Fields {
		if v != nil {
			s.walkField(v)
		}
	}
}

func (s *state) walkField(def *ast.FieldDefinition) {
	s.fieldFunc(def)
	s.walkFieldArgs(def.Arguments)
	s.typeFunc(def.Type)
	if def.Type.Elem != nil {
		s.walkArray(def.Type.Elem)
	}
}

func (s *state) walkArray(ty *ast.Type) {
	if ty.Elem != nil {
		s.walkArray(ty.Elem)
	}
}

func (s *state) walkFieldArgs(args ast.ArgumentDefinitionList) {
	for _, v := range args {
		s.walkFieldArg(v)
	}
}

func (s *state) walkFieldArg(arg *ast.ArgumentDefinition) {
	s.argFunc(arg)
	s.typeFunc(arg.Type)
	if arg.Type.Elem != nil {
		s.walkArray(arg.Type.Elem)
	}
}
