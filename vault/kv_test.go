package vault_test

import (
	"context"
	"testing"

	"github.com/dkyanakiev/vaulty/vault"
	"github.com/dkyanakiev/vaulty/vault/vaultfakes"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	ctx := context.Background()
	path := "testpath"

	fakeKV2 := &vaultfakes.FakeKV2{}

	fakeKV2.GetReturns(&api.KVSecret{}, nil)

	v := &vault.Vault{
		KV2: fakeKV2,
	}

	secret, err := v.Get(ctx, path)

	assert.NoError(t, err)
	assert.NotNil(t, secret)
	fakeKV2.Get(ctx, path)

}

func TestGetMetadata(t *testing.T) {
	ctx := context.Background()
	path := "testpath"

	fakeKV2 := &vaultfakes.FakeKV2{}

	fakeKV2.GetMetadataReturns(&api.KVMetadata{}, nil)

	v := &vault.Vault{
		KV2: fakeKV2,
	}

	secret, err := v.GetMetadata(ctx, path)

	assert.NoError(t, err)
	assert.NotNil(t, secret)
	fakeKV2.GetMetadata(ctx, path)
}
