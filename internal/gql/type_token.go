package gql

import "github.com/SundaeSwap-finance/cardano-toolkit/internal/cardano"

type TokenResolver struct {
	token cardano.Token
}

func (t *TokenResolver) Asset() *AssetResolver { return &AssetResolver{asset: t.token.Asset} }
func (t *TokenResolver) Quantity() string      { return t.token.Quantity }
