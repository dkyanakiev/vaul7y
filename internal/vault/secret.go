package vault

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/dkyanakiev/vaulty/internal/models"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

func (v *Vault) ListSecrets(path string) (*api.Secret, error) {

	secret, err := v.vault.Logical().List(fmt.Sprintf("%s/metadata", path))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to list for path:%s , secrets: %s", path, err))
	}

	// If the secret is wrapped, return the wrapped response.
	if secret != nil && secret.WrapInfo != nil && secret.WrapInfo.TTL != 0 {
		//TODO: Handle this usecase
		fmt.Println("Wrapped")
		//return OutputSecret(c.UI, secret)
	}

	return secret, nil

}

func (v *Vault) ListNestedSecrets(mount, path string) ([]models.SecretPath, error) {
	var secretPaths []models.SecretPath
	mountPath := fmt.Sprintf("%s/metadata/%s", mount, path)
	mountPath = sanitizePath(mountPath)
	secrets, err := v.vault.Logical().List(mountPath)

	v.Logger.Debug().Msg(fmt.Sprintf("Listing secrets for path: %s", mountPath))

	if err != nil {
		v.Logger.Err(err).Msgf("failed to list secrets: %s", err)
		return nil, fmt.Errorf("failed to list secrets: %s", err)
	}

	if secrets == nil {
		v.Logger.Err(err).Msgf("no secrets returned from the vault for path: %s", mountPath)
		return nil, errors.New("no secrets returned from the vault")
	}

	keys, ok := secrets.Data["keys"].([]interface{})
	if !ok {
		v.Logger.Err(err).Msgf("unexpected type for keys")
		return nil, errors.New("unexpected type for keys")
	}

	for _, key := range keys {
		keyStr, ok := key.(string)
		if !ok {
			return nil, errors.New("unexpected type for key")
		}

		isPath := strings.Contains(keyStr, "/")
		secretPath := models.SecretPath{
			PathName: keyStr,
			IsSecret: !isPath,
		}
		secretPaths = append(secretPaths, secretPath)
	}

	return secretPaths, nil
}

func (v *Vault) GetSecretInfo(mount, path string) (*api.Secret, error) {
	secretPath := fmt.Sprintf("%s/data/%s", mount, path)
	secretPath = sanitizePath(secretPath)
	secretData, err := v.vault.Logical().Read(secretPath)
	if err != nil {
		v.Logger.Err(err).Msgf("failed to read secret: %s", err)
		return nil, errors.New(fmt.Sprintf("Failed to read secret: %v", err))
	}

	if secretData == nil {
		v.Logger.Err(err).Msgf("no data found at %s", secretPath)
		return nil, errors.New(fmt.Sprintf("No data found at %s", secretPath))
	}
	//TODO: Add logging
	return secretData, nil
}

func (v *Vault) UpdateSecretObjectKV2(mount string, path string, patch bool, data map[string]interface{}) error {
	ctx := context.Background()

	data = prepareDataForWrite(data)
	v.Logger.Debug().Msg(fmt.Sprintf("Patch FLAG: %v", patch))
	secretPath := fmt.Sprintf("%s/data/%s", mount, path)

	secretPath = sanitizePath(secretPath)
	if patch {
		v.Logger.Debug().Msgf("Attempting to patch secret: %s", secretPath)
		// patchOption := func() (string, interface{}) {
		// 	return "method", "patch"
		// }

		v.Logger.Debug().Msgf("Data is: %v", data)

		_, err := v.PatchWithoutWrap(ctx, mount, path, data)
		if err != nil {

			if strings.Contains(err.Error(), "permission denied") {
				v.Logger.Err(err).Msgf("You do not have the necessary permissions to perform this operation")
				return errors.New("You do not have the necessary permissions to perform this operation")
			} else {
				v.Logger.Err(err).Msgf("Failed to patch secret: %v", err)
				return errors.New(fmt.Sprintf("Failed to patch secret: %v", err))
			}
		}

	} else {

		_, err := v.vault.Logical().WriteWithContext(ctx, secretPath, data)
		if err != nil {
			if strings.Contains(err.Error(), "permission denied") {
				v.Logger.Err(err).Msgf("You do not have the necessary permissions to perform this operation")
				return errors.New("You do not have the necessary permissions to perform this operation")
			} else {
				v.Logger.Err(err).Msgf("Failed to update secret: %v", err)
				return errors.New(fmt.Sprintf("Failed to update secret: %v", err))
			}
		}
	}

	if patch {
		v.Logger.Info().Msg("Secret updated successfully")
	} else {
		v.Logger.Info().Msg("Secret patched successfully")
	}

	return nil
}

