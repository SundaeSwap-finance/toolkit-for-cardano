package cardano

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/savaki/zapctx"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
)

const (
	dirWallets = "wallets" // dirWallets contains wallets
)

func (c CLI) CreateWallet(ctx context.Context, initialFunds string) (wallet string, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("created wallet",
			zap.Duration("elapsed", time.Now().Sub(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())

	addr := ksuid.New().String()
	if err := os.MkdirAll(filepath.Join(c.Dir, dirWallets), 0755); err != nil {
		return "", fmt.Errorf("unable to create wallet: unable to create directory, %v: %w", dirWallets, err)
	}

	// Payment address keys
	// cardano-cli address key-gen \
	// --verification-key-file addresses/${ADDR}.vkey \
	// --signing-key-file      addresses/${ADDR}.skey
	_, err = c.exec(
		"address", "key-gen",
		"--verification-key-file", fmt.Sprintf("%v/%v.vkey", dirWallets, addr),
		"--signing-key-file", fmt.Sprintf("%v/%v.skey", dirWallets, addr),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: failed to create payment address keys: %w", err)
	}

	// Stake address keys
	// cardano-cli stake-address key-gen \
	// --verification-key-file addresses/${ADDR}-stake.vkey \
	// --signing-key-file      addresses/${ADDR}-stake.skey
	_, err = c.exec(
		"stake-address", "key-gen",
		"--verification-key-file", fmt.Sprintf("%v/%v-stake.vkey", dirWallets, addr),
		"--signing-key-file", fmt.Sprintf("%v/%v-stake.skey", dirWallets, addr),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: failed to create stake address key: %w", err)
	}

	// Payment addresses
	// cardano-cli address build \
	// --payment-verification-key-file addresses/${ADDR}.vkey \
	// --stake-verification-key-file addresses/${ADDR}-stake.vkey \
	// --testnet-magic 42 \
	// --out-file addresses/${ADDR}.addr
	_, err = c.exec(
		"address", "build",
		"--payment-verification-key-file", fmt.Sprintf("%v/%v.vkey", dirWallets, addr),
		"--stake-verification-key-file", fmt.Sprintf("%v/%v-stake.vkey", dirWallets, addr),
		"--testnet-magic", c.TestnetMagic,
		"--out-file", fmt.Sprintf("%v/%v.addr", dirWallets, addr),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: failed to create payment address: %w", err)
	}

	// Stake addresses
	// cardano-cli stake-address build \
	// --stake-verification-key-file addresses/${ADDR}-stake.vkey \
	// --testnet-magic 42 \
	// --out-file addresses/${ADDR}-stake.addr
	_, err = c.exec(
		"stake-address", "build",
		"--stake-verification-key-file", fmt.Sprintf("%v/%v-stake.vkey", dirWallets, addr),
		"--testnet-magic", c.TestnetMagic,
		"--out-file", fmt.Sprintf("%v/%v-stake.addr", dirWallets, addr),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: failed to create stake address: %w", err)
	}

	// Stake addresses registration certs
	// cardano-cli stake-address registration-certificate \
	// --stake-verification-key-file addresses/${ADDR}-stake.vkey \
	// --out-file addresses/${ADDR}-stake.reg.cert
	_, err = c.exec(
		"stake-address", "registration-certificate",
		"--stake-verification-key-file", fmt.Sprintf("%v/%v-stake.vkey", dirWallets, addr),
		"--out-file", fmt.Sprintf("%v/%v-stake.reg.cert", dirWallets, addr),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: failed to create stake address registration cert: %w", err)
	}

	if _, err := c.FundWallet(ctx, addr, initialFunds); err != nil {
		return "", fmt.Errorf("failed to create wallet: %w", err)
	}

	return addr, nil
}

var reQuantity = regexp.MustCompile(`^\d+$`)

func (c CLI) FundWallet(ctx context.Context, address, quantity string) (tx Tx, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("funded wallet",
			zap.String("address", address),
			zap.String("quantity", quantity),
			zap.Duration("elapsed", time.Now().Sub(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())

	address, err = c.NormalizeAddress(address)
	if err != nil {
		return Tx{}, err
	}

	if quantity == "" || quantity == "0" {
		return Tx{}, nil
	}

	if address == c.TreasuryAddr {
		return Tx{}, fmt.Errorf("unable to fund wallet: illegall attempt to fund treasury addr")
	}
	if !reQuantity.MatchString(quantity) {
		return Tx{}, fmt.Errorf("unable to fund wallet: invalid quantity, %v", quantity)
	}
	if len(quantity) > 13 {
		return Tx{}, fmt.Errorf("unable to fund wallet: quantity requested exceeds maximum 9,999,999,999")
	}

	utxos, err := c.Utxos(c.TreasuryAddr)
	if err != nil {
		return Tx{}, fmt.Errorf("unable to fund wallet: %w", err)
	}
	if len(utxos) == 0 {
		return Tx{}, fmt.Errorf("unable to fund wallet: treasury addr is empty, %v", c.TreasuryAddr)
	}

	for _, utxo := range utxos {
		if len(utxo.Value) <= 10 {
			continue
		}
		return c.transferFunds(ctx, utxo, address, quantity)
	}

	return Tx{}, fmt.Errorf("unable to fund wallet: insufficient funds in treasury addr, %v", c.TreasuryAddr)
}
