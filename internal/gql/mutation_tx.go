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
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/SundaeSwap-finance/toolkit-for-cardano/internal/cardano"
)

type TxIn struct {
	Address string
	Index   int32
}

func (t TxIn) String() string { return fmt.Sprintf("%v#%v", t.Address, t.Index) }

func (t TxIn) ToString(utxos cardano.Utxos) string {
	utxo, err := utxos.Find(t.Address, t.Index)
	if err != nil {
		return fmt.Sprintf("utxo-not-found:%v#%v", t.Address, t.Index)
	}

	buf := bytes.NewBuffer(nil)
	if utxo.DatumHash != "" {
		fmt.Fprintf(buf, "%v@", utxo.DatumHash)
	}

	fmt.Fprintf(buf, "%v#%v+%v:lovelace", utxo.Address, utxo.Index, utxo.Value)
	for _, token := range utxo.Tokens {
		fmt.Fprintf(buf, "+%v:%v", token.Quantity, token.Asset.ID())
	}
	return buf.String()
}

type TxInDatum struct {
	TxIn
	DatumHash string
}

func (t TxInDatum) String() string { return fmt.Sprintf("%v@%v#%v", t.DatumHash, t.Address, t.Index) }

type TxOut struct {
	Address  string
	Quantity string
}

type TxBuildArgs struct {
	Fee   string
	TxIn  []TxIn
	TxOut []TxOut
}

func (r *Resolver) TxBuild(args TxBuildArgs) (*TxResolver, error) {
	var opts []cardano.BuildOption
	opts = append(opts, cardano.Fee(args.Fee))
	for _, txIn := range args.TxIn {
		opts = append(opts, cardano.TxIn(txIn.Address, txIn.Index))
	}
	for _, txOut := range args.TxOut {
		opts = append(opts, cardano.TxOut(txOut.Address, txOut.Quantity))
	}

	data, err := r.config.CLI.Build(opts...)
	if err != nil {
		return nil, err
	}

	tx, err := cardano.ParseTx(data)
	if err != nil {
		return nil, err
	}

	return &TxResolver{
		body: base64.StdEncoding.EncodeToString(data),
		id:   tx.ID,
		raw:  data,
	}, nil
}

type TxSignArgs struct {
	Raw    string
	Wallet string
}

func (r *Resolver) TxSign(ctx context.Context, args TxSignArgs) (*TxResolver, error) {
	data, err := base64.StdEncoding.DecodeString(args.Raw)
	if err != nil {
		return nil, fmt.Errorf("unable to sign tx: unable to base64 decode string: %w", err)
	}

	data, err = r.config.CLI.Sign(ctx, data, args.Wallet)
	if err != nil {
		return nil, err
	}

	tx, err := cardano.ParseTx(data)
	if err != nil {
		return nil, err
	}

	return &TxResolver{
		body: base64.StdEncoding.EncodeToString(data),
		id:   tx.ID,
	}, nil
}

type TxSubmitArgs struct {
	Signed string
}

func (r *Resolver) TxSubmit(ctx context.Context, args TxSubmitArgs) (*Resolver, error) {
	data, err := base64.StdEncoding.DecodeString(args.Signed)
	if err != nil {
		return nil, fmt.Errorf("unable to submit tx: unable to base64 decode string: %w", err)
	}

	if err := r.config.CLI.Submit(ctx, data); err != nil {
		return nil, err
	}

	return r, nil
}
