package contracts

import (
	"github.com/DSiSc/craft/types"
)

type MetaData interface {

	// register system contracts
	RegisterContract(types.ContractType, types.Address)

	// get contract address by contract id
	GetContractById(contractType types.ContractType) types.Address

	// issue a proposal to update contract address
	UpdateContractAddress(uint64, types.ContractType, types.Address)

	// voteForContractProposal
	UpdateWhiteListAddress(uint64, types.Address)
}
