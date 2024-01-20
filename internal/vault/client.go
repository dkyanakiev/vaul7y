package vault

import (
	"context"
	"strings"

	"github.com/dkyanakiev/vaulty/internal/models"
	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog"
)

//go:generate counterfeiter . Client
type Client interface {
	Address() string
}

type Vault struct {
	vault    *api.Client
	Client   Client
	KV2      KV2
	Sys      Sys
	NsClient NamespaceClient
	Logical  Logical
	Secret   Secret
	Logger   *zerolog.Logger
}

//go:generate counterfeiter . Logical
type Logical interface {
	List(path string) (*api.Secret, error)
}

type Secret interface {
	//ListSecrets(string) (*api.Secret, error)
}

//go:generate counterfeiter . Sys
type Sys interface {
	ListMounts() (map[string]*api.MountOutput, error)
	ListPolicies() ([]string, error)
	GetPolicy(name string) (string, error)
	Health() (*api.HealthResponse, error)
	// ListNamespaces() ([]models.Namespace, error)
	//ListMounts() ([]*api.Sys, error)
}

//go:generate counterfeiter . KV2
type KV2 interface {
	Get(context.Context, string) (*api.KVSecret, error)
	GetMetadata(context.Context, string) (*api.KVMetadata, error)
	Patch(context.Context, string, map[string]interface{}, ...KVOption) (*api.KVSecret, error)
	Put(context.Context, string, map[string]interface{}, ...KVOption) (*api.KVSecret, error)

	// GetVersion(context.Context, string, int) (*api.KVSecret, error)
	// GetVersionsAsList(context.Context, string) ([]*api.KVVersionMetadata, error)
}

type NamespaceClient interface {
	List() ([]*models.Namespace, error)
}

type KVOption func() (key string, value interface{})

func New(opts ...func(*Vault) error) (*Vault, error) {
	vault := Vault{}
	for _, opt := range opts {
		err := opt(&vault)
		if err != nil {
			return nil, err
		}
	}

	return &vault, nil
}

func Default(v *Vault, log *zerolog.Logger, vault_addr string, vault_token string, vault_ns string) error {
	cfg := api.DefaultConfig()
	cfg.Address = vault_addr

	client, err := api.NewClient(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Vault client")
		return err
	}

	client.SetToken(vault_token)

	// Check if the client is successfully created by making a request to Vault
	health, err := client.Sys().Health()
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to Vault")
		return err
	}
	// Check for enterprise version and set namespace
	if strings.Contains(health.Version, "ent") {
		client.SetNamespace(vault_ns)
	}

	log.Debug().Msg("Vault client successfully created and connected")

	v.vault = client
	v.Client = client
	v.Sys = client.Sys()
	v.Logical = client.Logical()
	v.Logger = log

	return nil
}

func (v *Vault) Address() string {
	return v.Client.Address()
}

func (v *Vault) Version() (string, error) {
	health, err := v.Sys.Health()
	if err != nil {
		return "", err
	}
	return health.Version, nil
}

func (v *Vault) ListNamespaces() ([]string, error) {
	namespaces := []string{}
	secret, err := v.vault.Logical().List("sys/namespaces")
	if err != nil {
		return nil, err
	}

	if secret != nil {
		keys, ok := secret.Data["keys"].([]interface{})
		if !ok || len(keys) == 0 {
			return namespaces, nil
		}
		for _, namespace := range keys {
			trimmedNamespace := strings.TrimSuffix(namespace.(string), "/")
			namespaces = append(namespaces, trimmedNamespace)
		}
	} else {
		return namespaces, nil
	}
	return namespaces, nil
}
