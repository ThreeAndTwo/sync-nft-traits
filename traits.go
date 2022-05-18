package sync_nft_traits

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeAndTwo/sync-nft-traits/chains"
	"github.com/ThreeAndTwo/sync-nft-traits/types"
	"strconv"
)

type (
	ITraits interface {
		Sync() error
		Calc() (map[string]*types.NftData, error)
	}
)

type traits struct {
	chain       chains.IChain
	uri         string
	totalSupply int64
	nftMap      map[string]*types.NftData
}

func initHeader() map[string]string {
	header := make(map[string]string)
	header["content-type"] = "application/json"
	return header
}

func (t *traits) checkSync(config *types.NftInfo) bool {
	return config.URI != "" && config.TotalSupply != nil
}

//ipfs://QmYGgEFqTRkWvNZ6u7gfk9HDdh55bQAbYVyc16TF1zX658/
func fmtBaseTokenURI(nftInfo *types.NftInfo) {
	if nftInfo.URI[0:6] == "ipfs://" {
		nftInfo.URI = "https://infura-ipfs.io/ipfs/" + nftInfo.URI[7:len(nftInfo.URI)-2]
		return
	}

	nftInfo.URI = nftInfo.URI[0 : len(nftInfo.URI)-2]
}

func (t *traits) Sync() error {
	nftInfo, err := t.chain.GetNftInfo()
	if err != nil {
		return err
	}

	if !t.checkSync(nftInfo) {
		return types.ReadDataOnChainErr
	}

	if nftInfo == nil {
		return types.ReadDataOnChainErr
	}

	// reset baseTokenURI
	fmtBaseTokenURI(nftInfo)

	t.uri = nftInfo.URI
	t.totalSupply = nftInfo.TotalSupply.Int64()
	t.nftMap = make(map[string]*types.NftData)
	for k := int64(0); k < t.totalSupply; k++ {
		if err = t.getTrait(k); err != nil {
			// clean cache
			t.nftMap = map[string]*types.NftData{}
			return err
		}
	}
	return nil
}

// getTrait 区分大小写
func (t *traits) getTrait(index int64) error {
	fmt.Printf("index: %d \n", index)
	resp, err := newNet(fmt.Sprintf("%s/%d", t.uri, index), initHeader(), nil).get(initHeader())
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	if err = json.Unmarshal([]byte(resp), &data); err != nil {
		return err
	}

	if _, ok := data["attributes"]; !ok {
		return fmt.Errorf("non attributes for %d, data: %s error", index, resp)
	}

	attr := data["attributes"].([]interface{})

	for _, item := range attr {

		detail := item.(map[string]interface{})
		key := detail["trait_type"].(string) + ":" + detail["value"].(string)

		_, ok := t.nftMap[key]
		if ok {
			t.nftMap[key].Count++
			continue
		}

		_data := &types.NftData{
			Key:         detail["trait_type"].(string),
			Val:         detail["value"].(string),
			Count:       1,
			TotalSupply: t.totalSupply,
		}
		t.nftMap[key] = _data
	}
	return nil
}

func (t *traits) Calc() (map[string]*types.NftData, error) {
	if t.nftMap == nil {
		return nil, fmt.Errorf("nft attribute is null")
	}

	if t.totalSupply == 0 {
		return nil, fmt.Errorf("totalSupply should be gt 0")
	}

	for _, data := range t.nftMap {
		rate, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", float64(data.Count)/float64(t.totalSupply)), 64)
		data.RarityRate = rate
	}
	return t.nftMap, nil
}

func NewTraits(chainId int64, rpc, contract, mcContract string) (ITraits, error) {
	if chainId <= 0 || rpc == "" || contract == "" || mcContract == "" {
		return nil, types.ParamsErr
	}

	iChain, err := chains.NewChains(chainId, rpc, contract, mcContract)
	if err != nil {
		return nil, err
	}
	return &traits{chain: iChain}, err
}
