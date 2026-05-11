package state

import (
	"sync"

	"github.com/dkyanakiev/vaul7y/internal/models"
	"github.com/hashicorp/vault/api"
	"github.com/rivo/tview"
)

type State struct {
	mu sync.RWMutex
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

	AuthMethods   map[string]*models.AuthMethod
	TokenTTL      int64
	TokenPolicies []string
	SealStatus    string
	ClusterName   string

	Elements *Elements
	Toggle   *Toggle
	Filter   *Filter
	Version  string
}

type Toggle struct {
	Search       bool
	JumpToPolicy bool
	JumpToPath   bool
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

func (s *State) Lock()    { s.mu.Lock() }
func (s *State) Unlock()  { s.mu.Unlock() }
func (s *State) RLock()   { s.mu.RLock() }
func (s *State) RUnlock() { s.mu.RUnlock() }
