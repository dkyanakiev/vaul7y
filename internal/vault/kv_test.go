package vault_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/dkyanakiev/vaul7y/internal/vault"
	"github.com/dkyanakiev/vaul7y/internal/vault/vaultfakes"
	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func discardLogger() *zerolog.Logger {
	l := zerolog.New(io.Discard)
	return &l
}

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

func TestGet_Error(t *testing.T) {
	ctx := context.Background()
	path := "testpath"

	fakeKV2 := &vaultfakes.FakeKV2{}
	vaultErr := errors.New("permission denied")
	fakeKV2.GetReturns(nil, vaultErr)

	v := &vault.Vault{
		KV2:    fakeKV2,
		Logger: discardLogger(),
	}

	secret, err := v.Get(ctx, path)

	assert.ErrorIs(t, err, vaultErr)
	assert.Nil(t, secret)
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

func TestGetMetadata_Error(t *testing.T) {
	ctx := context.Background()
	path := "testpath"

	fakeKV2 := &vaultfakes.FakeKV2{}
	vaultErr := errors.New("secret not found")
	fakeKV2.GetMetadataReturns(nil, vaultErr)

	v := &vault.Vault{
		KV2:    fakeKV2,
		Logger: discardLogger(),
	}

	meta, err := v.GetMetadata(ctx, path)

	assert.ErrorIs(t, err, vaultErr)
	assert.Nil(t, meta)
}
