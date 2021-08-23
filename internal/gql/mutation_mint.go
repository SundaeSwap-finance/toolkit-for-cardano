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
