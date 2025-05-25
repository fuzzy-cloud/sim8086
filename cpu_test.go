package cpu

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	filenames := []string{
		"0037_single_register_mov",
		"0038_many_register_mov",
		"0039_more_movs",
		"0040_challenge_movs",
		"0041_add_sub_cmp_jnz",
	}

	for _, name := range filenames {
		t.Run(name, func(t *testing.T) {
			tmpDir := t.TempDir()

			stream, err := os.ReadFile("listings/" + name)
			require.NoError(t, err)

			asmFile, err := os.ReadFile("listings/" + name + ".asm")
			require.NoError(t, err)

			filtered := make([]string, 0, 10)
			shouldAddNewLine := false
			for l := range strings.Lines(string(asmFile)) {
				if l[0] == ';' {
					continue
				}
				if l[0] == '\n' {
					if shouldAddNewLine {
						shouldAddNewLine = false
					} else {
						continue
					}
				}
				if l == "bits 16\n" {
					shouldAddNewLine = true
				}
				filtered = append(filtered, l)
			}
			want := strings.Join(filtered, "")
			want = strings.TrimSpace(want)

			got, err := disassemble(stream)
			require.NoError(t, err)

			tmpFileAsmPath := path.Join(tmpDir, name+".asm")
			err = os.WriteFile(tmpFileAsmPath, []byte(got), 0o666)
			require.NoError(t, err)

			cmd := exec.Command("nasm", tmpFileAsmPath)
			err = cmd.Run()
			require.NoError(t, err)

			tmpFileBinPath := path.Join(tmpDir, name)
			binary, err := os.ReadFile(tmpFileBinPath)

			require.NoError(t, err)
			if !assert.Equal(t, stream, binary) {
				assert.Equal(t, want, got)
			}
		})
	}
}
