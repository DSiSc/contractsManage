package contracts

import (
	"fmt"
	"github.com/DSiSc/contractsManage/utils"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG/system/contract/util"
	"github.com/DSiSc/repository"
	"math"
	"math/big"
)

// DPOS Contract utils
var ContractMsp = &msp{}
var (
	initManagerMethodHash         = util.ExtractMethodHash(util.Hash([]byte("initManagers(address[])")))
	voteAddManagerMethodHash      = util.ExtractMethodHash(util.Hash([]byte("voteAddManager(address)")))
	voteRemoveManagerMethodHash   = util.ExtractMethodHash(util.Hash([]byte("voteRemoveManager(address)")))
	authorizeMemberMethodHash     = util.ExtractMethodHash(util.Hash([]byte("authorizeMember(address)")))
	revokeAuthorizationMethodHash = util.ExtractMethodHash(util.Hash([]byte("revokeAuthorization(address)")))
	isAuthorizedMethodHash        = util.ExtractMethodHash(util.Hash([]byte("isAuthorized(address)")))
)

// dposPbft dpos contract utils
type msp struct {
}

// MspInitTransaction build init manager list transaction
func (self *msp) MspInitTransaction(nonce uint64, caller, contractAddr *types.Address, members []ConsensusMember) (*types.Transaction, error) {
	managers := make([]types.Address, 0)
	for _, member := range members {
		managers = append(managers, member.Addr)
	}
	// init managers
	if tx, err := self.initManagers(nonce, caller, contractAddr, managers); err != nil {
		return nil, err
	} else {
		return tx, nil
	}
}

func (self *msp) initManagers(nonce uint64, caller, contractAddr *types.Address, managers []types.Address) (*types.Transaction, error) {
	input := initManagerMethodHash[:4]
	args, err := AbiEncode(managers)
	if err != nil {
		return nil, fmt.Errorf("failed to encode init member contract call params, as: %v", err)
	}
	input = append(input, args...)

	d := types.TxData{
		AccountNonce: nonce,
		Recipient:    contractAddr,
		From:         caller,
		Payload:      input,
		Amount:       new(big.Int),
		GasLimit:     math.MaxUint64,
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
	return &types.Transaction{Data: d}, nil
}

func (self *msp) IsAuthorized(addr *types.Address) bool {
	mspContractAddr, _ := ContractsAddr.Load(types.MspContract)
	if mspContractAddr == nil {
		return true
	}
	input := isAuthorizedMethodHash[:4]
	valBytes, err := AbiEncode(*addr)
	input = append(input, valBytes...)
	if err != nil {
		return false
	}
	chain, err := repository.NewLatestStateRepository()
	if err != nil {
		panic(fmt.Errorf("failed to create init-state block chain, as: %v", err))
	}
	block := chain.GetCurrentBlock()
	nonce := utils.GetAccountNonce(block, types.Address{})
	result, err := utils.EvmCall(nonce, mspContractAddr.(*types.Address), big.NewInt(0),
		uint64(0), big.NewInt(0), input, nil)
	if nil != err {
		return false
	}
	var isAuthorized bool
	if err := AbiDecode(result, &isAuthorized); err == nil {
		return isAuthorized
	}
	return false
}
