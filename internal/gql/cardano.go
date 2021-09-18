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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/SundaeSwap-finance/toolkit-for-cardano/internal/cardano"
	"github.com/savaki/zapctx"
	"go.uber.org/zap"
)

type BuildMintTxInput struct {
	AssetName string
	Fee       string
	Quantity  string
	TxIn      cardano.Utxo
	Wallet    string
}

func (r *Resolver) buildMintTx(ctx context.Context, input BuildMintTxInput) (raw []byte, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("build mint-tx",
			zap.Duration("elapsed", time.Now().Sub(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())
	keyHash, err := r.config.CLI.KeyHash(ctx, input.Wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to build mint tx: %w", err)
	}

	script, scriptCleanup, err := buildScript(r.config.CLI.DataDir(), keyHash)
	if err != nil {
		return nil, fmt.Errorf("failed to build mint tx: %w", err)
	}
	defer scriptCleanup()

	policyID, err := r.config.CLI.PolicyID(ctx, script)
	if err != nil {
		return nil, fmt.Errorf("failed to build mint tx: %w", err)
	}

	address, err := r.config.CLI.NormalizeAddress(input.Wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to build mint tx: %w", err)
	}

	value, ok := big.NewInt(0).SetString(input.TxIn.Value, 10)
	if !ok {
		return nil, fmt.Errorf("mint failed to parse txIn amount, %v", input.TxIn.Value)
	}

	fee, ok := big.NewInt(0).SetString(input.Fee, 10)
	if !ok {
		return nil, fmt.Errorf("mint failed to parse fee, %v", fee)
	}

	remain := big.NewInt(0).Sub(value, fee)

	mintedTokens := fmt.Sprintf("%v %v.%v", input.Quantity, policyID, input.AssetName)
	return r.config.CLI.Build(
		cardano.Fee(input.Fee),
		cardano.TxIn(input.TxIn.Address, input.TxIn.Index),
		cardano.TxOut(address, remain.String(), mintedTokens),
		cardano.Mint(mintedTokens),
		cardano.MintScriptFile(fmt.Sprintf(script)),
	)
}

func (r *Resolver) fundWallet(ctx context.Context, address, quantity string) (err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("funded wallet",
			zap.String("address", address),
			zap.String("quantity", quantity),
			zap.Duration("elapsed", time.Now().Sub(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())

	if _, err := r.config.CLI.FundWallet(ctx, address, quantity); err != nil {
		return err
	}

	return nil
}

func buildScript(dir, keyHash string) (string, func(), error) {
	f, err := ioutil.TempFile(filepath.Join(dir, "/tmp"), "script")
	if err != nil {
		return "", func() {}, fmt.Errorf("failed to create script: %w", err)
	}
	defer f.Close()

	type Record struct {
		Type    string `json:"type,omitempty"`
		KeyHash string `json:"keyHash,omitempty"`
	}
	record := Record{
		Type:    "sig",
		KeyHash: keyHash,
	}
	if err := json.NewEncoder(f).Encode(record); err != nil {
		f.Close()
		_ = os.Remove(f.Name())
		return "", func() {}, fmt.Errorf("failed to write token policy script: %w", err)
	}

	filename := f.Name()
	return filename, func() { _ = os.Remove(filename) }, nil
}
