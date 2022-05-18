package types

import "errors"

var (
	McCallResIsNullErr = errors.New("mc result is null")
	UnSupportErr       = errors.New("unSupport chain")
	ParamsErr          = errors.New("params error")
	ReadDataOnChainErr = errors.New("pull tokenURI or totalSupply error for nft")
)
