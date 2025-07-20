package main

import (
	cpu "cpu8086"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodingRule(t *testing.T) {
	cases := []struct {
		input string
		want  cpu.DecodingRule
	}{
		{
			input: "MOV | 100010 d w | mod reg rm | disp-lo | disp-hi",
			want: cpu.DecodingRule{
				Mnemonic: cpu.MOV,
				Bytes: [6]cpu.ByteDecoding{
					{
						NotEmpty: true,
						Parts: [3]cpu.Part{
							{
								NotEmpty: true,
								Kind:     cpu.PartLiteral,
								Mask:     0b111_111,
								Shift:    2,
								Literal:  0b100010,
							},
							{
								NotEmpty: true,
								Kind:     cpu.PartD,
								Mask:     0b1,
								Shift:    1,
								Literal:  -1,
							},
							{
								NotEmpty: true,
								Kind:     cpu.PartW,
								Mask:     0b1,
								Shift:    0,
								Literal:  -1,
							},
						},
					},
					{
						NotEmpty: true,
						Parts: [3]cpu.Part{
							{
								NotEmpty: true,
								Kind:     cpu.PartMOD,
								Mask:     0b11,
								Shift:    6,
								Literal:  -1,
							},
							{
								NotEmpty: true,
								Kind:     cpu.PartREG,
								Mask:     0b111,
								Shift:    3,
								Literal:  -1,
							},
							{
								NotEmpty: true,
								Kind:     cpu.PartRM,
								Mask:     0b111,
								Shift:    0,
								Literal:  -1,
							},
						},
					},
					{
						NotEmpty: true,
						Parts: [3]cpu.Part{
							{
								NotEmpty: true,
								Kind:     cpu.PartDISP_LO,
								Mask:     0b1111_1111,
								Shift:    0,
								Literal:  -1,
							},
							{}, // empty
							{}, // empty
						},
					},
					{
						NotEmpty: true,
						Parts: [3]cpu.Part{
							{
								NotEmpty: true,
								Kind:     cpu.PartDISP_HI,
								Mask:     0b1111_1111,
								Shift:    0,
								Literal:  -1,
							},
							{}, // empty
							{}, // empty
						},
					},
					{}, // empty
					{}, // empty
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := parseDecodingRule(tc.input)

			require.NoError(t, err)
			require.Equal(t, tc.want.Mnemonic, got.Mnemonic)

			require.Equal(t, tc.want.Bytes[0], got.Bytes[0])
			require.Equal(t, tc.want.Bytes[1], got.Bytes[1])
			require.Equal(t, tc.want.Bytes[2], got.Bytes[2])
			require.Equal(t, tc.want.Bytes[3], got.Bytes[3])
			require.Equal(t, tc.want.Bytes[4], got.Bytes[4])
			require.Equal(t, tc.want.Bytes[5], got.Bytes[5])
		})
	}
}
