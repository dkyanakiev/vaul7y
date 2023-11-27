package vault

import (
	"context"

	"github.com/hashicorp/vault/api"
)

func (v *Vault) Get(ctx context.Context, path string) (*api.KVSecret, error) {
	secret, err := v.KV2.Get(ctx, path)
	if err != nil {
		v.Logger.Panicln("filed to retrieve secret: %w", err)
	}
	return secret, nil
}

func (v *Vault) GetMetadata(ctx context.Context, path string) (*api.KVMetadata, error) {
	secret, err := v.KV2.GetMetadata(ctx, path)
	if err != nil {
		v.Logger.Println("filed to retrieve secret: %w", err)
	}
	return secret, nil
}
