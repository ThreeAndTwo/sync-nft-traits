package chains

import (
	"github.com/ThreeAndTwo/sync-nft-traits/chains/ethereum"
	"github.com/ThreeAndTwo/sync-nft-traits/types"
)

func NewChains(chainId int64, rpc, contract, mcContract string) (IChain, error) {
	switch chainId {
	case 1, 56, 66, 128, 250:
		return ethereum.NewEvm(rpc, contract, mcContract)
	default:
		return nil, types.UnSupportErr
	}
}
