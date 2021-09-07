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
	Delegation   *string
}

func (r *Resolver) WalletCreate(ctx context.Context, args WalletCreateArgs) (string, error) {
	var initialFunds string
	if args.InitialFunds != nil {
		initialFunds = *args.InitialFunds
	}

	s, err := r.config.CLI.CreateWallet(ctx, initialFunds, StringValue(args.Name))
	if err != nil {
		return s, err
	}
	if *args.Delegation == "NONE" {
		return s, nil
	}
	_, err = r.WalletRegister(ctx, WalletRegisterArgs{Address: StringValue(args.Name)})
	if err != nil {
		return s, err
	}
	if *args.Delegation == "REGISTERED" {
		return s, nil
	}
	_, err = r.WalletDelegate(ctx, WalletDelegateArgs{Address: StringValue(args.Name)})
	if err != nil {
		return s, err
	}
	return s, nil
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

type WalletRegisterArgs struct {
	Address string
}

func (r *Resolver) WalletRegister(ctx context.Context, args WalletRegisterArgs) (*Resolver, error) {
	_, err := r.config.CLI.RegisterStake(ctx, args.Address)
	return r, err
}

type WalletDelegateArgs struct {
	Address string
}

func (r *Resolver) WalletDelegate(ctx context.Context, args WalletDelegateArgs) (*Resolver, error) {
	_, err := r.config.CLI.Delegate(ctx, args.Address)
	return r, err
}
