package vault_test

import (
	"testing"

	"github.com/dkyanakiev/vaulty/internal/vault"
	"github.com/dkyanakiev/vaulty/internal/vault/vaultfakes"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

func TestListMounts(t *testing.T) {

	fakeSys := &vaultfakes.FakeSys{}
	fakeSys.ListMountsReturns(map[string]*api.MountOutput{}, nil)

	v := &vault.Vault{
		Sys: fakeSys,
	}

	_, err := v.Sys.ListMounts()
	assert.NoError(t, err)

}
