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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/savaki/zapctx"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/blake2b"
)

type txIn struct {
	Address string
	Index   int32
}

type txOut struct {
	Address  string
	Quantity string
	Tokens   []string
}

type BuildOptions struct {
	fee            string
	mint           string
	mintScriptFile string
	txIn           []txIn
	txOut          []txOut
}

func buildOptions(opts ...BuildOption) BuildOptions {
	var options BuildOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.fee == "" {
		options.fee = "0"
	}

	return options
}

type BuildOption func(options *BuildOptions)

func Fee(fee string) BuildOption {
	return func(options *BuildOptions) {
		options.fee = fee
	}
}

func Mint(s string) BuildOption {
	return func(options *BuildOptions) {
		options.mint = s
	}
}

func MintScriptFile(filename string) BuildOption {
	return func(options *BuildOptions) {
		options.mintScriptFile = filename
	}
}

func TxIn(address string, index int32) BuildOption {
	return func(options *BuildOptions) {
		options.txIn = append(options.txIn, txIn{
			Address: address,
			Index:   index,
		})
	}
}

func TxOut(address string, quantity string, tokens ...string) BuildOption {
	return func(options *BuildOptions) {
		options.txOut = append(options.txOut, txOut{
			Address:  address,
			Quantity: quantity,
			Tokens:   tokens,
		})
	}
}

func (c CLI) Build(opts ...BuildOption) ([]byte, error) {
	filename := filepath.Join(c.Dir, "tmp", ksuid.New().String())
	if !c.Debug {
		defer func() { os.Remove(filename) }()
	}

	options := buildOptions(opts...)
	args := []string{
		"transaction", "build-raw",
		"--fee", options.fee,
		"--alonzo-era",
		"--out-file", filename,
	}

	for _, in := range options.txIn {
		args = append(args, "--tx-in", fmt.Sprintf("%v#%v", in.Address, in.Index))
	}
	for _, in := range options.txOut {
		address, err := c.NormalizeAddress(in.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to build tx: %w", err)
		}

		output := fmt.Sprintf("%v+%v", address, in.Quantity)
		if len(in.Tokens) > 0 {
			output += "+" + strings.Join(in.Tokens, "+")
		}
		args = append(args, "--tx-out", output)
	}
	if options.mint != "" {
		args = append(args, "--mint="+options.mint)
	}
	if options.mintScriptFile != "" {
		args = append(args, "--mint-script-file="+options.mintScriptFile)
	}

	fmt.Println()
	fmt.Println(strings.Join(c.Cmd, " "), strings.Join(args, " "))
	fmt.Println()

	if _, err := c.exec(args...); err != nil {
		return nil, fmt.Errorf("failed to build transaction: %w", err)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction: unable to read file, %v: %w", filename, err)
	}

	return data, nil
}

