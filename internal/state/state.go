package state

import (
	"github.com/dkyanakiev/vaulty/internal/models"
	"github.com/hashicorp/vault/api"
	"github.com/rivo/tview"
)

type State struct {
	VaultAddress string
	VaultVersion string
	Mounts       map[string]*models.MountOutput
	SecretsData  []models.SecretPath
	KV2          []models.KVSecret
	// Central/Root ns for the Vault instance
	RootNamespace    string
	DefaultNamespace string
	Namespaces       []string
	// Current ns for the Vault instance
	SelectedNamespace  string
	SelectedMount      string
	SelectedPath       string
	SelectedObject     string
	SelectedPolicyName string
	SelectedSecret     *api.Secret
	SelectedSecretMeta *models.Metadata
	PolicyList         []string
	PolicyACL          string
	NewSecretName      string
	Enterprise         bool

	Elements *Elements
	Toggle   *Toggle
	Filter   *Filter
	Version  string
}

type Toggle struct {
	Search       bool
	JumpToPolicy bool
	TextInput    bool
}

type Filter struct {
	Object    string
	Policy    string
	Namespace string
}

type Elements struct {
	DropDownNamespace *tview.DropDown
	TableMain         *tview.Table
	TextMain          *tview.TextView
}

func New() *State {
	return &State{
		Elements: &Elements{},
		Toggle:   &Toggle{},
		Filter:   &Filter{},
	}
}
