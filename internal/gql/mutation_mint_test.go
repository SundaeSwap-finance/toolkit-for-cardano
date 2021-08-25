package gql

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/SundaeSwap-finance/cardano-toolkit/internal/cardano"
	"github.com/segmentio/ksuid"
	"github.com/tj/assert"
)

type Mock struct {
	Cardano
	fee      string // fee is the value returned by #MinFee
	quantity string // quantity of lovelace every utxo returned by #Utxos will have
	options  []cardano.BuildOptions
}

func (m *Mock) Build(opts ...cardano.BuildOption) ([]byte, error) {
	m.options = append(m.options, cardano.MakeBuildOptions(opts...))
	return []byte("{}"), nil
}

func (m Mock) DataDir() string {
	_ = os.MkdirAll("/tmp/data", 0755)
	_ = os.MkdirAll("/tmp/tmp", 0755)
	return "/tmp"
}

func (m Mock) KeyHash(ctx context.Context, wallet string) (keyHash string, err error) {
	return "KeyHash", nil
}

func (m Mock) MinFee(ctx context.Context, filename string, txIn, txOut, witnesses int32) (fee string, err error) {
	if m.fee == "" {
		return "0", nil
	}
	return m.fee, nil
}

func (m Mock) NormalizeAddress(address string) (string, error) {
	return strings.ToLower(address), nil
}

func (m Mock) PolicyID(ctx context.Context, filename string) (policyID string, err error) {
	return "PolicyID", nil
}

func (m Mock) Sign(ctx context.Context, raw []byte, wallets ...string) (data []byte, err error) {
	return raw, nil
}

func (m Mock) Submit(ctx context.Context, signed []byte) (err error) {
	return nil
}

func (m Mock) Utxos(address string, excludes ...func(cardano.Utxo) bool) (utxos cardano.Utxos, err error) {
	txHash := ksuid.New().String()
	return cardano.Utxos{
		{
			Address: txHash,
			Index:   0,
			Value:   m.quantity,
		},
		{
			Address: txHash,
			Index:   1,
			Value:   m.quantity,
		},
	}, nil
}

func TestResolver_Mint(t *testing.T) {
	t.Run("with fee", func(t *testing.T) {
		var (
			ctx  = context.Background()
			mock = &Mock{
				fee:      "1000",
				quantity: "10000000",
			}
			config   = Config{CLI: mock}
			resolver = &Resolver{config: config}
		)

		args := MintArgs{
			AssetName: "BLAH",
			Quantity:  "100",
			Wallet:    "Test",
		}
		_, err := resolver.Mint(ctx, args)
		assert.Nil(t, err)
		assert.Len(t, mock.options, 2)

		option := mock.options[1]
		assert.Equal(t, option.TxOut[0].Quantity, "9999000")
	})

	t.Run("no fees", func(t *testing.T) {
		var (
			ctx  = context.Background()
			mock = &Mock{
				fee:      "0",
				quantity: "10000000",
			}
			config   = Config{CLI: mock}
			resolver = &Resolver{config: config}
		)

		args := MintArgs{
			AssetName: "BLAH",
			Quantity:  "100",
			Wallet:    "Test",
		}
		_, err := resolver.Mint(ctx, args)
		assert.Nil(t, err)
		assert.Len(t, mock.options, 2)

		option := mock.options[1]
		assert.Equal(t, option.TxOut[0].Quantity, "10000000")
	})
}
