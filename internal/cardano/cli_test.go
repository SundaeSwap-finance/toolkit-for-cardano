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
	"encoding/json"
	"os"
	"testing"

	"github.com/tj/assert"
)

func TestReUtxo(t *testing.T) {
	testCases := map[string]struct {
		Input string
		Want  string
	}{
		"complex": {
			Input: `
 + 1000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + 2000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
 + 3000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
 + TxOutDatumHashNone
`,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			matches := reTokens.FindAllStringSubmatch(tc.Input, -1)
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			encoder.Encode(matches)
		})
	}
}

func TestNormalize(t *testing.T) {
	cli := &CLI{Dir: "testdata"}
	got, err := cli.NormalizeAddress("sample")
	assert.Nil(t, err)
	assert.Equal(t, "blah", got)

	got, err = cli.NormalizeAddress("foo")
	assert.Nil(t, err)
	assert.Equal(t, "foo", got)
}
