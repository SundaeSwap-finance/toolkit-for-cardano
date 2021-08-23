package gql

func (r *Resolver) Tip() (*TipResolver, error) {
	tip, err := r.config.CLI.QueryTip()
	if err != nil {
		return nil, err
	}

	return &TipResolver{tip: tip}, nil
}
