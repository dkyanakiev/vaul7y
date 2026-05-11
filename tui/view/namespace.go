package view

import (
	"fmt"

	"github.com/dkyanakiev/vaul7y/tui/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) Namespaces() {
	v.viewSwitch()
	v.logger.Debug().Msg("view: Namespaces")
	v.Layout.Body.SetTitle("Vault Nmespaces")
	v.Layout.Container.SetFocus(v.components.NamespaceTable.Table.Primitive())

	v.state.Elements.TableMain = v.components.NamespaceTable.Table.Primitive().(*tview.Table)
	v.components.NamespaceTable.Logger = v.logger
	v.components.Commands.Update(component.NamespaceObjectCommands)
	v.Layout.Container.SetInputCapture(v.InputNamespaces)

	update := func() {
		v.state.RLock()
		namespaces := v.state.Namespaces
		v.state.RUnlock()

		v.logger.Debug().Msgf("Current ns list: %v", namespaces)
		v.components.NamespaceTable.Props.Data = namespaces
		v.components.NamespaceTable.Render()
		v.Draw()
		v.components.NamespaceTable.Table.ScrollToTop()
	}

	v.components.Search.Props.ChangedFunc = func(text string) {
		//v.state.Filter.Namespace = text
		update()
	}

	v.Watcher.SubscribeToNamespaces(func() { v.Layout.Container.QueueUpdateDraw(update) })

	update()

	// v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
	// 	v.state.SelectedNamespace = text
	// 	v.Namespaces()
	// })

	// v.addToHistory(v.state.SelectedNamespace, models.TopicNamespace, v.Namespaces)
}

func (v *View) inputNamespaces(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyEsc:
		//v.GoBack()
	case tcell.KeyCtrlD:
		v.state.RLock()
		defaultNs := v.state.DefaultNamespace
		v.state.RUnlock()
		v.logger.Debug().Msgf("Going back to default namespace: %v", defaultNs)
		v.state.Lock()
		v.state.SelectedNamespace = defaultNs
		v.state.Unlock()
		v.components.TogglesInfo.Props.Namespace = defaultNs
		v.components.TogglesInfo.Render()
		v.Mounts()
		return nil
	case tcell.KeyCtrlW:
		v.state.RLock()
		rootNs := v.state.RootNamespace
		v.state.RUnlock()
		v.logger.Debug().Msgf("Going back to root namespace : %v", rootNs)
		v.state.Lock()
		v.state.SelectedNamespace = rootNs
		v.state.Unlock()
		v.components.TogglesInfo.Props.Namespace = rootNs
		v.components.TogglesInfo.Render()
		v.Mounts()
		return nil
	case tcell.KeyEnter:
		selectdNs := v.components.NamespaceTable.GetIDForSelection()
		v.state.RLock()
		currentNs := v.state.SelectedNamespace
		v.state.RUnlock()
		v.logger.Debug().Msgf("Selected namespace is: %v", selectdNs)
		newNs := fmt.Sprintf("%s/%s", currentNs, selectdNs)
		v.logger.Debug().Msgf("Changing namespace to: %s", newNs)
		v.state.Lock()
		v.state.SelectedNamespace = newNs
		v.state.Unlock()
		v.components.TogglesInfo.Props.Namespace = newNs
		v.components.TogglesInfo.Render()
		v.Mounts()
		return nil
	}

	return event
}

func getNamespaceNameIndex(name string, ns []string) int {
	var index int
	for i, n := range ns {
		if n == name {
			index = i
		}
	}

	return index
}

// func (v *View) filterNamespaces(data []*models.Namespace) []*models.Namespace {
// 	filter := v.state.Filter.Namespace
// 	if filter != "" {
// 		rx, _ := regexp.Compile(filter)
// 		result := []*models.Namespace{}
// 		for _, ns := range v.state.Namespaces {
// 			switch true {
// 			case rx.MatchString(ns.Name),
// 				rx.MatchString(ns.Description):
// 				result = append(result, ns)
// 			}
// 		}

// 		return result
// 	}

// 	return data
// }
