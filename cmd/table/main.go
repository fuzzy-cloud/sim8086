package main

import (
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

	os.WriteFile("../../table.gen.go", []byte(code), 0o600)
	fmt.Println("done")
}

func generateCode() (string, error) {
	var b strings.Builder

	fmt.Fprintf(&b, "package cpu\n\n")
	fmt.Fprintf(&b, "var Rules = []DecodingRule{\n")

	for line := range strings.Lines(rawTable) {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, ";") {
			continue
		}

		rule, err := ParseDecodingRule(line)
		if err != nil {
			err = fmt.Errorf("failed to parse a deconding rule.\nline: %q\nerr: %w", line, err)
			return "", err
		}

		ruleStr := fmt.Sprintf("%#v", rule)
		ruleStr = strings.ReplaceAll(ruleStr, "cpu.", "")

		// omit excessive types
		// ruleStr = strings.ReplaceAll(ruleStr, "\tDecodingRule{", "\t{")
		// ruleStr = strings.ReplaceAll(ruleStr, "\tByteDecoding{", "\t{")
		// ruleStr = strings.ReplaceAll(ruleStr, "\tPart{", "\t{")

		// trigger formatting
		ruleStr = strings.ReplaceAll(ruleStr, "{", "{\n")
		ruleStr = strings.ReplaceAll(ruleStr, ", ", ",\n")
		ruleStr = strings.ReplaceAll(ruleStr, "}},", ",\n},\n},")

		fmt.Fprintf(&b, "\t%s,\n", ruleStr)
	}

	fmt.Fprintf(&b, "}\n\n")

	formatted, err := format.Source([]byte(b.String()))
	if err != nil {
		err = fmt.Errorf("failed to format a generated code.\n\ncode: %q\n\nerr: %w", b.String(), err)
		return "", err
	}

	return string(formatted), nil
}
