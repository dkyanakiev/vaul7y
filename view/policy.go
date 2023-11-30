package view

import (
	"regexp"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) VPolicy() {
	v.viewSwitch()
	v.Layout.Body.Clear()
	v.Layout.Body.SetTitle("Vault Policies")
	v.Layout.Container.SetInputCapture(v.InputVaultPolicy)
	v.components.Commands.Update(component.PolicyCommands)
	search := v.components.Search
	//table := v.components.in

	update := func() {
		if v.state.Toggle.Search {
			v.state.Filter.Policy = v.FilterText
		}
		v.components.PolicyTable.Props.Data = v.filterPolicies()
		v.components.PolicyTable.Render()
		v.Draw()
	}

	search.Props.ChangedFunc = func(text string) {
		v.FilterText = text
		update()
	}

	v.Watcher.SubscribeToPolicies(update)
	update()

	v.state.Elements.TableMain = v.components.PolicyTable.Table.Primitive().(*tview.Table)
	v.Layout.Container.SetFocus(v.components.PolicyTable.Table.Primitive())
}

func (v *View) inputPolicy(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyEsc:
		//v.GoBack()
	case tcell.KeyEnter:
		if v.components.PolicyTable.Table.Primitive().HasFocus() {
			v.PolicyACL(v.components.PolicyTable.GetIDForSelection())
			v.Watcher.Unsubscribe()
			return nil
		}
	case tcell.KeyRune:
		switch event.Rune() {
		case '/':
			if !v.Layout.Footer.HasFocus() {
				if !v.state.Toggle.Search {
					v.state.Toggle.Search = true
					v.components.Search.InputField.SetText("")
					v.Search()
				} else {
					v.Layout.Container.SetFocus(v.components.Search.InputField.Primitive())
				}
				return nil
			}
		case 'i':
			if v.components.PolicyTable.Table.Primitive().HasFocus() {
				v.PolicyACL(v.components.PolicyTable.GetIDForSelection())
				v.Watcher.Unsubscribe()
				return nil
			}
		}

	}

	return event
}

func (v *View) filterPolicies() []string {
	data := v.state.PolicyList
	filter := v.state.Filter.Policy
	if filter != "" {
		rx, _ := regexp.Compile(filter)
		result := []string{}
		for _, p := range data {
			switch true {
			case rx.MatchString(p):
				result = append(result, p)
			}
		}
		return result
	}
	return data
}
