package gql

import "github.com/SundaeSwap-finance/cardano-toolkit/internal/cardano"

type TipResolver struct {
	tip *cardano.Tip
}

func (t *TipResolver) Block() int32 { return t.tip.Block }
func (t *TipResolver) Epoch() int32 { return t.tip.Epoch }
func (t *TipResolver) Era() string  { return t.tip.Era }
func (t *TipResolver) Hash() string { return t.tip.Hash }
func (t *TipResolver) Slot() int32  { return t.tip.Slot }
