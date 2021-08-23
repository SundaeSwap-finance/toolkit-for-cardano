package cardano

import (
	"encoding/json"
	"os"
	"testing"
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
