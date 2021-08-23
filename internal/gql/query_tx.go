package gql

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const delay = time.Millisecond * 1500

type TxFeeArgs struct {
	Raw       string
	TxIn      int32
	TxOut     int32
	Witnesses int32
}

func (r *Resolver) TxFee(ctx context.Context, args TxFeeArgs) (string, error) {
	data, err := base64.StdEncoding.DecodeString(args.Raw)
	if err != nil {
		return "", fmt.Errorf("failed to calculate fee: %w", err)
	}

	f, err := ioutil.TempFile(filepath.Join(r.config.CLI.Dir, "/tmp"), "script")
	if err != nil {
		return "", fmt.Errorf("failed to calculate fee: %w", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return "", fmt.Errorf("failed to calculate fee: %w", err)
	}

	return r.config.CLI.MinFee(ctx, f.Name(), args.TxIn, args.TxOut, args.Witnesses)
}
