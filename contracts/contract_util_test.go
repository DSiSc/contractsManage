package contracts

import (
	"encoding/hex"
	"fmt"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbiEncode(t *testing.T) {
	assert := assert.New(t)
	addrList := []types.Address{
		util.HexToAddress("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b"),
		util.HexToAddress("0xE5fb3dAe9382c39f307e8b0d4C659eD5eC53a725"),
	}
	valBytes, err := AbiEncode(addrList)
	assert.Nil(err)
	assert.Equal("0000000000000000000000000000000000000000000000000000000000000020"+
		"0000000000000000000000000000000000000000000000000000000000000002"+
		"000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"+
		"000000000000000000000000e5fb3dae9382c39f307e8b0d4c659ed5ec53a725", fmt.Sprintf("%x", valBytes))
}

func TestAbiEncode2(t *testing.T) {
	assert := assert.New(t)
	addr := util.HexToAddress("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	valBytes, err := AbiEncode(addr)
	assert.Nil(err)
	assert.Equal("000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", fmt.Sprintf("%x", valBytes))
}

func TestAbiEncode3(t *testing.T) {
	assert := assert.New(t)
	addr := util.HexToAddress("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	url := "tcp://192.168.1.1:8080"
	valBytes, err := AbiEncode(addr, url)
	assert.Nil(err)
	assert.Equal("000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000167463703a2f2f3139322e3136382e312e313a3830383000000000000000000000", fmt.Sprintf("%x", valBytes))
}

func TestAbiEncode4(t *testing.T) {
	assert := assert.New(t)
	idle, commit, viewChange := uint64(2000), uint64(2000), uint64(5000)
	valBytes, err := AbiEncode(idle, commit, viewChange)
	assert.Nil(err)
	assert.Equal("00000000000000000000000000000000000000000000000000000000000007d000000000000000000000000000000000000000000000000000000000000007d00000000000000000000000000000000000000000000000000000000000001388", fmt.Sprintf("%x", valBytes))
}

func TestAbiDecode(t *testing.T) {
	assert := assert.New(t)

	input, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000020" +
		"0000000000000000000000000000000000000000000000000000000000000003" +
		"000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b" +
		"000000000000000000000000e5fb3dae9382c39f307e8b0d4c659ed5ec53a725" +
		"00000000000000000000000059b3f85ba6eb737fd0fad93bc4b5f92fd8c591de")
	addrList := []types.Address{}
	err := AbiDecode(input, &addrList)
	assert.Nil(err)
	assert.Equal(3, len(addrList))
	assert.Equal(util.HexToAddress("a94f5374fce5edbc8e2a8697c15331677e6ebf0b"), addrList[0])
	assert.Equal(util.HexToAddress("e5fb3dae9382c39f307e8b0d4c659ed5ec53a725"), addrList[1])
	assert.Equal(util.HexToAddress("59b3f85ba6eb737fd0fad93bc4b5f92fd8c591de"), addrList[2])
}

func TestAbiDecode2(t *testing.T) {
	assert := assert.New(t)
	input, _ := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000167463703a2f2f3139322e3136382e312e313a3830383000000000000000000000")
	var url string
	err := AbiDecode(input, &url)
	assert.Nil(err)
	assert.NotEqual(0, len(url))
	assert.Equal("tcp://192.168.1.1:8080", url)
}

func TestAbiDecode3(t *testing.T) {
	assert := assert.New(t)
	var idle, commit, viewChange uint64
	input, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000007d000000000000000000000000000000000000000000000000000000000000007d00000000000000000000000000000000000000000000000000000000000001388")
	err := AbiDecode(input, &idle, &commit, &viewChange)
	assert.Nil(err)
	assert.Equal(uint64(2000), idle)
	assert.Equal(uint64(2000), commit)
	assert.Equal(uint64(5000), viewChange)
}

func TestAbiDecode4(t *testing.T) {
	assert := assert.New(t)
	var val bool
	input, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	err := AbiDecode(input, &val)
	assert.Nil(err)
	assert.Equal(false, val)
	input, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	err = AbiDecode(input, &val)
	assert.Nil(err)
	assert.Equal(true, val)
}
