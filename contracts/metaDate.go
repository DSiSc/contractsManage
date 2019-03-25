package contracts

import (
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/contractsManage/utils"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	"math/big"
)

const (
	GetContractById  = "8ba9ac6f"
	ContractMetaDate = "1e7936bd"
	RegisteContract  = "f38bbf6d"
)

type MetaData interface {

	// get contract address by contract id
	GetContractById(types.ContractType) types.Address

	// issue a proposal to update contract address
	UpdateContractAddress(uint64, types.ContractType, types.Address)

	// voteForContractProposal
	UpdateWhiteListAddress(uint64, types.Address)
}

type MetaDataContract struct {
	Address types.Address
	MetaMap map[types.ContractType]types.Address
}

func NewMetaDataContract() MetaData {
	return &MetaDataContract{
		Address: utils.HexToAddress(types.MetaDataContractAddress),
		MetaMap: make(map[types.ContractType]types.Address),
	}
}

func (instance *MetaDataContract) GetContractById(contractType types.ContractType) types.Address {
	chain, err := blockchain.NewLatestStateBlockChain()
	if err != nil {
		panic(fmt.Errorf("failed to create init-state block chain, as: %v", err))
	}
	block := chain.GetCurrentBlock()
	nonce := utils.GetAccountNonce(block, types.Address{})
	callCode := fmt.Sprintf("%s%s", GetContractById, utils.EncodeUint256(int(contractType)))
	result, err := utils.EvmCall(nonce+1, &instance.Address, big.NewInt(0),
		uint64(0), big.NewInt(0), utils.Hex2Bytes(callCode), nil)
	if nil != err {
		log.Error("evm called to get contract by id error with %v.", err)
		return utils.HexToAddress(types.JustiitaContractDefaultAddress)
	}
	address := utils.BytesToAddress(result[12:])
	return address
}

func (instance *MetaDataContract) UpdateContractAddress(uint64, types.ContractType, types.Address) {

}

func (instance *MetaDataContract) UpdateWhiteListAddress(uint64, types.Address) {

}
