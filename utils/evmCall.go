package utils

import (
	"github.com/DSiSc/apigateway/core/types"
	"github.com/DSiSc/blockchain"
	typec "github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG"
	"github.com/DSiSc/validator/worker"
	"github.com/DSiSc/validator/worker/common"
)

func EvmCall(tx *typec.Transaction, blockNr types.BlockNumber) ([]byte, uint64, bool, error) {
	bc, err := blockchain.NewLatestStateBlockChain()
	if err != nil {
		return nil, 0, true, err
	}
	var block *typec.Block
	if blockNr == types.LatestBlockNumber {
		block = bc.GetCurrentBlock()
	} else {
		height := blockNr.Touint64()
		block, err = bc.GetBlockByHeight(height)
		if err != nil {
			return nil, 0, true, err
		}
	}

	bchash, err := blockchain.NewBlockChainByBlockHash(block.HeaderHash)
	if err != nil {
		return nil, 0, true, err
	}

	context := evm.NewEVMContext(*tx, block.Header, bchash, block.Header.CoinBase)
	evmEnv := evm.NewEVM(context, bchash)
	gp := new(common.GasPool).AddGas(uint64(65536))
	result, gas, failed, err, _ := worker.ApplyTransaction(evmEnv, tx, gp)
	return result, gas, failed, err
}
