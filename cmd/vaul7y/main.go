package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dkyanakiev/vaulty/internal/config"
	"github.com/dkyanakiev/vaulty/internal/state"
	"github.com/dkyanakiev/vaulty/internal/vault"
	"github.com/dkyanakiev/vaulty/internal/watcher"
	"github.com/dkyanakiev/vaulty/tui/component"
	"github.com/dkyanakiev/vaulty/tui/view"
	"github.com/gdamore/tcell/v2"
	"github.com/jessevdk/go-flags"
	"github.com/rivo/tview"
)

var refreshIntervalDefault = time.Second * 30
var version = "0.0.7"

type options struct {
	Version bool `short:"v" long:"version" description:"Show Damon version"`
}

func main() {

	var opts options
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println("vaul7y", version)
		os.Exit(0)
	}

	// Check for required Vault env vars
	checkForVaultAddress()

	logFile, logger := config.SetupLogger()
	defer logFile.Close()
	tview.Styles.PrimitiveBackgroundColor = tcell.NewRGBColor(40, 44, 48)

	vaultClient, err := vault.New(func(v *vault.Vault) error {
		return vault.Default(v, logger)
	})

	state := initializeState(vaultClient)
	toggles := component.NewTogglesInfo()
	selections := component.NewSelections(state)
	namespaces := component.NewNamespaceTable()
	commands := component.NewCommands()
	vaultInfo := component.NewVaultInfo()
	mounts := component.NewMountsTable()
	policies := component.NewPolicyTable()
	policyAcl := component.NewPolicyAclTable()
	secrets := component.NewSecretsTable()
	secretObj := component.NewSecretObjTable()
	logo := component.NewLogo(version)
	info := component.NewInfo()
	failure := component.NewInfo()
	errorComp := component.NewError()
	components := &view.Components{
		VaultInfo:      vaultInfo,
		Commands:       commands,
		Selections:     selections,
		NamespaceTable: namespaces,
		MountsTable:    mounts,
		PolicyTable:    policies,
		PolicyAclTable: policyAcl,
		SecretsTable:   secrets,
		SecretObjTable: secretObj,
		Info:           info,
		Error:          errorComp,
		Failure:        failure,
		Logo:           logo,
		Logger:         logger,
		TogglesInfo:    toggles,
	}
	watcher := watcher.NewWatcher(state, vaultClient, refreshIntervalDefault, logger)
	view := view.New(components, watcher, vaultClient, state, logger)
	view.Init(version)

	//view.Init("0.0.1")
	err = view.Layout.Container.Run()
	if err != nil {
		log.Fatal("cannot initialize view.")
	}

}

func initializeState(client *vault.Vault) *state.State {
	state := state.New()
	addr := client.Address()
	version, _ := client.Version()
	state.VaultAddress = addr
	state.VaultVersion = version
	//TODO
	state.RootNamespace = "admin"
	state.Namespaces, _ = client.ListNamespaces()
	//	state.Namespace = "default"

	return state
}

func checkForVaultAddress() {
	if os.Getenv("VAULT_ADDR") == "" {
		fmt.Println("VAULT_ADDR is not set. Please set it and try again.")
		os.Exit(1)
	}

	if os.Getenv("VAULT_TOKEN") == "" {
		fmt.Println("VAULT_TOKEN is not set. Please set it and try again.")
		os.Exit(1)
	}

}
