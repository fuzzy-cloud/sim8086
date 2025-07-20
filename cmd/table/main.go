package main

import (
	"cpu8086/table"
	_ "embed"
	"fmt"
	"go/format"
	"os"
	"strings"
)

//go:embed table.sim8086
var rawTable string

func main() {
	code, err := generateCode()
	if err != nil {
		panic(err)
	}

	os.WriteFile("../../table/table.gen.go", []byte(code), 0o600)
	fmt.Println("done")
}

func generateCode() (string, error) {
	var b strings.Builder

	fmt.Fprintf(&b, "package table\n\n")
	fmt.Fprintf(&b, "import cpu \"cpu8086\"\n\n")
	fmt.Fprintf(&b, "var Rules = []cpu.DecodingRule{\n")

	for line := range strings.Lines(rawTable) {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, ";") {
			continue
		}

		rule, err := table.ParseDecodingRule(line)
		if err != nil {
			err = fmt.Errorf("failed to parse a deconding rule.\nline: %q\nerr: %w", line, err)
			return "", err
		}

		fmt.Fprintf(&b, "\t%#v,\n", rule)
	}

	fmt.Fprintf(&b, "}\n\n")

	formatted, err := format.Source([]byte(b.String()))
	if err != nil {
		err = fmt.Errorf("failed to format a generated code.\n\ncode: %q\n\nerr: %w", b.String(), err)
		return "", err
	}

	return string(formatted), nil
}
