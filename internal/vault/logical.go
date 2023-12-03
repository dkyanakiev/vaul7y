package vault

import (
	"context"
	"io"
	"net/http"

	"github.com/hashicorp/vault/api"
)

func (v *Vault) List(path string) (*api.Secret, error) {
	return v.ListWithContext(context.Background(), path)
}

func (v *Vault) ListWithContext(ctx context.Context, path string) (*api.Secret, error) {
	// ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	// defer cancelFunc()

	r := v.vault.NewRequest("LIST", "/v1/"+path)

	// Set this for broader compatibility, but we use LIST above to be able to
	// handle the wrapping lookup function
	r.Method = http.MethodGet
	r.Params.Set("list", "true")

	// resp, err := v.vault.RawRequestWithContext(ctx, r)
	resp, err := v.vault.Logical().ReadRawWithContext(ctx, path)
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, parseErr
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}
