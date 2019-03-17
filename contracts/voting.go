package contracts

import (
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/contractsManage/utils"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG"
	"math"
	"math/big"
	"sync"
)

type NodeInfo struct {
	Address types.Address `json:"address"     gencodec:"required"`
	Url     string        `json:"url"`
	Id      uint64        `json:"id"`
}

type Voting interface {

	// get node num
	NodeNumber() uint64

	// get nodes list which sorted
	GetNodeList(count uint64) ([]NodeInfo, error)
}

type VotingContract struct {
	mutex           sync.RWMutex
	contractAddress types.Address
	handlerFunc     map[string][]byte
	nodeNumber      uint64
}

func NewVotingContract() Voting {
	chain, err := blockchain.NewLatestStateBlockChain()
	if err != nil {
		panic(fmt.Errorf("failed to create init-state block chain, as: %v", err))
	}
	result, err := chain.Get([]byte(types.JustitiaVoting))
	contract := &VotingContract{
		contractAddress: utils.BytesToAddress(result),
	}
	contract.handleRegister()
	return contract
}

func (vote *VotingContract) handleRegister() {
	vote.handlerFunc = make(map[string][]byte)
	vote.handlerFunc["totalNodes"] = utils.Hex2Bytes("9592d424")
	vote.handlerFunc["candidateState"] = utils.Hex2Bytes("8ff49c88000000000000000000000000")
	vote.handlerFunc["GetCandidateByRanking"] = utils.Hex2Bytes("ae3364a4")
}

func (vote *VotingContract) NodeNumber() uint64 {
	chain, err := blockchain.NewLatestStateBlockChain()
	if err != nil {
		panic(fmt.Errorf("failed to create init-state block chain, as: %v", err))
	}
	block := chain.GetCurrentBlock()
	tx := utils.NewTransactionForCall(vote.contractAddress, vote.handlerFunc["totalNodes"])
	context := evm.NewEVMContext(tx, block.Header, chain, types.Address{})
	evmEnv := evm.NewEVM(context, chain)
	sender := utils.NewRefAddress(*tx.Data.From)
	result, _, err := evmEnv.Call(sender, *tx.Data.Recipient, tx.Data.Payload, math.MaxUint64, big.NewInt(0))
	if nil != err {
		log.Error("error")
		return uint64(4)
	}
	nodeNum := utils.BigEndianToUin64(result)
	vote.mutex.Lock()
	vote.nodeNumber = nodeNum
	vote.mutex.Unlock()
	return nodeNum
}

func (vote *VotingContract) GetNodeList(count uint64) ([]NodeInfo, error) {
	vote.mutex.RLock()
	defer vote.mutex.RUnlock()
	var NodeList = make([]NodeInfo, 0)
	if count > vote.nodeNumber {
		log.Error("invalid parameter count %d while node number is %d.", count, vote.nodeNumber)
		return make([]NodeInfo, 0), fmt.Errorf("invalid parameter count %d while node number is %d", count, vote.nodeNumber)
	}
	chain, err := blockchain.NewLatestStateBlockChain()
	if err != nil {
		panic(fmt.Errorf("failed to create init-state block chain, as: %v", err))
	}
	block := chain.GetCurrentBlock()
	base := "ae3364a4000000000000000000000000000000000000000000000000000000000000000"
	for index := 0; uint64(index) < count; index++ {
		payload := fmt.Sprintf("%s%d", base, index)
		tx := utils.NewTransactionForCall(vote.contractAddress, utils.Hex2Bytes(payload))
		context := evm.NewEVMContext(tx, block.Header, chain, types.Address{})
		evmEnv := evm.NewEVM(context, chain)
		sender := utils.NewRefAddress(*tx.Data.From)
		out, _, err := evmEnv.Call(sender, *tx.Data.Recipient, tx.Data.Payload, math.MaxUint64, big.NewInt(0))
		if nil != err {
			panic("error")
		}
		nodeinfo := decodeNodeResult(out)
		NodeList = append(NodeList, nodeinfo)
	}
	return NodeList, nil
}

func decodeNodeResult(out []byte) NodeInfo {
	return NodeInfo{
		Address: utils.BytesToAddress(out[32-types.AddressLength : 32]),
		Id:      utils.BigEndianToUin64(out[64:96]),
		Url:     string(out[127:143]),
	}
}