func sanitizePath(p string) string {
	return path.Clean(p)
}

func prepareDataForWrite(data map[string]interface{}) map[string]interface{} {
	// Check if the 'data' key exists in the map
	if _, ok := data["data"]; !ok {
		// If the 'data' key does not exist, create it and move all existing keys into it
		newData := make(map[string]interface{})
		newData["data"] = data
		return newData
	}

	// If the 'data' key already exists, return the original map
	return data
}

func (v *Vault) getSecretVersion(mount, path string) (int, error) {
	var currentVersion int
	// Get the metadata of the secret
	secret, err := v.vault.Logical().List(fmt.Sprintf("%s/metadata/%s", mount, path))
	if err != nil {
		return currentVersion, errors.New(fmt.Sprintf("failed to read secret to check current version: %s", err))
	}

	if secret != nil && secret.Data != nil {
		if versions, ok := secret.Data["versions"].(map[string]interface{}); ok {
			// The versions are stored as a map where the keys are the version numbers
			// The highest version number is the current version
			for version := range versions {
				ver, err := strconv.Atoi(version)
				if err != nil {
					// Handle error
					return currentVersion, errors.New(fmt.Sprintf("failed to read secret to check current version: %s", err))

				}
				if ver > currentVersion {
					currentVersion = ver
				}
			}
		}
	}

	return currentVersion, nil
}

func (v *Vault) getCurrentVersion(mount, path string) (int, error) {
	// Get the metadata of the secret
	secret, err := v.vault.Logical().Read(fmt.Sprintf("%s/metadata/%s", mount, path))
	if err != nil {
		return 0, fmt.Errorf("failed to read secret metadata: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return 0, fmt.Errorf("no data found in secret metadata")
	}

	// Get the versions map from the metadata
	versions, ok := secret.Data["versions"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("no versions found in secret metadata")
	}

	// The versions are stored as a map where the keys are the version numbers
	// The highest version number is the current version
	var currentVersion int
	for version := range versions {
		v, err := strconv.Atoi(version)
		if err != nil {
			// Handle error
		}
		if v > currentVersion {
			currentVersion = v
		}
	}

	return currentVersion, nil
}

func (v *Vault) mergePatchWithoutWrap(ctx context.Context, mountPath string, secretPath string, newData map[string]interface{}, opts ...KVOption) (*api.KVSecret, error) {
	pathToMergePatch := fmt.Sprintf("%s/data/%s", mountPath, secretPath)

	// take any other additional options provided
	// and pass them along to the patch request

	options := make(map[string]interface{})
	for _, opt := range opts {
		k, v := opt()
		options[k] = v
	}
	if len(opts) > 0 {
		newData["options"] = options
	}

	secret, err := v.vault.Logical().JSONMergePatch(ctx, pathToMergePatch, newData)
	if err != nil {
		var re *api.ResponseError

		if errors.As(err, &re) {
			switch re.StatusCode {
			// 403
			case http.StatusForbidden:
				return nil, fmt.Errorf("received 403 from Vault server; please ensure that token's policy has \"patch\" capability: %w", err)

			// 404
			case http.StatusNotFound:
				return nil, fmt.Errorf("%w: performing merge patch to %s", api.ErrSecretNotFound, pathToMergePatch)

			// 405
			case http.StatusMethodNotAllowed:
				// If it's a 405, that probably means the server is running a pre-1.9
				// Vault version that doesn't support the HTTP PATCH method.
				// Fall back to the old way of doing it.
				return v.readThenWrite(ctx, mountPath, secretPath, newData)
			}
		}

		return nil, fmt.Errorf("error performing merge patch to %s: %w", pathToMergePatch, err)
	}

	metadata, err := extractVersionMetadata(secret)
	if err != nil {
		return nil, fmt.Errorf("secret was written successfully, but unable to view version metadata from response: %w", err)
	}

	kvSecret := &api.KVSecret{
		Data:            nil, // secret.Data in this case is the metadata
		VersionMetadata: metadata,
		Raw:             secret,
	}

	kvSecret.CustomMetadata = extractCustomMetadata(secret)

	return kvSecret, nil
}

func (v *Vault) readThenWrite(ctx context.Context, mountPath string, secretPath string, newData map[string]interface{}) (*api.KVSecret, error) {
	// First, read the secret.
	existingVersion, err := v.vault.KVv2(mountPath).Get(ctx, secretPath)
	if err != nil {
		return nil, fmt.Errorf("error reading secret as part of read-then-write patch operation: %w", err)
	}

	// Make sure the secret already exists
	if existingVersion == nil || existingVersion.Data == nil {
		return nil, fmt.Errorf("%w: at %s as part of read-then-write patch operation", api.ErrSecretNotFound, secretPath)
	}

	// Verify existing secret has metadata
	if existingVersion.VersionMetadata == nil {
		return nil, fmt.Errorf("no metadata found at %s; patch can only be used on existing data", secretPath)
	}

	// Copy new data over with existing data
	combinedData := existingVersion.Data
	for k, v := range newData {
		combinedData[k] = v
	}

	updatedSecret, err := v.vault.KVv2(mountPath).Put(ctx, secretPath, combinedData, api.WithCheckAndSet(existingVersion.VersionMetadata.Version))
	if err != nil {
		return nil, fmt.Errorf("error writing secret to %s: %w", secretPath, err)
	}

	return updatedSecret, nil
}

func extractVersionMetadata(secret *api.Secret) (*api.KVVersionMetadata, error) {
	var metadata *api.KVVersionMetadata

	if secret.Data == nil {
		return nil, nil
	}

	// Logical Writes return the metadata directly, Reads return it nested inside the "metadata" key
	var metadataMap map[string]interface{}
	metadataInterface, ok := secret.Data["metadata"]
	if ok {
		metadataMap, ok = metadataInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected type for 'metadata' element: %T (%#v)", metadataInterface, metadataInterface)
		}
	} else {
		metadataMap = secret.Data
	}

	// deletion_time usually comes in as an empty string which can't be
	// processed as time.RFC3339, so we reset it to a convertible value
	if metadataMap["deletion_time"] == "" {
		metadataMap["deletion_time"] = time.Time{}
	}

	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeHookFunc(time.RFC3339),
		Result:     &metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("error setting up decoder for API response: %w", err)
	}

	err = d.Decode(metadataMap)
	if err != nil {
		return nil, fmt.Errorf("error decoding metadata from API response into VersionMetadata: %w", err)
	}

	return metadata, nil
}

