package ethereum

import (
	"time"

	"github.com/ThreeAndTwo/sync-nft-traits/types"
	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/ethrpc/provider/httprpc"
	"github.com/alethio/web3-multicall-go/multicall"
	"github.com/ethereum/go-ethereum/common"
)

type Evm struct {
	contract common.Address
	mcGetter *ChainMCallGetter
}

func NewEvm(rpc, contract, mcallContract string) (*Evm, error) {
	_contract := common.HexToAddress(contract)
	provider, err := httprpc.New(rpc)
	if err != nil {
		return nil, err
	}

	provider.SetHTTPTimeout(10 * time.Second)
	nodeClient, err := ethrpc.New(provider)
	if err != nil {
		return nil, err
	}

	m, err := multicall.New(nodeClient, multicall.ContractAddress(mcallContract))
	if err != nil {
		return nil, err
	}

	_chainInfo := &types.ChainInfo{
		McClient: &m,
	}
	mcGetter := newMultiCallGetter(_chainInfo)

	return &Evm{
		contract: _contract,
		mcGetter: mcGetter,
	}, nil
}

func (e *Evm) GetNftInfo() (*types.NftInfo, error) {
	info := e.mcGetter.GetNftInfoForContract(e.contract)
	if info == nil {
		return nil, types.McCallResIsNullErr
	}
	return info, nil
}
