package sync_nft_traits

import (
	"encoding/json"
	"testing"
)

const rpc = "https://eth-hk1.csnodes.com/v1/973eeba6738a7d8c3bd54f91adcbea89"
const mcContract = "0x5eb3fa2dfecdde21c950813c665e9364fa609bd2"

func TestNewTraits(t *testing.T) {
	tests := []struct {
		name     string
		chainId  int64
		contract string
	}{
		{
			name:     "test duskbreaks: using cdn",
			chainId:  1,
			contract: "0x0bEed7099AF7514cCEDF642CfEA435731176Fb02",
		},
		{
			name:     "test TIMEPieces: using ipfs",
			chainId:  1,
			contract: "0xDd69da9a83ceDc730bc4d3C56E96D29Acc05eCDE",
		},
		{
			name:     "test aaa: not set tokenURI",
			chainId:  1,
			contract: "0x7cad06b811b5d9d3ff197c1a046abcbc0efbcbc9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nt, err := NewTraits(tt.chainId, rpc, tt.contract, mcContract)
			if err != nil {
				t.Fatalf("new traits error: %s", err)
			}

			err = nt.Sync()
			if err != nil {
				t.Fatalf("sync nft error: %s", err)
			}

			calc, err := nt.Calc()
			if err != nil {
				t.Fatalf("Calc nft rate error: %s", err)
			}

			mData, _ := json.Marshal(calc)
			t.Logf("mData: %s", string(mData))
		})
	}
}
