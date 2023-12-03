package vault

import (
	"fmt"

	"github.com/dkyanakiev/vaulty/internal/models"
	"github.com/hashicorp/vault/api"
)

func (v *Vault) ListMounts() (map[string]*models.MountOutput, error) {
	apiMountList, err := v.vault.Sys().ListMounts()
	if err != nil {
		return nil, fmt.Errorf("filed to retrieve secret mounts: %w", err)
	}
	// Convert api.MountOutput to models.MountOutput
	mountList := make(map[string]*models.MountOutput)
	for k, v := range apiMountList {
		mountList[k] = toMount(v)
	}

	return mountList, nil

}

func (v *Vault) AllMounts() (map[string]*models.MountOutput, error) {

	mounts, err := v.ListMounts()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve mounts: %w", err)
	}
	return mounts, nil
}

func toMount(m *api.MountOutput) *models.MountOutput {

	return &models.MountOutput{
		UUID:        m.UUID,
		Type:        m.Type,
		Description: m.Description,
		Accessor:    m.Accessor,
		Config: models.MountConfigOutput{
			//TODO:Fill out if needed
			DefaultLeaseTTL:   m.Config.DefaultLeaseTTL,
			ListingVisibility: m.Config.ListingVisibility,
		},
		Options:               m.Options,
		Local:                 m.Local,
		SealWrap:              m.SealWrap,
		ExternalEntropyAccess: m.ExternalEntropyAccess,
		PluginVersion:         m.PluginVersion,
		RunningVersion:        m.RunningVersion,
		RunningSha256:         m.RunningSha256,
		DeprecationStatus:     m.DeprecationStatus,
	}

}
