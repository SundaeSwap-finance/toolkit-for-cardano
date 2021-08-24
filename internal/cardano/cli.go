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

package cardano

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/savaki/zapctx"
	"go.uber.org/zap"
)

var (
	reID         = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	reGit        = regexp.MustCompile(`(?m)^(git.*)`)
	reRevision   = regexp.MustCompile(`(?m)^(.*ghc\S+)`)
	reWhitespace = regexp.MustCompile(`\s+`)
)

type CLI struct {
	Cmd              []string
	Dir              string
	SocketPath       string
	TestnetMagic     string
	TreasuryAddr     string
	TreasurySkeyFile string
	Debug            bool
}

//func New(dir, testnetMagic, treasuryAddr, treasuryKey string, base ...string) *CLI {
//	return &CLI{
//		Base:            base[0:len(base):len(base)],
//		Dir:             dir,
//		TestnetMagic:    testnetMagic,
//		TreasuryAddr:    treasuryAddr,
//		TreasuryKeyFile: treasuryKey,
//	}
//}

type Asset struct {
	AssetName   string `json:"asset_name,omitempty"`
	Description string `json:"description,omitempty"`
	Logo        string `json:"logo,omitempty"`
	Name        string `json:"name,omitempty"`
	PolicyId    string `json:"policy_id,omitempty"`
	Ticker      string `json:"ticker,omitempty"`
	Url         string `json:"url,omitempty"`
}

func (a Asset) ID() string {
	return fmt.Sprintf("%v.%v", a.PolicyId, a.AssetName)
}

type Tip struct {
	Block int32
	Epoch int32
	Era   string
	Hash  string
	Slot  int32
}

type Token struct {
	Asset    *Asset `json:"asset,omitempty"`
	Quantity string `json:"quantity,omitempty"`
}

type Utxo struct {
	Address   string  `json:"address,omitempty"`
	DatumHash string  `json:"datum_hash,omitempty"`
	Index     int32   `json:"index,omitempty"`
	Tokens    []Token `json:"tokens,omitempty"`
	Value     string  `json:"value,omitempty"`
}

func (u Utxo) TxIn() string {
	return fmt.Sprintf("%v#%v", u.Address, u.Index)
}

func (u Utxo) String() string {
	buf := bytes.NewBuffer(nil)

	if u.DatumHash != "" {
		fmt.Fprintf(buf, "%v@", u.DatumHash)
	}
	fmt.Fprintf(buf, "%v#%v+%v:lovelace", u.Address, u.Index, u.Value)
	for _, token := range u.Tokens {
		fmt.Fprintf(buf, "+%v:%v", token.Quantity, token.Asset.ID())
	}
	return buf.String()
}

type Utxos []Utxo

func (uu Utxos) Filter(includes ...func(utxo Utxo) bool) (filtered Utxos) {
loop:
	for _, utxo := range uu {
		for _, fn := range includes {
			if !fn(utxo) {
				continue loop
			}
		}
		filtered = append(filtered, utxo)
	}
	return filtered
}

func (uu Utxos) Find(address string, index int32) (Utxo, error) {
	for _, u := range uu {
		if u.Address == address && u.Index == index {
			return u, nil
		}
	}
	return Utxo{}, fmt.Errorf("unable to find utxo: %v#%v", address, index)
}

type Version struct {
	Revision string
	Git      string
}

func (c CLI) exec(args ...string) (*bytes.Buffer, error) {
	stmt := append(c.Cmd, args...)
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command(stmt[0], stmt[1:]...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("CARDANO_NODE_SOCKET_PATH=%v", c.SocketPath))
	cmd.Dir = c.Dir
	cmd.Stdout = buf
	cmd.Stderr = buf

	if c.Debug {
		cmd.Stdout = io.MultiWriter(buf, os.Stdout)
		cmd.Stderr = io.MultiWriter(buf, os.Stderr)
	}

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("exec failed: %v -> %v: %w", strings.Join(stmt, " "), buf.String(), err)
	}

	return buf, nil
}

func (c CLI) KeyHash(ctx context.Context, wallet string) (keyHash string, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("generated key hash",
			zap.Duration("elapsed", time.Now().Sub(begin).Round(time.Millisecond)),
			zap.String("wallet", wallet),
			zap.String("keyHash", keyHash),
			zap.Error(err),
		)
	}(time.Now())

	filename := paymentVerificationKeyFile(c.Dir, wallet)

	buf, err := c.exec("address", "key-hash", "--payment-verification-key-file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to generate key hash for wallet, %v: %w", wallet, err)
	}

	return strings.TrimSpace(buf.String()), nil
}

func (c CLI) MinFee(ctx context.Context, filename string, txIn, txOut, witnesses int32) (fee string, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("calculated min fee",
			zap.Duration("elapsed", time.Now().Sub(begin).Round(time.Millisecond)),
			zap.String("fee", fee),
			zap.Error(err),
		)
	}(time.Now())

	protocol, err := c.ProtocolParameters(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to calculate min fee: %w", err)
	}

	args := []string{
		"transaction", "calculate-min-fee",
		"--tx-body-file", filename,
		"--tx-in-count", strconv.FormatInt(int64(txIn), 10),
		"--tx-out-count", strconv.FormatInt(int64(txOut), 10),
		"--witness-count", strconv.FormatInt(int64(witnesses), 10),
		"--testnet-magic", c.TestnetMagic,
		"--protocol-params-file", protocol,
	}
	buf, err := c.exec(args...)
	if err != nil {
		return "", fmt.Errorf("unable to calculate min fee: %w", err)
	}

	fmt.Println("min-fees", buf)
	parts := reWhitespace.Split(buf.String(), -1)
	return parts[0], nil
}

