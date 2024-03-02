package goutils

import "strconv"

func SliceUniqInt(s []int) []int {
	if len(s) == 0 {
		return s
	}
	seen := make(map[int]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func SliceUniqStr(s []string) []string {
	if len(s) == 0 {
		return s
	}
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func FormatInt(number int64) string {
	str := strconv.FormatInt(number, 10)
	lStr := len(str)
	digits := lStr
	if number < 0 {
		digits--
	}
	commas := (digits+2)/3 - 1
	lBuf := lStr + commas
	var sBuf [32]byte // pre allocate buffer at stack rather than make([]byte,n)
	buf := sBuf[0:lBuf]
	// copy str from the end
	for si, bi, c3 := lStr-1, lBuf-1, 0; ; {
		buf[bi] = str[si]
		if si == 0 {
			return string(buf)
		}
		si--
		bi--
		// insert comma every 3 chars
		c3++
		if c3 == 3 && (si > 0 || number > 0) {
			buf[bi] = ','
			bi--
			c3 = 0
		}
	}
}
