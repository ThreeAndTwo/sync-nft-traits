package ethereum

import (
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/ThreeAndTwo/sync-nft-traits/types"
	"github.com/alethio/web3-multicall-go/multicall"
	"github.com/ethereum/go-ethereum/common"
)

type ChainMCallGetter struct {
	MChainInfo *types.ChainInfo
	McCounter  int
}

func newMultiCallGetter(_chainInfo *types.ChainInfo) *ChainMCallGetter {
	return &ChainMCallGetter{MChainInfo: _chainInfo}
}

func checkAddress(token common.Address) bool {
	return len(token.Hex()) == 42
}

func (mc *ChainMCallGetter) GetNftInfoForContract(contract common.Address) *types.NftInfo {
	if mc == nil {
		return nil
	}

	mc.McCounter++
	if nil == mc || nil == mc.MChainInfo {
		return nil
	}

	if mc.McCounter >= 3 {
		return nil
	}

	var vcs multicall.ViewCalls
	if !checkAddress(contract) {
		return nil
	}

	tokenURI := mc.newViewCall("tokenId", contract.Hex(), "tokenURI(uint256)(string)", []interface{}{1})
	totalSupply := mc.newViewCall("totalSupply", contract.Hex(), "totalSupply()(uint256)", []interface{}{})

	vcs = append(vcs, tokenURI)
	vcs = append(vcs, totalSupply)

	callRes, err := mc.callMultiContract(&vcs)
	if err != nil {
		time.Sleep(1 * time.Second)
		return mc.GetNftInfoForContract(contract)
	}

	if len(callRes) != 2 {
		time.Sleep(1 * time.Second)
		return mc.GetNftInfoForContract(contract)
	}

	nftInfo := &types.NftInfo{
		URI:         callRes[0].(string),
		TotalSupply: callRes[1].(*big.Int),
	}
	return nftInfo
}

func (mc *ChainMCallGetter) newViewCall(id, target, method string, arguments []interface{}) multicall.ViewCall {
	return multicall.NewViewCall(id, target, method, arguments)
}

func bigInt2String(num *multicall.BigIntJSONString) string {
	byteNum, _ := num.MarshalJSON()
	splitStrNum := strings.Split(string(byteNum)[1:], `"`)
	return splitStrNum[0]
}

func (mc *ChainMCallGetter) callMultiContract(vcs *multicall.ViewCalls) ([]interface{}, error) {
	var callRes []interface{}
	block := "latest"
	mcClient := *mc.MChainInfo.McClient
	res, err := mcClient.Call(*vcs, block)
	if err != nil {
		return nil, err
	}

	if nil == res || len(res.Calls) == 0 {
		return nil, types.McCallResIsNullErr
	}

	for _, item := range res.Calls {
		if !item.Success {
			break
		}

		if len(item.Decoded) == 0 {
			// TODO: shouldn't common.address
			callRes = append(callRes, common.HexToAddress("xxxxx"))
		}

		for _, value := range item.Decoded {
			if reflect.TypeOf(value).String() == "*multicall.BigIntJSONString" {
				valStr := bigInt2String(value.(*multicall.BigIntJSONString))
				n := new(big.Int)
				val, _ := n.SetString(valStr, 10)
				callRes = append(callRes, val)
				break
			}
			callRes = append(callRes, value)
		}
	}
	return callRes, nil
}
