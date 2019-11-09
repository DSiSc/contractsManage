package contracts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DSiSc/contractsManage/utils"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG/system/contract/util"
	"github.com/DSiSc/justitia-pbft/common"
	"github.com/DSiSc/repository"
	"math"
	"math/big"
)

// DPOS Contract utils
var DposContractUtils = &dposPbft{}

var (
	emptyAddr                   = types.Address{}
	initMemberMethodHash        = util.ExtractMethodHash(util.Hash([]byte("initMembers(address[])")))
	initTimerMethodHash         = util.ExtractMethodHash(util.Hash([]byte("initConsensusTimer(uint64,uint64,uint64,uint64)")))
	initMemberUrlMethodHash     = util.ExtractMethodHash(util.Hash([]byte("initMemberUrl(address,string)")))
	getMembersMethodHash        = util.ExtractMethodHash(util.Hash([]byte("getMembers()")))
	getMemberUrlMethodHash      = util.ExtractMethodHash(util.Hash([]byte("getMemberUrl(address)")))
	getConsensusTimerMethodHash = util.ExtractMethodHash(util.Hash([]byte("getConsensusTimer()")))
)

// ConsensusMember dpos consensus member info
type ConsensusMember struct {
	Addr types.Address `json:"addr"`
	// url used to communicate
	URL string `json:"url"`
}

type consensusMemberJson struct {
	Addr string `json:"addr"`
	// url used to communicate
	URL string `json:"url"`
}

func (self *ConsensusMember) MarshalJSON() ([]byte, error) {
	replica := consensusMemberJson{
		Addr: common.AddressToHex(self.Addr),
		URL:  self.URL,
	}
	return json.Marshal(replica)
}

func (self *ConsensusMember) UnmarshalJSON(data []byte) error {
	replica := new(consensusMemberJson)
	if err := json.Unmarshal(data, replica); err != nil {
		return err
	}
	self.Addr = common.HexToAddress(replica.Addr)
	self.URL = replica.URL
	return nil
}

// ConsensusConf dpos consensus config
type ConsensusConf struct {
	BlockDelay        uint64 `json:"blockDelay"`
	IdleTimeOut       uint64 `json:"idleTimeOut"`
	CommitTimeOut     uint64 `json:"commitTimeOut"`
	ViewChangeTimeOut uint64 `json:"viewChangeTimeOut"`
}

// dposPbft dpos contract utils
type dposPbft struct {
}

// InitMemberTransaction build init member list transaction
func (self *dposPbft) DposInitTransactions(nonce uint64, caller, contractAddr *types.Address, members []ConsensusMember, conf ConsensusConf) ([]*types.Transaction, error) {
	memberAddrs := make([]types.Address, 0)
	for _, member := range members {
		memberAddrs = append(memberAddrs, member.Addr)
	}

	txs := make([]*types.Transaction, 0)

	var tx *types.Transaction
	var err error

	// init members
	if tx, err = self.initMembers(nonce, caller, contractAddr, memberAddrs); err != nil {
		return nil, err
	}
	txs = append(txs, tx)
	nonce++

	// init member url
	for _, member := range members {
		if tx, err = self.initMemberUrl(nonce, caller, contractAddr, member); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
		nonce++
	}

	// init consensus timer
	if tx, err = self.initConsensusTimer(nonce, caller, contractAddr, conf); err != nil {
		return nil, err
	}
	txs = append(txs, tx)
	nonce++

	return txs, nil
}

