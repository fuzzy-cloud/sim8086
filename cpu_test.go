package cpu

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	filenames := []string{
		"0037_single_register_mov",
		"0038_many_register_mov",
	}

	for _, name := range filenames {
		t.Run(name, func(t *testing.T) {
			stream, err := os.ReadFile("listings/" + name)
			require.NoError(t, err)

			asmFile, err := os.ReadFile("listings/" + name + ".asm")
			require.NoError(t, err)

			insideComment := false
			want := strings.TrimLeftFunc(string(asmFile), func(r rune) bool {
				if r == ';' {
					insideComment = true
				}
				if r == '\n' {
					insideComment = false
					return true
				}
				return insideComment
			})
			want = strings.TrimSpace(want)

			t.Logf("\n%b\n", stream)

			got, err := disassemble(stream)
			require.NoError(t, err)
			require.Equal(t, want, got)
		})
	}
}
