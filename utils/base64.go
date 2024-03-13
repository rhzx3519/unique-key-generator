package utils

import (
	"math/rand"
	"slices"
)

func init() {
	for i := 0; i < 26; i++ {
		lookup = append(lookup, byte(i+'a'))
	}
	for i := 0; i < 10; i++ {
		lookup = append(lookup, byte(i+'0'))
	}
	for i := 0; i < 26; i++ {
		lookup = append(lookup, byte(i+'A'))
	}
	lookup = append(lookup, '.')
	lookup = append(lookup, '-')
}

var lookup []byte

// [a-z0-9A-Z.-]
func Base64Encode(i int64) string {
	if i == 0 {
		return "a"
	}
	var bytes []byte
	for ; i != 0; i /= 64 {
		bytes = append(bytes, lookup[i%64])
	}
	slices.Reverse(bytes)
	return string(bytes)
}

func Base64Decode(k string) int64 {
	var res int64
	if k == "" {
		return res
	}
	n := len(k)
	for i := 0; i < n; i++ {
		b := k[i]
		var j int
		if b >= 'a' && b <= 'z' {
			j = int(b - 'a')
		} else if b >= '0' && b <= '9' {
			j = int(b-'0') + 26
		} else if b >= 'A' && b <= 'Z' {
			j = int(b-'A') + 36
		} else if b == '.' {
			j = 62
		} else if b == '-' {
			j = 63
		}
		res = res*64 + int64(j)
	}
	return res
}

func base64DigitEncode(d int) byte {
	return lookup[d%64]
}

func base64DigitDecode(b byte) int {
	var j int
	if b >= 'a' && b <= 'z' {
		j = int(b - 'a')
	} else if b >= '0' && b <= '9' {
		j = int(b-'0') + 26
	} else if b >= 'A' && b <= 'Z' {
		j = int(b-'A') + 36
	} else if b == '.' {
		j = 62
	} else if b == '-' {
		j = 63
	}
	return j
}

func RandomBase64Str(length int) string {
	var tmp []byte
	for j := 0; j < length; j++ {
		tmp = append(tmp, base64DigitEncode(rand.Intn(64)))
	}
	return string(tmp)
}
