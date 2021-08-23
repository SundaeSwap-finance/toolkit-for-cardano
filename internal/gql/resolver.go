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
