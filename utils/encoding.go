package utils

import (
	"fmt"
	"strconv"
)

func getHexDigit(a int) (hex string) {
	if a < 10 {
		return strconv.Itoa(a)
	}

	return string([]byte{'a' + byte(a-10)})
}

func toHex(num int) string {
	if 0 == num {
		return "0"
	} else if num < 0 {
		num = 4294967296 + num
	}

	var ret string
	for num > 0 {
		a := num % 16
		num = num / 16

		hexDigit := getHexDigit(a)
		ret = hexDigit + ret
	}

	return ret
}

func EncodeUint256(num int) string {
	var defaultUint256 = "0000000000000000000000000000000000000000000000000000000000000000"
	hexStr := toHex(num)
	encoding := fmt.Sprintf("%s%s", defaultUint256[:len(defaultUint256)-len(hexStr)], hexStr)
	return encoding
}

func EncodingAddress(address string) string {
	var defaultAddress = "0000000000000000000000000000000000000000000000000000000000000000"
	encoding := fmt.Sprintf("%s%s", defaultAddress[:len(defaultAddress)-len(address)], address)
	return encoding
}
