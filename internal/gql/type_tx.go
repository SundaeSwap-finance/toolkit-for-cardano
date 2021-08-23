package gql

type TxResolver struct {
	body string
	id   string
	raw  []byte
}

func (t *TxResolver) Body() string { return t.body }
func (t *TxResolver) Id() string   { return t.id }
