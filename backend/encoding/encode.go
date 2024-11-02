package encoding

import (
	"crypto/md5"
	"encoding/hex"
)

const radix = 62

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(num int) string {
	res := ""
	for num > 0 {
		rem := num % radix
		res = string(alphabet[rem]) + res
		num /= radix
	}
	return res
}

func EncodeMd5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}
