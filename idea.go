package idea

import (
	"encoding/binary"
	"fmt"
)

const BlockSize = 8
const KeySize = 16
const subBlockQty = 52
const keyShifts = 25

type IdeaCipher struct {
	key          []byte
	encSubBlocks []uint16
	decSubBlocks []uint16
	TestMode     bool
}

func printRound(rn int, rnd []uint16) {
	if rn == 0 {
		fmt.Printf("PT:\t")
	} else if rn == 9 {
		fmt.Printf("CT:\t")
	} else {
		fmt.Printf("R%d:\t", rn)
	}
	for _, v := range rnd {
		fmt.Printf("%d\t", v)
	}
	fmt.Printf("\n")
}

func imul(a, b uint16) uint16 {
	var p int64
	var q uint64
	if a == 0 {
		p = int64(65537 - int64(b))
	} else if b == 0 {
		p = int64(65537 - int64(a))
	} else {
		q = uint64(a) * uint64(b)
		p = int64((q & 65535) - (q >> 16))
		if p <= 0 {
			p = p + 65537
		}
	}
	return uint16(p & 65535)
}

func inv(a uint16) uint16 {
	var n1, n2, q, r, b1, b2, t int64
	if a == 0 {
		b2 = 0
	} else {
		n1 = 65537
		n2 = int64(a)
		b2 = 1
		b1 = 0
		for r = 1; r != 0; {
			r = n1 % n2
			q = (n1 - r) / n2
			if r == 0 {
				if b2 < 0 {
					b2 = 65537 + b2
				}
			} else {
				n1 = n2
				n2 = r
				t = b2
				b2 = b1 - (q * b2)
				b1 = t
			}
		}
	}
	return uint16(b2)

}

func NewIdeaCipher(key []byte) IdeaCipher {
	if len(key) != KeySize {
		panic("Invalid key size slice")
	}
	var b IdeaCipher
	b.TestMode = false
	b.key = make([]byte, KeySize)
	b.key = key
	b.encSubBlocks = make([]uint16, subBlockQty)
	j := 0
	for i := 0; i < subBlockQty; i++ {
		b.encSubBlocks[i] = binary.BigEndian.Uint16(b.key[j:])
		j += 2
		if j >= KeySize {
			for k := 0; k < keyShifts; k++ {
				var endcarry uint8 = 0
				for l := 0; l < KeySize; l++ {
					if l == 0 {
						endcarry = b.key[0] & 0x80
						endcarry >>= 7
					} else {
						thiscarry := b.key[l] & 0x80
						thiscarry >>= 7
						b.key[l-1] |= thiscarry
					}
					b.key[l] <<= 1
				}
				if endcarry == 1 {
					b.key[15] |= 1
				}
			}
			j = 0
		}
	}
	//buf := make([]uint16, subBlockQty)
	b.decSubBlocks = make([]uint16, subBlockQty)
	b.decSubBlocks[0] = inv(b.encSubBlocks[48])
	b.decSubBlocks[1] = inv(b.encSubBlocks[49])
	b.decSubBlocks[2] = uint16(65536 - int(b.encSubBlocks[50]))
	b.decSubBlocks[3] = uint16(65536 - int(b.encSubBlocks[51]))
	b.decSubBlocks[4] = b.encSubBlocks[46]
	b.decSubBlocks[5] = b.encSubBlocks[47]

	b.decSubBlocks[6] = inv(b.encSubBlocks[42])
	b.decSubBlocks[7] = inv(b.encSubBlocks[43])
	b.decSubBlocks[8] = uint16(65536 - int(b.encSubBlocks[44]))
	b.decSubBlocks[9] = uint16(65536 - int(b.encSubBlocks[45]))
	b.decSubBlocks[10] = b.encSubBlocks[40]
	b.decSubBlocks[11] = b.encSubBlocks[41]

	b.decSubBlocks[12] = inv(b.encSubBlocks[36])
	b.decSubBlocks[13] = inv(b.encSubBlocks[37])
	b.decSubBlocks[14] = uint16(65536 - int(b.encSubBlocks[38]))
	b.decSubBlocks[15] = uint16(65536 - int(b.encSubBlocks[39]))
	b.decSubBlocks[16] = b.encSubBlocks[34]
	b.decSubBlocks[17] = b.encSubBlocks[35]

	b.decSubBlocks[18] = inv(b.encSubBlocks[30])
	b.decSubBlocks[19] = inv(b.encSubBlocks[31])
	b.decSubBlocks[20] = uint16(65536 - int(b.encSubBlocks[32]))
	b.decSubBlocks[21] = uint16(65536 - int(b.encSubBlocks[33]))
	b.decSubBlocks[22] = b.encSubBlocks[28]
	b.decSubBlocks[23] = b.encSubBlocks[29]

	b.decSubBlocks[24] = inv(b.encSubBlocks[24])
	b.decSubBlocks[25] = inv(b.encSubBlocks[25])
	b.decSubBlocks[26] = uint16(65536 - int(b.encSubBlocks[26]))
	b.decSubBlocks[27] = uint16(65536 - int(b.encSubBlocks[27]))
	b.decSubBlocks[28] = b.encSubBlocks[22]
	b.decSubBlocks[29] = b.encSubBlocks[23]

	b.decSubBlocks[30] = inv(b.encSubBlocks[18])
	b.decSubBlocks[31] = inv(b.encSubBlocks[19])
	b.decSubBlocks[32] = uint16(65536 - int(b.encSubBlocks[20]))
	b.decSubBlocks[33] = uint16(65536 - int(b.encSubBlocks[21]))
	b.decSubBlocks[34] = b.encSubBlocks[16]
	b.decSubBlocks[35] = b.encSubBlocks[17]

	b.decSubBlocks[36] = inv(b.encSubBlocks[12])
	b.decSubBlocks[37] = inv(b.encSubBlocks[13])
	b.decSubBlocks[38] = uint16(65536 - int(b.encSubBlocks[14]))
	b.decSubBlocks[39] = uint16(65536 - int(b.encSubBlocks[15]))
	b.decSubBlocks[40] = b.encSubBlocks[10]
	b.decSubBlocks[41] = b.encSubBlocks[11]

	b.decSubBlocks[42] = inv(b.encSubBlocks[6])
	b.decSubBlocks[43] = inv(b.encSubBlocks[7])
	b.decSubBlocks[44] = uint16(65536 - int(b.encSubBlocks[8]))
	b.decSubBlocks[45] = uint16(65536 - int(b.encSubBlocks[9]))
	b.decSubBlocks[46] = b.encSubBlocks[4]
	b.decSubBlocks[47] = b.encSubBlocks[5]

	b.decSubBlocks[48] = inv(b.encSubBlocks[0])
	b.decSubBlocks[49] = inv(b.encSubBlocks[1])
	b.decSubBlocks[50] = uint16(65536 - int(b.encSubBlocks[2]))
	b.decSubBlocks[51] = uint16(65536 - int(b.encSubBlocks[3]))

	return b
}

