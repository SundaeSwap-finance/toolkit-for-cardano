package gql

func (r *Resolver) Version() (*VersionResolver, error) {
	version, err := r.config.CLI.Version()
	if err != nil {
		return nil, err
	}
	return newVersionResolver(r.config, version), nil
}
