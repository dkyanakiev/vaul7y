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

func Default(v *Vault, log *zerolog.Logger) error {
	//ctx := context.Background()
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}
	//TODO: This needs to be changed if we want to support namespaces
	client.SetNamespace("admin")
	v.vault = client
	v.Client = client
	//KV1 is not setup as its adviced against
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
