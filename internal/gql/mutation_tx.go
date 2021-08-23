package gql

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/SundaeSwap-finance/cardano-toolkit/internal/cardano"
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
