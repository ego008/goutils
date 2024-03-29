package goutils

import (
	"math/rand"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetIdxBits = 6                     // 6 bits to represent a charset index
	charsetIdxMask = 1<<charsetIdxBits - 1 // All 1-bits, as many as charsetIdxBits
	charsetIdxMax  = 63 / charsetIdxBits   // # of letter indices fitting in 63 bits
)

// ExtendByteSlice extends b to needLen bytes
func ExtendByteSlice(b []byte, needLen int) []byte {
	b = b[:cap(b)]
	if n := needLen - cap(b); n > 0 {
		b = append(b, make([]byte, n)...)
	}
	return b[:needLen]
}

// RandBytes returns dst with a string random bytes
// Make sure that dst has the length you need.
func RandBytes(dst []byte) []byte {
	n := len(dst)

	for i, cache, remain := n-1, src.Int63(), charsetIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), charsetIdxMax
		}
		if idx := int(cache & charsetIdxMask); idx < len(charset) {
			dst[i] = charset[idx]
			i--
		}
		cache >>= charsetIdxBits
		remain--
	}

	return dst
}

// Bconcat concat a list of byte
func Bconcat(slices [][]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

// Bconcat2 concat a list of byte
func Bconcat2(slices ...[]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func CloneBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
