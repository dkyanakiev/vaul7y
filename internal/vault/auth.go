package vault

import (
	"fmt"

	"github.com/dkyanakiev/vaul7y/internal/models"
)

func (v *Vault) ListAuthMethods() (map[string]*models.AuthMethod, error) {
	auths, err := v.vault.Sys().ListAuth()
	if err != nil {
		return nil, fmt.Errorf("failed to list auth methods: %w", err)
	}

	result := make(map[string]*models.AuthMethod, len(auths))
	for path, mount := range auths {
		result[path] = &models.AuthMethod{
			Type:        mount.Type,
			Description: mount.Description,
			Accessor:    mount.Accessor,
			Local:       mount.Local,
		}
	}
	return result, nil
}
