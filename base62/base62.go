package base62

var (
	base62 = []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd',
		'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y', 'z'}
	uint8Index = []uint64{
		0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2,
		3, 4, 5, 6, 7, 8, 9, 0, 0, 0,
		0, 0, 0, 0, 10, 11, 12, 13, 14,
		15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
		25, 26, 27, 28, 29, 30, 31, 32, 33, 34,
		35, 0, 0, 0, 0, 0, 0, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
		50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
		60, 61, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, // 256
	}
	pow62Index = []uint64{
		1, 62, 3844, 238328, 14776336, 916132832, 56800235584,
		3521614606208, 218340105584896, 13537086546263552,
		839299365868340224, 9223372036854775808,
	}
)

// Encode encodes a number to base62.
func Encode(value uint64) string {
	var res [14]byte
	var i int
	for i = len(res) - 1; ; i-- {
		res[i] = base62[value%62]
		value /= 62
		if value == 0 {
			break
		}
	}

	return string(res[i:])
}

// EncodeByte encodes a number to base62 byte.
func EncodeByte(value uint64) []byte {
	var res [14]byte
	var i int
	for i = len(res) - 1; ; i-- {
		res[i] = base62[value%62]
		value /= 62
		if value == 0 {
			break
		}
	}

	return res[i:]
}

// Decode decodes a base36-encoded string.
func Decode(s string) uint64 {
	if len(s) > 12 {
		s = s[:11]
	}
	res := uint64(0)
	l := len(s) - 1
	for idx := 0; idx < len(s); idx++ {
		c := s[l-idx]
		res += uint8Index[c] * pow62Index[idx]
	}
	return res
}
