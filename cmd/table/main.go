package main

import (
	_ "embed"
	"strings"
)

//go:embed table.sim8086
var table string

func main() {
	for line := range strings.Lines(table) {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, ";") {
			continue
		}

		// fmt.Printf("%+v\n", rule)
	}
}
