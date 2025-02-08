package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/logandavies181/graphqlsp/state"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
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

func definition(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	fmt.Fprintln(os.Stderr, params.TextDocument.URI)
	fmt.Fprintln(os.Stderr, params.Position.Line)
	fmt.Fprintln(os.Stderr, params.Position.Character)

	url, err := url.Parse(params.TextDocument.URI)
	if err != nil {
		return nil, fmt.Errorf("could not parse file uri", err)
	}

	s, err := state.NewFromFile(url.Path)
	if err != nil {
		return nil, fmt.Errorf("could not load state from file %s: %w", url.Path, err)
	}

	pos := s.GetDefinitionOf(int(params.Position.Line) + 1, int(params.Position.Character) + 1)
	if pos == nil {
		return nil, nil
	}

	return protocol.Location{
		URI: params.TextDocument.URI,
		Range: protocol.Range{
			Start: protocol.Position{
				Line: uint32(pos.Line - 1),
				Character: uint32(pos.Col - 1),
			},
			End: protocol.Position{
				Line: uint32(pos.Line - 1),
				Character: uint32(pos.Col - 1 + pos.Len),
			},
		},
	}, nil
}
