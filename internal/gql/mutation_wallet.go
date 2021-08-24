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
	"regexp"
	"time"
)

var reValueNotConserved = regexp.MustCompile(`ValueNotConservedUTxO\s*\(Value\s+0`)

type WalletCreateArgs struct {
	InitialFunds *string
	Name         *string
}

func (r *Resolver) WalletCreate(ctx context.Context, args WalletCreateArgs) (string, error) {
	var initialFunds string
	if args.InitialFunds != nil {
		initialFunds = *args.InitialFunds
	}

	return r.config.CLI.CreateWallet(ctx, initialFunds, StringValue(args.Name))
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