func (self *dposPbft) initMembers(nonce uint64, caller, contractAddr *types.Address, members []types.Address) (*types.Transaction, error) {
	input := initMemberMethodHash[:4]
	args, err := AbiEncode(members)
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

func (self *dposPbft) initMemberUrl(nonce uint64, caller, contractAddr *types.Address, member ConsensusMember) (*types.Transaction, error) {
	input := initMemberUrlMethodHash[:4]
	args, err := AbiEncode(member.Addr, member.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to encode init member url contract call params, as: %v", err)
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

func (self *dposPbft) initConsensusTimer(nonce uint64, caller, contractAddr *types.Address, conf ConsensusConf) (*types.Transaction, error) {
	input := initTimerMethodHash[:4]
	args, err := AbiEncode(conf.BlockDelay, conf.IdleTimeOut, conf.CommitTimeOut, conf.ViewChangeTimeOut)
	if err != nil {
		return nil, fmt.Errorf("failed to encode init consensus timer conf params, as: %v", err)
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

// GetMembers get member list
func (self *dposPbft) GetMembers(hash types.Hash) []types.Address {
	input := getMembersMethodHash[:4]
	chain, err := repository.NewRepositoryByBlockHash(hash)
	if err != nil {
		panic(fmt.Errorf("failed to create init-state block chain, as: %v", err))
	}
	block := chain.GetCurrentBlock()
	nonce := utils.GetAccountNonce(block, types.Address{})
	dposContractAddr, _ := ContractsAddr.Load(types.DposBftVotingContract)
	result, err := utils.EvmCall(nonce, dposContractAddr.(*types.Address), big.NewInt(0),
		uint64(0), big.NewInt(0), input, nil)
	if nil != err {
		return nil
	}
	addrList := []types.Address{}
	if err := AbiDecode(result, &addrList); err != nil {
		return nil
	}
	validAddrList := make([]types.Address, 0)
	for _, addr := range addrList {
		if !bytes.Equal(addr[:], emptyAddr[:]) {
			validAddrList = append(validAddrList, addr)
		}
	}
	return validAddrList
}

// GetMemberUrl get member url
func (self *dposPbft) GetMemberUrl(hash types.Hash, address types.Address) string {
	input := getMemberUrlMethodHash[:4]
	valBytes, err := AbiEncode(address)
	input = append(input, valBytes...)
	if err != nil {
		return ""
	}
	chain, err := repository.NewRepositoryByBlockHash(hash)
	if err != nil {
		panic(fmt.Errorf("failed to create init-state block chain, as: %v", err))
	}
	block := chain.GetCurrentBlock()
	nonce := utils.GetAccountNonce(block, types.Address{})
	dposContractAddr, _ := ContractsAddr.Load(types.DposBftVotingContract)
	result, err := utils.EvmCall(nonce, dposContractAddr.(*types.Address), big.NewInt(0),
		uint64(0), big.NewInt(0), input, nil)
	if nil != err {
		return ""
	}
	var url string
	if err := AbiDecode(result, &url); err != nil {
		return ""
	}
	return url
}

// GetMembersAndFaultTolerants get member list and fault tolerants.
func (self *dposPbft) GetMembersAndFaultTolerants(hash types.Hash) ([]types.Address, int) {
	members := self.GetMembers(hash)
	return members, (len(members) - 1) / 3
}

// GetMembersAndFaultTolerants get member list and fault tolerants.
func (self *dposPbft) GetConsensusTimer() (*ConsensusConf, error) {
	input := getConsensusTimerMethodHash[:4]
	chain, err := repository.NewLatestStateRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create init-state block chain, as: %v", err)
	}
	block := chain.GetCurrentBlock()
	nonce := utils.GetAccountNonce(block, types.Address{})
	dposContractAddr, _ := ContractsAddr.Load(types.DposBftVotingContract)
	result, err := utils.EvmCall(nonce, dposContractAddr.(*types.Address), big.NewInt(0),
		uint64(0), big.NewInt(0), input, nil)
	if nil != err {
		return nil, fmt.Errorf("failed to get consensus timer from contract, as: %v", err)
	}
	var delay, idle, commit, viewChange uint64
	if err := AbiDecode(result, &delay, &idle, &commit, &viewChange); err != nil {
		return nil, fmt.Errorf("failed to decode consensus time conf result, as: %v", err)
	}
	return &ConsensusConf{
		BlockDelay:        delay,
		IdleTimeOut:       idle,
		CommitTimeOut:     commit,
		ViewChangeTimeOut: viewChange,
	}, nil
}
