package view

import (
	"sync"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/dkyanakiev/vaulty/layout"
	"github.com/dkyanakiev/vaulty/models"
	"github.com/dkyanakiev/vaulty/state"
	"github.com/rs/zerolog"
)

const (
	historySize = 15
)

type Client interface {
	UpdateSecretObject(mount string, path string, update bool, data map[string]interface{}) error
}

type Watcher interface {
	Subscribe(notify func(), topics ...string)
	Unsubscribe()
	SubscribeHandler(handler models.Handler, handle func(string, ...interface{}))
	SubscribeToPolicies(notify func())
	SubscribeToPoliciesACL(notify func())
	SubscribeToMounts(notify func())
	SubscribeToSecrets(selectedMount, selectedPath string, notify func())
	SubscribeToSecret(selectedMount, selectedPath string, notify func())
	UpdateMounts()
}

type View struct {
	Client  Client
	Watcher Watcher
	Layout  *layout.Layout

	history    *History
	state      *state.State
	logger     *zerolog.Logger
	components *Components
	mutex      sync.Mutex

	FilterText string ""

	draw chan struct{}
}

type Components struct {
	MountsTable    *component.MountsTable
	PolicyTable    *component.PolicyTable
	PolicyAclTable *component.PolicyAclTable
	SecretsTable   *component.SecretsTable
	SecretObjTable *component.SecretObjTable

	Commands  *component.Commands
	VaultInfo *component.VaultInfo
	Search    *component.SearchField
	Error     *component.Error
	Info      *component.Info
	Failure   *component.Info

	Selections   *component.Selections
	JumpToPolicy *component.JumpToPolicy
	Logo         *component.Logo
	Logger       *zerolog.Logger
}

func New(components *Components, watcher Watcher, client Client, state *state.State, logger *zerolog.Logger) *View {
	components.Search = component.NewSearchField("")

	return &View{
		Client:     client,
		Watcher:    watcher,
		state:      state,
		Layout:     layout.New(layout.Default, layout.EnableMouse),
		draw:       make(chan struct{}, 1),
		logger:     logger,
		components: components,
		history: &History{
			HistorySize: historySize,
			Logger:      logger,
		},
	}
}

func (v *View) Draw() {
	v.draw <- struct{}{}
}

// DrawLoop refreshes the screen when it receives a
// signal on the draw channel. This function should
// be run inside a goroutine as tview.Application.Draw()
// can deadlock when called from the main thread.
func (v *View) DrawLoop(stop chan struct{}) {
	for {
		select {
		case <-v.draw:
			v.Layout.Container.Draw()
		case <-stop:
			return
		}

	}
}

func (v *View) GoBack() {
	v.history.pop()
}

func (v *View) addToHistory(ns string, topic string, update func()) {
	v.history.push(func() {
		v.state.SelectedNamespace = ns
		// update()

		// v.components.Selections.Props.Rerender = update
		v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
			v.state.SelectedNamespace = text
			update()
		})
		// v.Watcher.Subscribe(topic, update)

		//index := getNamespaceNameIndex(ns, v.state.Namespaces)
		//v.state.Elements.DropDownNamespace.SetCurrentOption(index)
	})
}

func (v *View) viewSwitch() {
	v.resetSearch()
}

func (v *View) Search() {
	search := v.components.Search
	v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 1)
	search.Render()
	v.Layout.Container.SetFocus(search.InputField.Primitive())
}

func (v *View) resetSearch() {
	if v.state.Toggle.Search {
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		v.Layout.Footer.RemoveItem(v.components.Search.InputField.Primitive())
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.state.Toggle.Search = false
	}
}

// func (v *View) addToHistory(ns string, topic api.Topic, update func()) {
// 	v.history.push(func() {
// 		//v.state.SelectedNamespace = ns
// 		// update()

// 		// v.components.Selections.Props.Rerender = update
// 		v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
// 			v.state.SelectedNamespace = text
// 			update()
// 		})
// 		// v.Watcher.Subscribe(topic, update)

// 		index := getNamespaceNameIndex(ns, v.state.Namespaces)
// 		v.state.Elements.DropDownNamespace.SetCurrentOption(index)
// 	})
// }
