package view

import (
	"github.com/dkyanakiev/vaulty/tui/component"
	"github.com/rivo/tview"
)

func (v *View) Namespaces() {
	v.viewSwitch()

	v.Layout.Body.SetTitle("namespaces")

	v.state.Elements.TableMain = v.components.NamespaceTable.Table.Primitive().(*tview.Table)
	v.components.Commands.Update(component.NoViewCommands)
	v.Layout.Container.SetInputCapture(v.InputNamespaces)

	update := func() {
		// v.components.NamespaceTable.Props.Data = v.filterNamespaces(v.state.Namespaces)
		v.components.NamespaceTable.Render()
		v.Draw()
	}

	v.components.Search.Props.ChangedFunc = func(text string) {
		v.state.Filter.Namespace = text
		update()
	}

	// v.Watcher.SubscribeToNamespaces(update)

	// update()

	// v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
	// 	v.state.SelectedNamespace = text
	// 	v.Namespaces()
	// })

	// v.addToHistory(v.state.SelectedNamespace, models.TopicNamespace, v.Namespaces)
	v.Layout.Container.SetFocus(v.components.NamespaceTable.Table.Primitive())
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
