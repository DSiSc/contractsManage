package contracts

import (
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/contractsManage/utils"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	"math/big"
	"sync"
	"regexp"
)

// byte code which match to function name
const (
	totalNodes            = "9592d424"
	GetCandidataByRanking = "ae3364a4"
	UrlReg = "[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}:[0-9]{2,5}"
)

type Voting interface {

	// get node num
	NodeNumber() uint64

	// get nodes list which sorted
	GetNodeList(count uint64) ([]NodeInfo, error)
}

type NodeInfo struct {
	Address types.Address `json:"address"     gencodec:"required"`
	Url     string        `json:"url"`
	Id      uint64        `json:"id"`
}

type VotingContract struct {
	mutex           sync.RWMutex
	contractAddress types.Address
	// nodeNumber will update when called NodeNumber()
	nodeNumber uint64
}

func NewVotingContract() Voting {
	metaData := NewMetaDataContract()
	voteAddress := metaData.GetContractById(types.VoteContractType)
	contract := &VotingContract{
		contractAddress: voteAddress,
	}
	return contract
}

func (vote *VotingContract) NodeNumber() uint64 {
	chain, err := blockchain.NewLatestStateBlockChain()
	if err != nil {
		panic(fmt.Errorf("failed to create init-state block chain, as: %v", err))
	}
	block := chain.GetCurrentBlock()
	callCode := utils.Hex2Bytes(totalNodes)
	nonce := utils.GetAccountNonce(block, types.Address{})
	result, err := utils.EvmCall(nonce, &vote.contractAddress, big.NewInt(0),
		uint64(0), big.NewInt(0), callCode, nil)
	if nil != err {
		log.Error("error")
		return types.MinimunNodesForDpos
	}
	nodeNum := utils.BigEndianToUin64(result)
	vote.mutex.Lock()
	defer vote.mutex.Unlock()
	vote.nodeNumber = nodeNum
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
	for index := 0; uint64(index) < count; index++ {
		callCode := fmt.Sprintf("%s%s", GetCandidataByRanking, utils.EncodeUint256(index))
		nonce := utils.GetAccountNonce(block, types.Address{})
		result, err := utils.EvmCall(nonce, &vote.contractAddress, big.NewInt(0),
			uint64(0), big.NewInt(0), utils.Hex2Bytes(callCode), nil)
		if nil != err {
			log.Error("error")
			return NodeList, fmt.Errorf("get node info failed with error %v", err)
		}
		nodeInfo := decodeNodeResult(result)
		NodeList = append(NodeList, nodeInfo)
	}
	return NodeList, nil
}

func decodeNodeResult(out []byte) NodeInfo {
	bigOffsetEnd := big.NewInt(0).SetBytes(out[96:128])
	offsetEnd := 128 + int(bigOffsetEnd.Uint64())
	rawUrl := string(out[128:offsetEnd])

	return NodeInfo{
		Address: utils.BytesToAddress(out[12:32]),
		Id:      utils.BigEndianToUin64(out[32:64]),
		Url:     string(rawUrl),
	}
}

func regURL(rawUrl string) (url string){
	reg := regexp.MustCompile(UrlReg)
	data := reg.Find([]byte(rawUrl))
	return string(data)
}