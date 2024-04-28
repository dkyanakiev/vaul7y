package view

import (
	"fmt"

	"github.com/dkyanakiev/vaulty/tui/component"
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
		// v.components.NamespaceTable.Props.Data = v.filterNamespaces(v.state.Namespaces)
		v.logger.Debug().Msgf("Current ns list: %v", v.state.Namespaces)
		v.components.NamespaceTable.Props.Data = v.state.Namespaces
		v.components.NamespaceTable.Render()
		v.Draw()
		v.components.NamespaceTable.Table.ScrollToTop()
	}

	v.components.Search.Props.ChangedFunc = func(text string) {
		//v.state.Filter.Namespace = text
		update()
	}

	v.Watcher.SubscribeToNamespaces(update)

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
		v.logger.Debug().Msgf("Going back to default namespace: %v", v.state.DefaultNamespace)
		v.state.SelectedNamespace = v.state.DefaultNamespace
		v.components.TogglesInfo.Props.Namespace = v.state.DefaultNamespace
		v.components.TogglesInfo.Render()
		v.state.SelectedNamespace = v.state.DefaultNamespace
		//v.Client.ChangeNamespace(v.state.DefaultNamespace)
		v.Watcher.Unsubscribe()
		v.Mounts()
		return nil
	case tcell.KeyCtrlW:
		v.logger.Debug().Msgf("Going back to root namespace : %v", v.state.RootNamespace)
		v.state.SelectedNamespace = v.state.RootNamespace
		v.components.TogglesInfo.Props.Namespace = v.state.RootNamespace
		v.components.TogglesInfo.Render()
		v.state.SelectedNamespace = v.state.RootNamespace
		// v.Client.ChangeNamespace(v.state.RootNamespace)
		v.Watcher.Unsubscribe()
		v.Mounts()
		return nil
	case tcell.KeyEnter:
		selectdNs := v.components.NamespaceTable.GetIDForSelection()
		v.logger.Debug().Msgf("Selected namespace is: %v", selectdNs)
		newNs := fmt.Sprintf("%s/%s", v.state.SelectedNamespace, selectdNs)
		v.logger.Debug().Msgf("Changing namespace to: %s", newNs)
		v.components.TogglesInfo.Props.Namespace = newNs
		v.components.TogglesInfo.Render()
		v.state.SelectedNamespace = newNs
		// v.Client.ChangeNamespace(newNs)
		v.Watcher.Unsubscribe()
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
