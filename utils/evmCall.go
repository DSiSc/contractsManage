package utils

import (
	typec "github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG"
	"github.com/DSiSc/repository"
	"math"
	"math/big"
)

type RefAddress struct {
	Addr typec.Address
}

func NewRefAddress(addr typec.Address) *RefAddress {
	return &RefAddress{Addr: addr}
}

func (self *RefAddress) Address() typec.Address {
	return self.Addr
}

func EvmCall(nonce uint64, to *typec.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int,
	data []byte, from *typec.Address) ([]byte, error) {

	tx := newTransaction(nonce, to, amount, gasLimit, gasPrice, data, &typec.Address{})
	bc, err := repository.NewLatestStateRepository()
	if err != nil {
		return nil, err
	}
	block := bc.GetCurrentBlock()
	bchash, err := repository.NewRepositoryByBlockHash(block.HeaderHash)
	if err != nil {
		return nil, err
	}
	context := evm.NewEVMContext(*tx, block.Header, bchash, block.Header.CoinBase)
	evmEnv := evm.NewEVM(context, bchash)
	sender := NewRefAddress(*tx.Data.From)
	result, _, err := evmEnv.Call(sender, *tx.Data.Recipient, tx.Data.Payload, math.MaxUint64, tx.Data.Amount)
	return result, err
}