func (c CLI) Sign(ctx context.Context, raw []byte, wallets ...string) (data []byte, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("signed tx",
			zap.Duration("elapsed", time.Now().Sub(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())

	filename := filepath.Join(c.Dir, "tmp", ksuid.New().String())
	if !c.Debug {
		defer func() { os.Remove(filename) }()
	}

	if err := ioutil.WriteFile(filename, raw, 0644); err != nil {
		return nil, fmt.Errorf("failed to sign tx: unable to write file: %w", err)
	}

	args := []string{
		"transaction", "sign",
		"--tx-body-file", filename,
		"--out-file", filename,
	}
	for _, wallet := range wallets {
		var signingKeyFile string
		switch wallet {
		case "":
			signingKeyFile = c.TreasurySkeyFile
		default:
			signingKeyFile = fmt.Sprintf("%v/%v/%v.skey", c.Dir, dirWallets, wallet)
		}
		args = append(args, "--signing-key-file", signingKeyFile)
	}

	if _, err := c.exec(args...); err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: unable to read file, %v: %w", filename, err)
	}

	return data, nil
}

func (c CLI) Submit(ctx context.Context, signed []byte) (err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("submitted tx",
			zap.Duration("elapsed", time.Now().Sub(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())

	filename := filepath.Join(c.Dir, "tmp", ksuid.New().String())
	if !c.Debug {
		defer func() { os.Remove(filename) }()
	}

	if err := ioutil.WriteFile(filename, signed, 0644); err != nil {
		return fmt.Errorf("failed to submit tx: unable to write file: %w", err)
	}

	args := []string{
		"transaction", "submit",
		"--cardano-mode",
		"--testnet-magic", c.TestnetMagic,
		"--tx-file", filename,
	}
	if _, err := c.exec(args...); err != nil {
		return fmt.Errorf("failed to submit tx: %w", err)
	}

	return nil
}

func (c CLI) transferFunds(ctx context.Context, utxo Utxo, address, quantity string) (Tx, error) {
	total := &big.Int{}
	total, ok := total.SetString(utxo.Value, 10)
	if !ok {
		return Tx{}, fmt.Errorf("failed to transfer funds: unable to parse source utxo value, %v", utxo.Value)
	}

	q := &big.Int{}
	q, ok = q.SetString(quantity, 10)
	if !ok {
		return Tx{}, fmt.Errorf("failed to transfer funds: unable to parse desired quantity, %v", quantity)
	}

	remainder := &big.Int{}
	remainder = remainder.Sub(total, q)

	raw, err := c.Build(
		TxIn(utxo.Address, utxo.Index),
		TxOut(c.TreasuryAddr, remainder.Text(10)),
		TxOut(address, quantity),
	)
	if err != nil {
		return Tx{}, fmt.Errorf("failed to transfer funds: %w", err)
	}

	filename := filepath.Join(c.Dir, "tmp", ksuid.New().String())
	defer os.Remove(filename)
	if _ = ioutil.WriteFile(filename, raw, 0644); err != nil {
		return Tx{}, fmt.Errorf("failed to write raw tx body: %w", err)
	}

	feeStr, err := c.MinFee(ctx, filename, 1, 2, 1)
	if err != nil {
		return Tx{}, fmt.Errorf("failed to transfer funds: %w", err)
	}

	fee, ok := big.NewInt(0).SetString(feeStr, 10)
	if !ok {
		return Tx{}, fmt.Errorf("failed to parse fee, %v", feeStr)
	}

	remainder = big.NewInt(0).Sub(remainder, fee)

	raw, err = c.Build(
		Fee(feeStr),
		TxIn(utxo.Address, utxo.Index),
		TxOut(c.TreasuryAddr, remainder.Text(10)),
		TxOut(address, quantity),
	)
	if err != nil {
		return Tx{}, fmt.Errorf("failed to transfer funds: %w", err)
	}

	signed, err := c.Sign(ctx, raw, "")
	if err != nil {
		return Tx{}, fmt.Errorf("failed to transfer funds: %w", err)
	}

	tx, err := ParseTx(signed)
	if err != nil {
		return Tx{}, fmt.Errorf("failed to parse transaction: %w", err)
	}

	if err := c.Submit(ctx, signed); err != nil {
		return Tx{}, fmt.Errorf("failed to transfer funds: %w", err)
	}

	return tx, nil
}

type Tx struct {
	ID string
}

func ParseTx(data []byte) (Tx, error) {
	var file struct{ CborHex string }
	if err := json.Unmarshal(data, &file); err != nil {
		return Tx{}, fmt.Errorf("failed to get tx id: %w", err)
	}

	data, err := hex.DecodeString(file.CborHex)
	if err != nil {
		return Tx{}, fmt.Errorf("failed to decode cbor hex: %w", err)
	}

	type Record struct {
		_     struct{} `cbor:",toarray"`
		Body  cbor.RawMessage
		Junk1 cbor.RawMessage
		Junk2 cbor.RawMessage
	}
	var record Record
	h, err := blake2b.New256(nil)
	if err != nil {
		return Tx{}, fmt.Errorf("failed to get tx id: unable to create blake2b decoder: %w", err)
	}

	if _, err := h.Write(record.Body); err != nil {
		return Tx{}, fmt.Errorf("failed to write record body: %w", err)
	}

	hash := h.Sum(nil)
	return Tx{
		ID: hex.EncodeToString(hash),
	}, nil
}
