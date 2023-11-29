package vault

import (
	"context"

	"github.com/hashicorp/vault/api"
)

func (v *Vault) Get(ctx context.Context, path string) (*api.KVSecret, error) {
	secret, err := v.KV2.Get(ctx, path)
	if err != nil {
		v.Logger.Err(err).Msgf("filed to retrieve secret: %s", err)
	}
	return secret, nil
}

func (v *Vault) GetMetadata(ctx context.Context, path string) (*api.KVMetadata, error) {
	secret, err := v.KV2.GetMetadata(ctx, path)
	if err != nil {
		v.Logger.Err(err).Msgf("filed to retrieve secret: %s", err)
	}
	return secret, nil
}
