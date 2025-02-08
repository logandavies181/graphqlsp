package state

import (
	"fmt"
	"os"
)

type location struct {
	s       any
	start   int
	end     int
	prelude bool
}

type locator map[int][]location

func (l locator) push(loc location, line int) {
	if overlaps(loc, l[line]) {
		return
	}

	l[line] = append(l[line], loc)
}

func (l locator) get(line, col int) any {
	if val, ok := l[line]; ok {
		for _, v := range val {
			if col >= v.start && col <= v.end {
				fmt.Fprintf(os.Stderr, "found: %v\n", v.s)
				return v.s
			}
		}
	}
	return nil
}

func overlaps(loc location, locs []location) bool {
	for _, v := range locs {
		if loc.start >= v.start && loc.end <= v.end {
			return true
		}
	}

	return false
}
