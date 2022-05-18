package chains

import "github.com/ThreeAndTwo/sync-nft-traits/types"

type IChain interface {
	GetNftInfo() (*types.NftInfo, error)
}
