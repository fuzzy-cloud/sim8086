package cpu

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	bs, err := os.ReadFile("listings/0037_single_register_mov")
	require.NoError(t, err)
	require.Len(t, bs, 2)
}
