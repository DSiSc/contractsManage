package utils

import (
	"encoding/hex"
	"github.com/DSiSc/craft/types"
	"math/big"
)

func SetBytes(b []byte, a *types.Address) {
	if len(b) > len(a) {
		b = b[len(b)-types.AddressLength:]
	}
	copy(a[types.AddressLength-len(b):], b)
}

func BytesToAddress(b []byte) types.Address {
	var a types.Address
	SetBytes(b, &a)
	return a
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

type RefAddress struct {
	Addr types.Address
}

func NewRefAddress(addr types.Address) *RefAddress {
	return &RefAddress{Addr: addr}
}

func (ref *RefAddress) Address() types.Address {
	return ref.Addr
}

func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

func HexToAddress(s string) types.Address {
	return BytesToAddress(FromHex(s))
}

// New a transaction
func NewTransactionForCall(to types.Address, data []byte) types.Transaction {
	if len(data) > 0 {
		data = CopyBytes(data)
	}
	from := HexToAddress("0x0000000000000000000000000000000000000000")
	d := types.TxData{
		Recipient: &to,
		From:      &from,
		Payload:   data,
		Amount:    new(big.Int),
		Price:     new(big.Int),
		V:         new(big.Int),
		R:         new(big.Int),
		S:         new(big.Int),
	}

	return types.Transaction{Data: d}
}

func BigEndianToUin64(b []byte) uint64 {
	_ = b[31] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[31]) | uint64(b[30])<<8 | uint64(b[29])<<16 | uint64(b[28])<<24 |
		uint64(b[27])<<32 | uint64(b[26])<<40 | uint64(b[25])<<48 | uint64(b[24])<<56
}

func Bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}
