package cardano

import (
	"bytes"
	"regexp"
	"strconv"
)

var (
	reUtxo   = regexp.MustCompile(`(?m)^([a-z0-9]+)\s+(\d+)\s+(\d+)\s+lovelace(.*)$`)
	reTokens = regexp.MustCompile(`\+\s*(\d+)\s+([a-z0-9]+).(\S+)`)
	reScript = regexp.MustCompile(`ScriptDataInAlonzoEra\s+"([^"]+)"`)
)

func ParseUtxos(buf *bytes.Buffer) []Utxo {
	var (
		matches = reUtxo.FindAllStringSubmatch(buf.String(), -1)
		utxos   []Utxo
	)

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		index, _ := strconv.ParseInt(match[2], 10, 32)
		utxo := Utxo{
			Address: match[1],
			Index:   int32(index),
			Value:   match[3],
		}

		if len(match) >= 5 {
			if extra := match[4]; extra != "" {
				if ss := reScript.FindStringSubmatch(extra); len(ss) == 2 {
					utxo.DatumHash = ss[1]
				}

				tokens := reTokens.FindAllStringSubmatch(extra, -1)
				for _, token := range tokens {
					if len(token) < 4 {
						continue
					}
					utxo.Tokens = append(utxo.Tokens, Token{
						Asset: &Asset{
							AssetName: token[3],
							PolicyId:  token[2],
						},
						Quantity: token[1],
					})
				}
			}
		}

		utxos = append(utxos, utxo)
	}
	return utxos
}
