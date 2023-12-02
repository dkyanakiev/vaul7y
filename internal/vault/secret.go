package vault

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/dkyanakiev/vaulty/internal/models"
	"github.com/hashicorp/vault/api"
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

func (v *Vault) UpdateSecretObject(mount string, path string, update bool, data map[string]interface{}) error {

	secretPath := fmt.Sprintf("%s/data/%s", mount, path)
	secretPath = sanitizePath(secretPath)
	if !update {
		data["options"] = map[string]interface{}{
			"cas": 0, // Use 'cas' (Check-And-Set) to patch the secret
		}
	}
	v.Logger.Debug().Msg(fmt.Sprintf("Writing secret to %s", secretPath))
	_, err := v.vault.Logical().Write(secretPath, data)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			v.Logger.Err(err).Msgf("You do not have the necessary permissions to perform this operation")
			return errors.New("You do not have the necessary permissions to perform this operation")
		} else {
			v.Logger.Err(err).Msgf("Failed to write secret: %v", err)
			return errors.New(fmt.Sprintf("Failed to write secret: %v", err))
		}
	}

	if update {
		v.Logger.Info().Msg("Secret updated successfully")
	} else {
		v.Logger.Info().Msg("Secret patched successfully")
	}

	return nil
}

func sanitizePath(p string) string {
	return path.Clean(p)
}
