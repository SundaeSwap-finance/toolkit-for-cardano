package gql

import "github.com/SundaeSwap-finance/cardano-toolkit/internal/cardano"

type UtxoResolver struct {
	utxo cardano.Utxo
}

func (u *UtxoResolver) Address() string { return u.utxo.Address }

func (u *UtxoResolver) DatumHash() *string {
	if u.utxo.DatumHash == "" {
		return nil
	}
	return &u.utxo.DatumHash
}

func (u *UtxoResolver) Index() int32 { return u.utxo.Index }

func (u *UtxoResolver) Value() string { return u.utxo.Value }

func (u *UtxoResolver) Tokens() []*TokenResolver {
	var resolvers []*TokenResolver
	for _, token := range u.utxo.Tokens {
		resolvers = append(resolvers, &TokenResolver{token: token})
	}
	return resolvers
}
