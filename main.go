package main

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/logandavies181/graphqlsp/state"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"github.com/vektah/gqlparser/v2/validator"
)

var (
	version string = "0.0.1"
	handler protocol.Handler
	tempDir string
	states  map[string]*state.State
)

func main() {
	handler = protocol.Handler{
		TextDocumentDefinition: definition,
		TextDocumentHover:      hover,
		Initialize:             initialize,
		Shutdown:               shutdown,
	}

	server := server.NewServer(&handler, "testls", true)

	server.RunStdio()
}

func loadFile(uri string) (*state.State, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("could not parse file uri: %w", err)
	}

	path := url.Path
	s, ok := states[path]
	if ok {
		return s, nil
	}

	s, err = state.NewFromFile(url.Path)
	if err != nil {
		return nil, fmt.Errorf("could not load state from file %s: %w", url.Path, err)
	}

	return s, nil
}

func preludeFilePath() string {
	return path.Join(tempDir, "prelude.graphql")
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	var err error
	tempDir, err = os.MkdirTemp("", "graphqlsp")
	if err != nil {
		return nil, fmt.Errorf("could not create temp dir for prelude: %w", err)
	}

	err = os.WriteFile(preludeFilePath(), []byte(validator.Prelude.Input), 0700)
	if err != nil {
		return nil, fmt.Errorf("could not write to temp file for prelude: %w", err)
	}

	states[preludeFilePath()] = state.PreludeState()

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
	return os.RemoveAll(tempDir)
}

func definition(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	s, err := loadFile(params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	pos := s.GetDefinitionOf(int(params.Position.Line)+1, int(params.Position.Character)+1)
	if pos == nil {
		return nil, nil
	}

	file := params.TextDocument.URI
	if pos.Prelude {
		file = "file://" + preludeFilePath()
	}

	return protocol.Location{
		URI: file,
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      uint32(pos.Line - 1),
				Character: uint32(pos.Col - 1),
			},
			End: protocol.Position{
				Line:      uint32(pos.Line - 1),
				Character: uint32(pos.Col - 1 + pos.Len),
			},
		},
	}, nil
}

func hover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	s, err := loadFile(params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	mu, pos := s.GetHoverOf(int(params.Position.Line)+1, int(params.Position.Character)+1)
	if pos == nil {
		return nil, nil
	}

	return &protocol.Hover{
		Contents: mu,
		Range: &protocol.Range{
			Start: protocol.Position{
				Line:      uint32(pos.Line - 1),
				Character: uint32(pos.Col - 1),
			},
			End: protocol.Position{
				Line:      uint32(pos.Line - 1),
				Character: uint32(pos.Col - 1 + pos.Len),
			},
		},
	}, nil
}
