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
	_ "embed"
	"fmt"
	"net/http"

	"github.com/SundaeSwap-finance/cardano-toolkit/internal/cardano"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

//go:embed schema.graphql
var textSchema string

type Config struct {
	Built   string
	CLI     *cardano.CLI
	Version string
}

type Resolver struct {
	config Config
}

func New(config Config) (http.Handler, error) {
	resolver := &Resolver{
		config: config,
	}

	schema, err := graphql.ParseSchema(textSchema, resolver)
	if err != nil {
		return nil, fmt.Errorf("unable to parse schema: %w", err)
	}

	return &relay.Handler{Schema: schema}, nil
}
