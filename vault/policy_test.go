package vault_test

import (
	"testing"

	"github.com/dkyanakiev/vaulty/vault"
	"github.com/dkyanakiev/vaulty/vault/vaultfakes"
	"github.com/stretchr/testify/assert"
)

func TestAllPolicies(t *testing.T) {

	fakeSys := &vaultfakes.FakeSys{}

	v := &vault.Vault{
		Sys: fakeSys,
	}

	_, err := v.AllPolicies()

	assert.NoError(t, err)

}

func TestGetPolicyInfo(t *testing.T) {

	fakeSys := &vaultfakes.FakeSys{}

	v := &vault.Vault{
		Sys: fakeSys,
	}

	_, err := v.GetPolicyInfo("test")

	assert.NoError(t, err)
}
