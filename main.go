package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)


var version string = "0.0.1"
var handler protocol.Handler

func main() {
	handler = protocol.Handler{
		TextDocumentDefinition: definition,
		Initialize: initialize,
		Shutdown:   shutdown,
	}

	server := server.NewServer(&handler, "testls", true)

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	capabilities.CompletionProvider = &protocol.CompletionOptions{}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    "testls",
			Version: &version,
		},
	}, nil
}

func shutdown(context *glsp.Context) error {
	return nil
}

type pos struct {
	line int
	col int
	len int

	name string
}

func getFromASTAt(uri string, line, char uint32) (pos, error) {
	_url, err := url.Parse(uri)
	if err != nil {
		return pos{}, err
	}

	dat, err := os.ReadFile(_url.Path)
	if err != nil {
		return pos{}, err
	}

	source := ast.Source{
		Name: "schema.graphql",
		Input: string(dat),
	}
	schema, err := gqlparser.LoadSchema(&source)
	if err != nil {
		return pos{}, err
	}

	locs := []pos{}
	for _, v := range schema.Types {
		vpos := v.Position
		locs = append(locs, pos{
			line: vpos.Line,
			col: vpos.Column,
			len: vpos.End - vpos.Start,
			name: v.Name,
		})
	}

	isWithin := func(line, char int, pos pos) bool {
		if line != pos.line {
			return false
		}

		if char < pos.col || char > pos.col + pos.len {
			return false
		}

		return true
	}

	for _, v := range locs {
		if !isWithin(int(line), int(char), v) {
			continue
		}

		v.
	}

	return pos{}, nil
}

func definition(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	fmt.Fprintln(os.Stderr, params.TextDocument.URI)
	fmt.Fprintln(os.Stderr, params.Position.Line)
	fmt.Fprintln(os.Stderr, params.Position.Character)
	pos, err := getFromASTAt(params.TextDocument.URI, params.Position.Line, params.Position.Character)
	if err != nil {
		return nil, err
	}

	return protocol.Location{
		URI: params.TextDocument.URI,
		Range: protocol.Range{
			Start: protocol.Position{
				Line: uint32(pos.line),
				Character: uint32(pos.col),
			},
			End: protocol.Position{
				Line: uint32(pos.line),
				Character: uint32(pos.col + pos.len),
			},
		},
	}, nil
}