func (t IdeaCipher) GetSubKeys() ([]uint16, []uint16) {
	return t.encSubBlocks, t.decSubBlocks
}

func (IdeaCipher) BlockSize() int {
	return BlockSize
}

func (t IdeaCipher) Encrypt(dst, src []byte) {
	if len(src) != BlockSize || len(dst) != BlockSize {
		panic("Invalid block size passed to cipher.")
	}
	pt := make([]uint16, BlockSize/2)
	for i, _ := range pt {
		pt[i] = binary.BigEndian.Uint16(src[i*2:])
	}
	t.rounds(pt, t.encSubBlocks)
	for i, v := range pt {
		binary.BigEndian.PutUint16(dst[i*2:], v)
	}
}

func (t IdeaCipher) Decrypt(dst, src []byte) {
	if len(src) != BlockSize || len(dst) != BlockSize {
		panic("Invalid block size passed to cipher.")
	}
	pt := make([]uint16, BlockSize/2)
	for i, _ := range pt {
		pt[i] = binary.BigEndian.Uint16(src[i*2:])
	}
	t.rounds(pt, t.decSubBlocks)
	for i, v := range pt {
		binary.BigEndian.PutUint16(dst[i*2:], v)
	}
}

func (t IdeaCipher) rounds(x, keySched []uint16) {
	if len(keySched) != subBlockQty {
		panic("Wrong length of sub key blocks passed")
	}
	for i := 0; i < 8; i++ {
		if t.TestMode == true {
			printRound(i, x)
		}
		x[0] = imul(x[0], keySched[(6*i)])
		x[1] = imul(x[1], keySched[(6*i)+1])
		x[2] = x[2] + keySched[(6*i)+2]
		x[3] = x[3] + keySched[(6*i)+3]
		q := x[0] ^ x[2]
		r := x[1] ^ x[3]
		q = imul(q, keySched[(6*i)+4])
		s := q + r
		s = imul(s, keySched[(6*i)+5])
		t := q + s
		u := x[0] ^ s
		v := x[2] ^ s
		w := x[1] ^ t
		y := x[3] ^ t
		x[0] = v
		x[1] = y
		x[2] = u
		x[3] = w
	}
	if t.TestMode == true {
		printRound(8, x)
	}
	x[0] = imul(x[0], keySched[48+0])
	x[1] = imul(x[1], keySched[48+1])
	x[2] = x[2] + keySched[48+2]
	x[3] = x[3] + keySched[48+3]
	if t.TestMode == true {
		printRound(9, x)
	}
}
