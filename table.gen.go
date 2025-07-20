package cpu

var Rules = []DecodingRule{
	DecodingRule{
		Mnemonic: 0x1,
		Bytes: [6]ByteDecoding{
			ByteDecoding{
				NotEmpty: true,
				Parts: [3]Part{
					Part{
						NotEmpty: true,
						Kind:     0xa,
						Mask:     63,
						Shift:    2,
						Literal:  34},
					Part{
						NotEmpty: true,
						Kind:     0x5,
						Mask:     1,
						Shift:    1,
						Literal:  -1},
					Part{
						NotEmpty: true,
						Kind:     0x6,
						Mask:     1,
						Shift:    0,
						Literal:  -1,
					},
				},
				Cond: 0x0},
			ByteDecoding{
				NotEmpty: true,
				Parts: [3]Part{
					Part{
						NotEmpty: true,
						Kind:     0x2,
						Mask:     3,
						Shift:    6,
						Literal:  -1},
					Part{
						NotEmpty: true,
						Kind:     0x4,
						Mask:     7,
						Shift:    3,
						Literal:  -1},
					Part{
						NotEmpty: true,
						Kind:     0x3,
						Mask:     7,
						Shift:    0,
						Literal:  -1,
					},
				},
				Cond: 0x0},
			ByteDecoding{
				NotEmpty: true,
				Parts: [3]Part{
					Part{
						NotEmpty: true,
						Kind:     0x7,
						Mask:     255,
						Shift:    0,
						Literal:  -1},
					Part{
						NotEmpty: false,
						Kind:     0x0,
						Mask:     0,
						Shift:    0,
						Literal:  0},
					Part{
						NotEmpty: false,
						Kind:     0x0,
						Mask:     0,
						Shift:    0,
						Literal:  0,
					},
				},
				Cond: 0x0},
			ByteDecoding{
				NotEmpty: true,
				Parts: [3]Part{
					Part{
						NotEmpty: true,
						Kind:     0x8,
						Mask:     255,
						Shift:    0,
						Literal:  -1},
					Part{
						NotEmpty: false,
						Kind:     0x0,
						Mask:     0,
						Shift:    0,
						Literal:  0},
					Part{
						NotEmpty: false,
						Kind:     0x0,
						Mask:     0,
						Shift:    0,
						Literal:  0,
					},
				},
				Cond: 0x0},
			ByteDecoding{
				NotEmpty: false,
				Parts: [3]Part{
					Part{
						NotEmpty: false,
						Kind:     0x0,
						Mask:     0,
						Shift:    0,
						Literal:  0},
					Part{
						NotEmpty: false,
						Kind:     0x0,
						Mask:     0,
						Shift:    0,
						Literal:  0},
					Part{
						NotEmpty: false,
						Kind:     0x0,
						Mask:     0,
						Shift:    0,
						Literal:  0,
					},
				},
				Cond: 0x0},
			ByteDecoding{
				NotEmpty: false,
				Parts: [3]Part{
					Part{
						NotEmpty: false,
						Kind:     0x0,
						Mask:     0,
						Shift:    0,
						Literal:  0},
					Part{
						NotEmpty: false,
						Kind:     0x0,
						Mask:     0,
						Shift:    0,
						Literal:  0},
					Part{
						NotEmpty: false,
						Kind:     0x0,
						Mask:     0,
						Shift:    0,
						Literal:  0,
					},
				},
				Cond: 0x0}}},
}
