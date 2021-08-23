package gql

import (
	"context"
	"regexp"
	"time"
)

var reValueNotConserved = regexp.MustCompile(`ValueNotConservedUTxO\s*\(Value\s+0`)

func (r *Resolver) WalletCreate(ctx context.Context, args struct{ InitialFunds *string }) (string, error) {
	var initialFunds string
	if args.InitialFunds != nil {
		initialFunds = *args.InitialFunds
	}

	return r.config.CLI.CreateWallet(ctx, initialFunds)
}

type WalletFundArgs struct {
	Address  string
	Quantity string
}

func (r *Resolver) WalletFund(ctx context.Context, args WalletFundArgs) (*Resolver, error) {

	for attempt := 1; attempt <= 5; attempt++ {
		if err := r.fundWallet(ctx, args.Address, args.Quantity); err != nil {
			if reValueNotConserved.MatchString(err.Error()) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(2 * time.Duration(attempt) * time.Second):
					continue
				}
			}
			return nil, err
		}
		return r, nil
	}

	return r, nil
}
