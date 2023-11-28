package main

import (
	"log"
	"time"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/dkyanakiev/vaulty/config"
	"github.com/dkyanakiev/vaulty/state"
	"github.com/dkyanakiev/vaulty/vault"
	"github.com/dkyanakiev/vaulty/view"
	"github.com/dkyanakiev/vaulty/watcher"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var refreshIntervalDefault = time.Second * 5

func main() {

	logFile, logger := config.SetupLogger()
	defer logFile.Close()
	tview.Styles.PrimitiveBackgroundColor = tcell.NewRGBColor(40, 44, 48)

	vaultClient, err := vault.New(func(v *vault.Vault) error {
		return vault.Default(v, logger)
	})
	state := initializeState(vaultClient)
	commands := component.NewCommands()
	vaultInfo := component.NewVaultInfo()
	mounts := component.NewMountsTable()
	policies := component.NewPolicyTable()
	policyAcl := component.NewPolicyAclTable()
	secrets := component.NewSecretsTable()
	secretObj := component.NewSecretObjTable()
	logo := component.NewLogo()
	info := component.NewInfo()
	failure := component.NewInfo()
	errorComp := component.NewError()
	components := &view.Components{
		VaultInfo:      vaultInfo,
		Commands:       commands,
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
	}
	watcher := watcher.NewWatcher(state, vaultClient, refreshIntervalDefault, logger)
	view := view.New(components, watcher, vaultClient, state, logger)

	view.Init("0.0.1")
	err = view.Layout.Container.Run()
	if err != nil {
		log.Fatal("cannot initialize view.")
	}

}

func initializeState(client *vault.Vault) *state.State {
	state := state.New()
	addr := client.Address()
	state.VaultAddress = addr
	state.Namespace = "default"

	return state
}

// // LOOK AT LATER
// func main() {
// 	vaultClient, _ := vault.New(vault.Default)
// 	//ctx := context.TODO()
// 	// mounts, _ := vaultClient.Sys.ListMounts()

// 	secret, _ := vaultClient.ListSecrets("kv0FF76557")
// 	log.Println(secret)

// 	secrets, _ := vaultClient.ListNestedSecrets("kv0FF76557", "")
// 	//secrets, err := vaultClient.Logical.List("randomkv/metadata/test/one")

// 	for _, value := range secrets {
// 		fmt.Printf("Key: %s\n", value.PathName)
// 		fmt.Printf("IsSecret: %t\n", value.IsSecret)
// 	}
// 	// val, err := vaultClient.KV2.Get(ctx, "path")
// 	// fmt.Println(val)
// 	// fmt.Println(err)

// 	// secretClient, err := vaultClient.Logical.List("credentials/metadata/")
// 	// if err != nil {
// 	// 	// TODO
// 	// 	fmt.Println(err)
// 	// }
// 	// vault.DataIterator(secretClient.Data["keys"])
// }
