package gql

import (
	"github.com/SundaeSwap-finance/cardano-toolkit/internal/cardano"
)

type VersionResolver struct {
	config  Config
	version cardano.Version
}

func newVersionResolver(config Config, version cardano.Version) *VersionResolver {
	return &VersionResolver{
		config:  config,
		version: version,
	}
}

func (v *VersionResolver) Built() string    { return v.config.Built }
func (v *VersionResolver) Git() string      { return v.version.Git }
func (v *VersionResolver) Revision() string { return v.version.Revision }
func (v *VersionResolver) Version() string  { return v.config.Version }
