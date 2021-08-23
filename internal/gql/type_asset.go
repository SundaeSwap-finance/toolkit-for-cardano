package gql

import "github.com/SundaeSwap-finance/cardano-toolkit/internal/cardano"

type AssetResolver struct {
	asset *cardano.Asset
}

func (a *AssetResolver) AssetId() string   { return a.asset.ID() }
func (a *AssetResolver) AssetName() string { return a.asset.AssetName }
func (a *AssetResolver) Name() *string     { return StringPtr(a.asset.Name) }
func (a *AssetResolver) PolicyId() string  { return a.asset.PolicyId }
func (a *AssetResolver) Ticker() *string   { return StringPtr(a.asset.Ticker) }
func (a *AssetResolver) Url() *string      { return StringPtr(a.asset.Url) }