func extractCustomMetadata(secret *api.Secret) map[string]interface{} {
	// Logical Writes return the metadata directly, Reads return it nested inside the "metadata" key
	customMetadataInterface, ok := secret.Data["custom_metadata"]
	if !ok {
		metadataInterface := secret.Data["metadata"]
		metadataMap, ok := metadataInterface.(map[string]interface{})
		if !ok {
			return nil
		}
		customMetadataInterface = metadataMap["custom_metadata"]
	}

	cm, ok := customMetadataInterface.(map[string]interface{})
	if !ok {
		return nil
	}

	return cm
}

func (v *Vault) PatchWithoutWrap(ctx context.Context, mountPath string, secretPath string, newData map[string]interface{}) (*api.KVSecret, error) {

	kvs, err := v.mergePatchWithoutWrap(ctx, mountPath, secretPath, newData)

	if err != nil {
		return nil, fmt.Errorf("unable to perform patch: %w", err)
	}
	if kvs == nil {
		return nil, fmt.Errorf("no secret was written to %s", secretPath)
	}

	return kvs, nil
}

func (v *Vault) CreateNewSecret(mount string, path string) error {
	secretPath := fmt.Sprintf("%s/data/%s", mount, path)
	secretPath = sanitizePath(secretPath)

	data := map[string]interface{}{
		"data": make(map[string]interface{}),
	}

	secret, err := v.vault.Logical().Write(secretPath, data)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	if secret == nil {
		return fmt.Errorf("no secret was written to %s", secretPath)
	}

	return nil
}

func (v *Vault) ChangeNamespace(ns string) []string {
	v.Logger.Debug().Msgf("Changing namespace to: %v", ns)
	v.vault.SetNamespace(ns)
	list, err := v.ListNamespaces()
	if err != nil {
		v.Logger.Err(err).Msgf("Failed to list namespaces: %v", err)
	}
	v.Logger.Debug().Msgf("New available namespaces are: %v", list)
	return list
}
