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
	"encoding/base64"
	"fmt"
	"time"

	"github.com/SundaeSwap-finance/cardano-toolkit/internal/cardano"
)

type MintArgs struct {
	AssetName string
	Quantity  string
	Wallet    string
}

func (r *Resolver) Mint(ctx context.Context, args MintArgs) (*Resolver, error) {
	utxos, err := r.config.CLI.Utxos(args.Wallet, excludeScripts(true), excludeTokens(true))
	if err != nil {
		return nil, err
	}

	if len(utxos) < 2 {
		if err := r.fundWallet(ctx, args.Wallet, "10000000"); err != nil {
			return nil, fmt.Errorf("failed to mint tokens: %w", err)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
			// ok
		}

		utxos, err = r.config.CLI.Utxos(args.Wallet, excludeScripts(true), excludeTokens(true))
		if err != nil {
			return nil, err
		}
	}

	input := BuildMintTxInput{
		Fee:       "0",
		Quantity:  args.Quantity,
		AssetName: args.AssetName,
		TxIn:      utxos[0],
		Wallet:    args.Wallet,
	}
	raw, err := r.buildMintTx(ctx, input)
	if err != nil {
		return nil, err
	}

	feeArgs := TxFeeArgs{
		Raw:       base64.StdEncoding.EncodeToString(raw),
		TxIn:      1,
		TxOut:     1,
		Witnesses: 1,
	}
	fee, err := r.TxFee(ctx, feeArgs)
	if err != nil {
		return nil, err
	}

	input.Fee = fee
	raw, err = r.buildMintTx(ctx, input)
	if err != nil {
		return nil, err
	}

	signArgs := TxSignArgs{
		Raw:    base64.StdEncoding.EncodeToString(raw),
		Wallet: args.Wallet,
	}
	signed, err := r.TxSign(ctx, signArgs)
	if err != nil {
		return nil, err
	}

	submitArgs := TxSubmitArgs{Signed: signed.body}
	if _, err := r.TxSubmit(ctx, submitArgs); err != nil {
		return nil, err
	}

	return r, nil
}

func excludeTokens(enabled bool) func(utxo cardano.Utxo) bool {
	if !enabled {
		return func(utxo cardano.Utxo) bool { return false }
	}
	return func(utxo cardano.Utxo) bool {
		return len(utxo.Tokens) > 0
	}
}

func excludeScripts(enabled bool) func(utxo cardano.Utxo) bool {
	if !enabled {
		return func(utxo cardano.Utxo) bool { return false }
	}
	return func(utxo cardano.Utxo) bool {
		return len(utxo.DatumHash) > 0
	}
}
