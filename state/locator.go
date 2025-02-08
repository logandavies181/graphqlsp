package state

import (
	"fmt"
	"os"

	"github.com/vektah/gqlparser/v2/ast"
)

type symbol interface {
	*ast.Definition | *ast.FieldDefinition | *ast.Type | *ast.ArgumentDefinition
}

type location[S symbol] struct {
	s     S
	start int
	end   int
}

type locator[S symbol] map[int][]location[S]

func (l locator[S]) push(loc location[S], line int) {
	if overlaps(loc, l[line]) {
		fmt.Fprintf(os.Stderr, "overlapping symbol, cannot push")
		return
	}

	l[line] = append(l[line], loc)
}

func (l locator[S]) get(line, col int) S {
	if val, ok := l[line]; ok {
		for _, v := range val {
			if col >= v.start && col <= v.end {
				return v.s
			}
		}
	}
	return nil
}

func overlaps[S symbol](loc location[S], locs []location[S]) bool {
	for _, v := range locs {
		if loc.start >= v.start && loc.end <= v.end {
			return true
		}
	}

	return false
}
