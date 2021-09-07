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

package cardano

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/savaki/zapctx"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
)

const (
	dirWallets = "wallets" // dirWallets contains wallets
)

// reWalletName ensures folks don't play games with the names
var reWalletName = regexp.MustCompile(`^[a-zA-Z0-9.\-_ ']*$`)

func (c CLI) CreateWallet(ctx context.Context, initialFunds, name string) (wallet string, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("created wallet",
			zap.String("name", name),
			zap.Duration("elapsed", time.Now().Sub(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())

	if !reWalletName.MatchString(name) {
		return "", fmt.Errorf("failed to create wallet: invalid name, must match ^[a-zA-Z0-9.\\-_ ']*$")
	}
	if name == "" {
		name = ksuid.New().String()
	}
	if err := os.MkdirAll(filepath.Join(c.Dir, dirWallets), 0755); err != nil {
		return "", fmt.Errorf("unable to create wallet: unable to create directory, %v: %w", dirWallets, err)
	}

	// don't allow existing wallets to be overwritten
	if _, err := os.Stat(fmt.Sprintf("%v/%v.skey", dirWallets, name)); !os.IsNotExist(err) {
		return "", fmt.Errorf("unable to create wallet, %v: wallet already exists", name)
	}

	// Payment address keys
	// cardano-cli address key-gen \
	// --verification-key-file addresses/${ADDR}.vkey \
	// --signing-key-file      addresses/${ADDR}.skey
	_, err = c.exec(
		"address", "key-gen",
		"--verification-key-file", fmt.Sprintf("%v/%v.vkey", dirWallets, name),
		"--signing-key-file", fmt.Sprintf("%v/%v.skey", dirWallets, name),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: failed to create payment address keys: %w", err)
	}

	// Stake address keys
	// cardano-cli stake-address key-gen \
	// --verification-key-file addresses/${ADDR}-stake.vkey \
	// --signing-key-file      addresses/${ADDR}-stake.skey
	_, err = c.exec(
		"stake-address", "key-gen",
		"--verification-key-file", fmt.Sprintf("%v/%v-stake.vkey", dirWallets, name),
		"--signing-key-file", fmt.Sprintf("%v/%v-stake.skey", dirWallets, name),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: failed to create stake address key: %w", err)
	}

	// Payment addresses
	// cardano-cli address build \
	// --payment-verification-key-file addresses/${ADDR}.vkey \
	// --stake-verification-key-file addresses/${ADDR}-stake.vkey \
	// --testnet-magic 42 \
	// --out-file addresses/${ADDR}.addr
	_, err = c.exec(
		"address", "build",
		"--payment-verification-key-file", fmt.Sprintf("%v/%v.vkey", dirWallets, name),
		"--stake-verification-key-file", fmt.Sprintf("%v/%v-stake.vkey", dirWallets, name),
		"--testnet-magic", c.TestnetMagic,
		"--out-file", fmt.Sprintf("%v/%v.addr", dirWallets, name),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: failed to create payment address: %w", err)
	}

	// Stake addresses
	// cardano-cli stake-address build \
	// --stake-verification-key-file addresses/${ADDR}-stake.vkey \
	// --testnet-magic 42 \
	// --out-file addresses/${ADDR}-stake.addr
	_, err = c.exec(
		"stake-address", "build",
		"--stake-verification-key-file", fmt.Sprintf("%v/%v-stake.vkey", dirWallets, name),
		"--testnet-magic", c.TestnetMagic,
		"--out-file", fmt.Sprintf("%v/%v-stake.addr", dirWallets, name),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: failed to create stake address: %w", err)
	}

	// Stake addresses registration certs
	// cardano-cli stake-address registration-certificate \
	// --stake-verification-key-file addresses/${ADDR}-stake.vkey \
	// --out-file addresses/${ADDR}-stake.reg.cert
	_, err = c.exec(
		"stake-address", "registration-certificate",
		"--stake-verification-key-file", fmt.Sprintf("%v/%v-stake.vkey", dirWallets, name),
		"--out-file", fmt.Sprintf("%v/%v-stake.reg.cert", dirWallets, name),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: failed to create stake address registration cert: %w", err)
	}

	// Stake delegation certs
	// cardano-cli stake-address delegation-certificate \
	// --stake-verification-key-file addresses/${ADDR}-stake.vkey
	// --cold-verification-key-file shelley/operator.vkey
	// --out-file addresses/${ADDR}-stake.delegate.cert
	_, err = c.exec(
		"stake-address", "delegation-certificate",
		"--stake-verification-key-file", fmt.Sprintf("%v/%v-stake.vkey", dirWallets, name),
		"--cold-verification-key-file", fmt.Sprintf("%v/shelley/operator.vkey", c.PoolDir),
		"--out-file", fmt.Sprintf("%v/%v-stake.delegate.cert", dirWallets, name),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: failed to create stake address delegation cert: %w", err)
	}

	if _, err := c.FundWallet(ctx, name, initialFunds); err != nil {
		return "", fmt.Errorf("failed to create wallet: %w", err)
	}

	return name, nil
}

var reQuantity = regexp.MustCompile(`^\d+$`)

func (c CLI) FundWallet(ctx context.Context, address, quantity string) (tx Tx, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("funded wallet",
			zap.String("address", address),
			zap.String("quantity", quantity),
			zap.Duration("elapsed", time.Since(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())

	address, err = c.NormalizeAddress(address)
	if err != nil {
		return Tx{}, err
	}

	if quantity == "" || quantity == "0" {
		return Tx{}, nil
	}

	if address == c.TreasuryAddr {
		return Tx{}, fmt.Errorf("unable to fund wallet: illegall attempt to fund treasury addr")
	}
	if !reQuantity.MatchString(quantity) {
		return Tx{}, fmt.Errorf("unable to fund wallet: invalid quantity, %v", quantity)
	}
	if len(quantity) > 13 {
		return Tx{}, fmt.Errorf("unable to fund wallet: quantity requested exceeds maximum 9,999,999,999")
	}

	utxos, err := c.Utxos(c.TreasuryAddr)
	if err != nil {
		return Tx{}, fmt.Errorf("unable to fund wallet: %w", err)
	}
	if len(utxos) == 0 {
		return Tx{}, fmt.Errorf("unable to fund wallet: treasury addr is empty, %v", c.TreasuryAddr)
	}

	for _, utxo := range utxos {
		if len(utxo.Value) <= 10 {
			continue
		}
		return c.transferFunds(ctx, utxo, address, quantity)
	}

	return Tx{}, fmt.Errorf("unable to fund wallet: insufficient funds in treasury addr, %v", c.TreasuryAddr)
}

func (c CLI) RegisterStake(ctx context.Context, address string) (tx Tx, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("registered stake address",
			zap.String("address", address),
			zap.Duration("elapsed", time.Since(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())
	addressMnemonic := address
	location := c.WalletLocation(addressMnemonic)
	address, err = c.NormalizeAddress(address)
	if err != nil {
		return Tx{}, err
	}

	// Fund the wallet with enough ADA to cover fees and registration deposit
	amt := big.NewInt(10000 * 1e6)
	tx, err = c.FundWallet(ctx, address, amt.String())
	if err != nil {
		return Tx{}, fmt.Errorf("failed to fund wallet: %w", err)
	}

	// Estimate the fee for the tx to register the stake fee
	cert := location + "-stake.reg.cert"
	raw, err := c.Build(
		TxIn(tx.ID, 1),
		TxOut(address, amt.String()),
		Certificate(cert),
	)
	if err != nil {
		return Tx{}, err
	}

	f, err := ioutil.TempFile(filepath.Join(c.DataDir(), "/tmp"), "script")
	if err != nil {
		return Tx{}, err
	}

	defer os.Remove(f.Name())
	defer f.Close()

	if _, err := f.Write(raw); err != nil {
		return Tx{}, err
	}

	fee, err := c.MinFee(ctx, f.Name(), 1, 1, 2)
	if err != nil {
		return Tx{}, err
	}
	feeValue, _ := big.NewInt(0).SetString(fee, 10)
	amt = big.NewInt(0).Sub(amt, feeValue)
	// Build, Sign, and Submit
	raw, err = c.Build(
		TxIn(tx.ID, 1),
		TxOut(address, amt.String()),
		Fee(fee),
		Certificate(cert),
	)
	if err != nil {
		return Tx{}, err
	}
	signed, err := c.Sign(
		ctx,
		raw,
		addressMnemonic,
		addressMnemonic+"-stake",
	)
	if err != nil {
		return Tx{}, err
	}
	tx, err = ParseTx(signed)
	if err != nil {
		return Tx{}, err
	}
	err = c.Submit(ctx, signed)
	if err != nil {
		return Tx{}, err
	}
	return tx, nil
}

func (c CLI) Delegate(ctx context.Context, address string) (tx Tx, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("delegated stake",
			zap.String("address", address),
			zap.Duration("elapsed", time.Since(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())
	addressMnemonic := address
	location := c.WalletLocation(addressMnemonic)
	address, err = c.NormalizeAddress(address)
	if err != nil {
		return Tx{}, err
	}

	// Find a utxo containing at least 2 ada
	amt := big.NewInt(2 * 1e6)
	utxos, err := c.Utxos(
		address,
		AtLeast(int32(amt.Int64())), // bounds checking?
		ExcludeScripts(true),
		ExcludeTokens(true),
	)
	if err != nil {
		return Tx{}, err
	}
	utxo := Utxo{}
	if len(utxos) < 1 {
		tx, err = c.FundWallet(ctx, address, amt.String())
		if err != nil {
			return Tx{}, fmt.Errorf("failed to fund wallet: %w", err)
		}
		utxo = Utxo{
			Address: tx.ID,
			Index:   1,
			Value:   "2000000",
		}
	} else {
		utxo = utxos[0]
		amt, _ = big.NewInt(0).SetString(utxo.Value, 10)
	}

	// Estimate the fee for the tx to register the stake fee
	cert := location + "-stake.delegate.cert"
	raw, err := c.Build(
		TxIn(utxo.Address, utxo.Index),
		TxOut(address, amt.String()),
		Certificate(cert),
	)
	if err != nil {
		return Tx{}, err
	}

	f, err := ioutil.TempFile(filepath.Join(c.DataDir(), "/tmp"), "script")
	if err != nil {
		return Tx{}, err
	}

	defer os.Remove(f.Name())
	defer f.Close()

	if _, err := f.Write(raw); err != nil {
		return Tx{}, err
	}

	fee, err := c.MinFee(ctx, f.Name(), 1, 1, 2)
	if err != nil {
		return Tx{}, err
	}
	feeValue, _ := big.NewInt(0).SetString(fee, 10)
	amt = big.NewInt(0).Sub(amt, feeValue)
	// Build, Sign, and Submit
	raw, err = c.Build(
		TxIn(utxo.Address, utxo.Index),
		TxOut(address, amt.String()),
		Fee(fee),
		Certificate(cert),
	)
	if err != nil {
		return Tx{}, err
	}
	signed, err := c.Sign(
		ctx,
		raw,
		addressMnemonic,
		addressMnemonic+"-stake",
	)
	if err != nil {
		return Tx{}, err
	}
	tx, err = ParseTx(signed)
	if err != nil {
		return Tx{}, err
	}
	err = c.Submit(ctx, signed)
	if err != nil {
		return Tx{}, err
	}
	return tx, nil
}