func (c *CLI) NormalizeAddress(address string) (string, error) {
	if c == nil {
		return address, nil
	}

	if !reID.MatchString(address) {
		return address, nil
	}

	filename := filepath.Join(c.Dir, dirWallets, address+".addr")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("unable to read wallet: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

func (c CLI) PolicyID(ctx context.Context, filename string) (policyID string, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("generated policy id",
			zap.Duration("elapsed", time.Now().Sub(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())

	buf, err := c.exec("transaction", "policyid", "--script-file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to generate policy id: %w", err)
	}
	return strings.TrimSpace(buf.String()), nil
}

func (c CLI) ProtocolParameters(ctx context.Context) (filename string, err error) {
	defer func(begin time.Time) {
		zapctx.FromContext(ctx).Info("generated protocol parameters",
			zap.Duration("elapsed", time.Now().Sub(begin).Round(time.Millisecond)),
			zap.Error(err),
		)
	}(time.Now())

	filename = filepath.Join(c.Dir, "protocol.parameters")
	if _, err := os.Stat(filename); err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("unable to read protocol parameters: %w", err)
		}

		args := []string{
			"query", "protocol-parameters", "--testnet-magic", c.TestnetMagic, "--out-file", filename,
		}
		if _, err := c.exec(args...); err != nil {
			return "", fmt.Errorf("unable to read protocol parameters: %w", err)
		}
	}
	return filename, nil
}

func (c CLI) QueryTip() (*Tip, error) {
	buf, err := c.exec("query", "tip", "--testnet-magic", c.TestnetMagic, "--cardano-mode")
	if err != nil {
		return nil, fmt.Errorf("query tip failed: %w", err)
	}

	var tip Tip
	if err := json.Unmarshal(buf.Bytes(), &tip); err != nil {
		return nil, fmt.Errorf("query tip failed: %w", err)
	}

	return &tip, nil
}

// Utxos retrieves the list of utxos from cardano node.
//
//		TxHash                                 TxIx        Amount
//		--------------------------------------------------------------------------------------
//		111b3dc09d55e1708a22c866f697f358ccfe94dda61df8c0b9bca5b9081989ba     0        1000000000 lovelace + 1000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + 2000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
//		4d746439745c787087ac001c91270767ba0b4d10849fb6e8ab7c327b39f337f8     0        1000000000 lovelace + 1000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
//		52afd623b02712d5b37f582eae28ec31ee222547ff59e46d30015f2e3ab583f2     1        1000000000 lovelace + TxOutDatumHashNone
//		77c564a216621c883e628de4586254f7a4ea49d12daa560c5f22bd9886bf06b8     1        1000000000 lovelace + TxOutDatumHashNone
//		7a35e91a434183f6af77f7f2193efddb9661db024a6cdd391a4d6f7809b2627b     1        1000000000 lovelace + TxOutDatumHashNone
//		79ca8a27a1c05030d0c4290ddaf8c6e49ea8080aa867e7068a0b16e425b77d60     0        1000000000000 lovelace + TxOutDatumHashNone
//		84ab3e643b0bbd7856fdde0e723e50ba40008fc01b7b1ac03b9b861211e13d3d     0        5010000000000 lovelace + TxOutDatumHashNone
//		a1d4e06b6a0a5acdd0d3669b3a09ad195754fbde2387361935941df41474e5bd     0        5010000000000 lovelace + TxOutDatumHashNone
//
func (c CLI) Utxos(address string, excludes ...func(Utxo) bool) (utxos Utxos, err error) {
	address, err = c.NormalizeAddress(address)
	if err != nil {
		return nil, err
	}

	args := []string{"query", "utxo", "--testnet-magic", c.TestnetMagic, "--cardano-mode"}
	if address == "" {
		args = append(args, "--whole-utxo")
	} else {
		args = append(args, "--address", address)
	}

	buf, err := c.exec(args...)
	if err != nil {
		return nil, fmt.Errorf("query tip failed: %w", err)
	}

loop:
	for _, item := range ParseUtxos(buf) {
		utxo := item
		for _, fn := range excludes {
			if fn(utxo) {
				continue loop
			}
		}
		utxos = append(utxos, utxo)
	}

	return utxos, nil
}

func (c CLI) Version() (version Version, err error) {
	buf, err := c.exec("version")
	if err != nil {
		return Version{}, fmt.Errorf("query tip failed: %w", err)
	}

	if match := reGit.FindStringSubmatch(buf.String()); len(match) == 2 {
		version.Git = match[1]
	}
	if match := reRevision.FindStringSubmatch(buf.String()); len(match) == 2 {
		version.Revision = match[1]
	}

	return version, nil
}

func paymentVerificationKeyFile(dir, wallet string) string {
	return filepath.Join(dir, dirWallets, wallet+".vkey")
}

func HasToken(assetID string) func(utxo Utxo) bool {
	return func(utxo Utxo) bool {
		for _, token := range utxo.Tokens {
			if token.Asset.ID() == assetID {
				return true
			}
		}
		return false
	}
}
