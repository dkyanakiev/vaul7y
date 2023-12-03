package vault_test

import (
	"testing"

	"github.com/dkyanakiev/vaulty/internal/vault"
	"github.com/dkyanakiev/vaulty/internal/vault/vaultfakes"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	path := "testpath"

	fakeLogical := &vaultfakes.FakeLogical{}

	v := &vault.Vault{
		Logical: fakeLogical,
	}
	_, err := v.Logical.List(path)
	assert.NoError(t, err)

}
