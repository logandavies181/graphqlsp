package state

import (
	"fmt"
	"os"

	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type State struct {
	schema  *ast.Schema
	locator locator
}

type Position struct {
	Line    int
	Col     int
	Len     int
	Prelude bool
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

	state := &State{
		schema:  schema,
		locator: make(locator),
	}

	state.walk(schema.Query, false)
	state.walk(schema.Mutation, false)
	for _, v := range schema.Types {
		state.walk(v, false)
	}

	return state, nil
}

func PreludeState() *State {
	source := ast.Source{
		Name:  "graphqlsp_internal",
		Input: "",
	}
	schema, err := gqlparser.LoadSchema(&source)
	if err != nil {
		panic(err)
	}

	state := &State{
		schema:  schema,
		locator: make(locator),
	}

	state.walk(schema.Query, true)
	state.walk(schema.Mutation, true)
	for _, v := range schema.Types {
		state.walk(v, true)
	}

	return state
}

func (s *State) GetDefinitionOf(line, col int) *Position {
	sym := s.locator.get(line, col)
	typeName := ""
	switch symTy := sym.(type) {
	case *ast.Type:
		typeName = symTy.Name()
	case *ast.Definition:
		typeName = symTy.Name
	default:
		return nil
	}

	defType, ok := s.schema.Types[typeName]
	if !ok {
		return nil
	}

	return &Position{
		Line:    defType.Position.Line,
		Col:     defType.Position.Column,
		Len:     defType.Position.End - defType.Position.Start,
		Prelude: defType.BuiltIn,
	}
}

func defKindToKeyword(kind ast.DefinitionKind) string {
	switch kind {
	case ast.Scalar:
		return "scalar"
	case ast.Object:
		return "type"
	case ast.Enum:
		return "enum"
	case ast.Union:
		return "union"
	case ast.Interface:
		return "interface"
	case ast.InputObject:
		return "input"
	default:
		return ""
	}
}

func formatDescriptionMarkdown(def ast.Definition) string {
	return fmt.Sprintf("```graphql\n%s %s\n```\n%s", defKindToKeyword(def.Kind), def.Name, def.Description)
}

func (s *State) GetHoverOf(line, col int) (*protocol.MarkupContent, *Position) {
	sym := s.locator.get(line, col)
	typeName := ""
	switch ty := sym.(type) {
	case *ast.Type:
		if ty == nil {
			return nil, nil
		}
		typeName = ty.Name()
	case *ast.Definition:
		typeName = ty.Name
	default:
		return nil, nil
	}

	defType, ok := s.schema.Types[typeName]
	if !ok {
		return nil, nil
	}

	mu := protocol.MarkupContent{
		Kind:  protocol.MarkupKindMarkdown,
		Value: formatDescriptionMarkdown(*defType),
	}

	return &mu, &Position{
		Line:    defType.Position.Line,
		Col:     defType.Position.Column,
		Len:     defType.Position.End - defType.Position.Start,
		Prelude: defType.BuiltIn,
	}
}

func (s *State) handleType(ty *ast.Type) {
	if ty.Position == nil {
		return
	}

	s.locator.push(location{
		s:     ty,
		start: ty.Position.Column,
		end:   ty.Position.Column + ty.Position.End - ty.Position.Start,
	}, ty.Position.Line)
}

func (s *State) handleField(ty *ast.FieldDefinition) {
	if ty.Position == nil {
		return
	}

	s.locator.push(location{
		s:     ty,
		start: ty.Position.Column,
		end:   ty.Position.Column + ty.Position.End - ty.Position.Start,
	}, ty.Position.Line)
}

func (s *State) handleDef(ty *ast.Definition) {
	if ty.Position == nil {
		return
	}

	s.locator.push(location{
		s:       ty,
		start:   ty.Position.Column,
		end:     ty.Position.Column + ty.Position.End - ty.Position.Start,
		prelude: ty.Position.Src.BuiltIn,
	}, ty.Position.Line)
}

func (s *State) handleArg(ty *ast.ArgumentDefinition) {
	if ty.Position == nil {
		return
	}

	s.locator.push(location{
		s:     ty,
		start: ty.Position.Column,
		end:   ty.Position.Column + ty.Position.End - ty.Position.Start,
	}, ty.Position.Line)
}

func (s *State) walk(def *ast.Definition, builtinmode bool) {
	if def == nil || (def.BuiltIn && !builtinmode) {
		return
	}

	switch def.Kind {
	case ast.Scalar:
		s.handleDef(def)
	case ast.Object:
		s.walkObj(def)
	}
}

func (s *State) walkObj(def *ast.Definition) {
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
