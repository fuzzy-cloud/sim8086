package main

import (
	cpu "cpu8086"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodingRule(t *testing.T) {
	const input = `MOV | 100010 d w | mod reg rm | disp-lo | disp-hi`

	got, err := parseDecodingRule(input)

	require.NoError(t, err)
	require.Equal(t, cpu.MOV, got.Mnemonic)
	require.Equal(t, 4, countNotEmptyBytes(got.Bytes))
}
