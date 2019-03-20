package contracts

import "github.com/DSiSc/craft/types"

type WhiteList interface {

	// issue a proposal to add whitelist
	IssueWhileListProposal(uint64, types.Address, uint64)

	// vote for white list proposal
	VoteForWhiteListProposal(uint64)

	// issue a proposal to update contract address
	IssueContractProposal(uint64, types.ContractType)

	// voteForContractProposal
	VoteForContractProposal(uint64)
}
