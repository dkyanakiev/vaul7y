package vault

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dkyanakiev/vaulty/internal/models"
	"github.com/hashicorp/vault/api"
)

func (v *Vault) ListMounts() (map[string]*models.MountOutput, error) {

	apiMountList, err := v.vault.Sys().ListMounts()
	if err != nil {
		v.Logger.Warn().Err(err).Msg("Unable to access sys/mounts, attempting to use fallback method.\n")
		return v.listMountsFallback()
	}

	// Convert api.MountOutput to MountOutput
	mountList := make(map[string]*models.MountOutput)
	for k, v := range apiMountList {
		mountList[k] = toMount(v)
	}
	return mountList, nil

}

func (v *Vault) listMountsFallback() (map[string]*models.MountOutput, error) {

	resp, err := v.vault.Logical().ReadRaw("/sys/internal/ui/mounts")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret mounts: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.UiMountsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Convert models.MountOutput to api.MountOutput
	mountList := make(map[string]*models.MountOutput)
	for k, v := range response.Data.Secret {
		mountList[k] = v
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
