package gql

func (r *Resolver) Wallets(args struct{ Query *string }) ([]string, error) {
	return r.config.CLI.FindAllWallets(StringValue(args.Query))
}
