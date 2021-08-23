package gql

import (
	"strings"

	"github.com/SundaeSwap-finance/cardano-toolkit/internal/cardano"
)

type UtxosArgs struct {
	AssetId        *string
	Address        *string
	ExcludeScripts *bool
	ExcludeTokens  *bool
}

func (r *Resolver) Utxos(args UtxosArgs) ([]*UtxoResolver, error) {
	var address string
	if args.Address != nil {
		address = *args.Address
	}

	utxos, err := r.config.CLI.Utxos(address,
		onlyAssetId(args.AssetId),
		excludeScripts(boolValue(args.ExcludeScripts)),
		excludeTokens(boolValue(args.ExcludeTokens)),
	)
	if err != nil {
		return nil, err
	}

	var resolvers []*UtxoResolver
	for _, utxo := range utxos {
		resolvers = append(resolvers, &UtxoResolver{utxo: utxo})
	}

	return resolvers, nil
}

func boolValue(b *bool) bool {
	return b != nil && *b
}

func onlyAssetId(s *string) func(utxo cardano.Utxo) bool {
	if s == nil {
		return func(utxo cardano.Utxo) bool { return false }
	}

	parts := strings.SplitN(*s, ".", 2)
	if len(parts) != 2 {
		return func(utxo cardano.Utxo) bool { return false }
	}

	policyId, assetName := parts[0], parts[1]
	return func(utxo cardano.Utxo) bool {
		for _, token := range utxo.Tokens {
			if asset := token.Asset; asset != nil {
				if asset.AssetName == assetName && asset.PolicyId == policyId {
					return false
				}
			}
		}
		return true
	}
}
