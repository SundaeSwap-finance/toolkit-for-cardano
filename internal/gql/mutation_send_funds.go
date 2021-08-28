// MIT License
//
// Copyright (c) 2021 SundaeSwap Finance
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gql

import (
	"context"
	"fmt"
	"math/big"

	"github.com/SundaeSwap-finance/toolkit-for-cardano/internal/cardano"
)

type SendFundArgs struct {
	Source string
	Target *string
	TxIn   []TxIn
}

func (r *Resolver) SendFunds(ctx context.Context, args SendFundArgs) (*Resolver, error) {
	utxos, err := r.config.CLI.Utxos("", excludeScripts(true))
	if err != nil {
		return nil, err
	}

	var (
		sum     = big.NewInt(0)
		tokens  = map[string]*big.Int{}
		options []cardano.BuildOption
	)
	for _, txIn := range args.TxIn {
		utxo, err := utxos.Find(txIn.Address, txIn.Index)
		if err != nil {
			return nil, fmt.Errorf("sendFunds failed: %w", err)
		}

		v, ok := big.NewInt(0).SetString(utxo.Value, 10)
		if !ok {
			return nil, fmt.Errorf("failed to parse utxo, %v", utxo.Value)
		}

		sum = big.NewInt(0).Add(sum, v)
		for _, token := range utxo.Tokens {
			v, ok := big.NewInt(0).SetString(token.Quantity, 10)
			if !ok {
				return nil, fmt.Errorf("failed to parse token quantity, %v", token.Quantity)
			}

			assetID := token.Asset.ID()
			total := tokens[assetID]
			if total == nil {
				total = big.NewInt(0)
			}
			tokens[assetID] = big.NewInt(0).Add(total, v)
		}

		options = append(options, cardano.TxIn(txIn.Address, txIn.Index))
	}

	if len(args.TxIn) > 0 {
		address := args.Source
		if args.Target != nil {
			if s := *args.Target; s != "" {
				address = s
			}
		}

		var tt []string
		for assetID, amount := range tokens {
			tt = append(tt, amount.String()+" "+assetID)
		}
		options = append(options, cardano.TxOut(address, sum.String(), tt...))
	}
	options = append(options, cardano.Fee("0"))

	raw, err := r.config.CLI.Build(options...)
	if err != nil {
		return nil, err
	}

	signed, err := r.config.CLI.Sign(ctx, raw, args.Source)
	if err != nil {
		return nil, err
	}

	if err := r.config.CLI.Submit(ctx, signed); err != nil {
		return nil, err
	}

	return r, nil
}
