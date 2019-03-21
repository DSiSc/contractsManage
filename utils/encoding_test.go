package utils

import (
	"github.com/DSiSc/craft/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeUint256(t *testing.T) {
	c := EncodeUint256(int(types.JustitiaRightContractType))
	assert.Equal(t, 64, len(c))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", c)

	c = EncodeUint256(int(types.VoteContractType))
	assert.Equal(t, 64, len(c))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000002", c)

	c = EncodeUint256(int(types.WhiteListContractType))
	assert.Equal(t, 64, len(c))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000003", c)

	c = EncodeUint256(int(types.MetaDataContractType))
	assert.Equal(t, 64, len(c))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000004", c)

}

func TestEncodingAddress(t *testing.T) {
	c := EncodingAddress(types.JustiitaContractDefaultAddress)
	assert.Equal(t, 64, len(c))
	assert.Equal(t, "000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb1", c)

	c = EncodingAddress(types.MetaDataContractAddress)
	assert.Equal(t, 64, len(c))
	assert.Equal(t, "0000000000000000000000008be503bcded90ed42eff31f56199399b2b0154ca", c)

	c = EncodingAddress(types.WhiteListContractTypeDefaultAddress)
	assert.Equal(t, 64, len(c))
	assert.Equal(t, "00000000000000000000000047e9fbef8c83a1714f1951f142132e6e90f5fa5d", c)

	c = EncodingAddress(types.VotingContractDefaultAddress)
	assert.Equal(t, 64, len(c))
	assert.Equal(t, "0000000000000000000000005a443704dd4b594b382c22a083e2bd3090a6fef3", c)
}
