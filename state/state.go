package state

import (
	"os"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type State struct {
	schema  *ast.Schema
	locator locator
}

type Position struct {
	Line int
	Col  int
}

func NewFromFile(fname string) (*State, error) {
	dat, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	source := ast.Source{
		Name:  fname,
		Input: string(dat),
	}
	schema, err := gqlparser.LoadSchema(&source)
	if err != nil {
		return nil, err
	}

	return &State{
		schema:    schema,
		locator:   make(locator),
	}, nil
}

func (s *State) GetDefinitionOf(line, col int) *Position {
	sym := s.locator.get(line, col)
	ty := lift[*ast.Type](sym)
	if ty == nil {
		return nil
	}

	defType, ok := s.schema.Types[ty.Name()]
	if !ok {
		return nil
	}

	return &Position{
		Line: defType.Position.Line,
		Col: defType.Position.Column,
	}
}

func (s *State) handleType(ty *ast.Type) {
	s.locator.push(location{
		start: ty.Position.Column,
		end: ty.Position.End - ty.Position.End,
	}, ty.Position.Line)
}

func (s *State) handleField(ty *ast.FieldDefinition) {
	s.locator.push(location{
		start: ty.Position.Column,
		end: ty.Position.End - ty.Position.End,
	}, ty.Position.Line)
}

func (s *State) handleDef(ty *ast.Definition) {
	s.locator.push(location{
		start: ty.Position.Column,
		end: ty.Position.End - ty.Position.End,
	}, ty.Position.Line)
}

func (s *State) handleArg(ty *ast.ArgumentDefinition) {
	s.locator.push(location{
		start: ty.Position.Column,
		end: ty.Position.End - ty.Position.End,
	}, ty.Position.Line)
}

func (s *State) walk(def *ast.Definition) {
	switch def.Kind {
	case ast.Scalar:
		s.handleDef(def)
	case ast.Object:
		s.walkObj(def)
	}
}

func (s *State) walkObj(def *ast.Definition) {
	if def.Name[0:2] == "__" {
		// no doing stuff with meta types which get populated into the AST
		return
	}

	s.handleDef(def)
	for _, v := range def.Fields {
		if v != nil {
			s.walkField(v)
		}
	}
}

func (s *State) walkField(def *ast.FieldDefinition) {
	s.handleField(def)
	s.walkFieldArgs(def.Arguments)
	s.handleType(def.Type)
	if def.Type.Elem != nil {
		s.walkArray(def.Type.Elem)
	}
}

func (s *State) walkArray(ty *ast.Type) {
	if ty.Elem != nil {
		s.walkArray(ty.Elem)
	}
}

func (s *State) walkFieldArgs(args ast.ArgumentDefinitionList) {
	for _, v := range args {
		s.walkFieldArg(v)
	}
}

func (s *State) walkFieldArg(arg *ast.ArgumentDefinition) {
	s.handleArg(arg)
	s.handleType(arg.Type)
	if arg.Type.Elem != nil {
		s.walkArray(arg.Type.Elem)
	}
}
