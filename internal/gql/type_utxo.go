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

import "github.com/SundaeSwap-finance/toolkit-for-cardano/internal/cardano"

type UtxoResolver struct {
	utxo cardano.Utxo
}

func (u *UtxoResolver) Address() string { return u.utxo.Address }

func (u *UtxoResolver) DatumHash() *string {
	if u.utxo.DatumHash == "" {
		return nil
	}
	return &u.utxo.DatumHash
}

func (u *UtxoResolver) Index() int32 { return u.utxo.Index }

func (u *UtxoResolver) Value() string { return u.utxo.Value }

func (u *UtxoResolver) Tokens() []*TokenResolver {
	var resolvers []*TokenResolver
	for _, token := range u.utxo.Tokens {
		resolvers = append(resolvers, &TokenResolver{token: token})
	}
	return resolvers
}
