package utils

import (
	"github.com/DSiSc/craft/types"
	"math/big"
)

// New a transaction
func newTransaction(nonce uint64, to *types.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int,
	data []byte, from *types.Address) *types.Transaction {

	if len(data) > 0 {
		data = CopyBytes(data)
	}
	d := types.TxData{
		AccountNonce: nonce,
		Recipient:    to,
		From:         from,
		Payload:      data,
		Amount:       new(big.Int),
		GasLimit:     gasLimit,
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
	if amount != nil {
		d.Amount.Set(amount)
	}
	if gasPrice != nil {
		d.Price.Set(gasPrice)
	}

	return &types.Transaction{Data: d}
}

func NewTransaction(nonce uint64, to *types.Address, amount *big.Int, gasLimit uint64,
	gasPrice *big.Int, data []byte, from types.Address) *types.Transaction {

	return newTransaction(nonce, to, amount, gasLimit, gasPrice, data, &from)
}
