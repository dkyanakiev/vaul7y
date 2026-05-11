package view

import (
	"regexp"
	"time"

	"github.com/dkyanakiev/vaul7y/tui/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) VPolicy() {
	v.viewSwitch()
	v.Layout.Body.Clear()
	v.Layout.Body.SetTitle("Vault Policies")
	v.Layout.Container.SetFocus(v.components.PolicyTable.Table.Primitive())
	v.Layout.Container.SetInputCapture(v.InputVaultPolicy)
	v.components.Commands.Update(component.PolicyCommands)
	search := v.components.Search
	//table := v.components.in

	update := func() {
		v.state.RLock()
		searching := v.state.Toggle.Search
		v.state.RUnlock()

		v.mutex.RLock()
		filterText := v.FilterText
		v.mutex.RUnlock()

		if searching {
			v.state.Lock()
			v.state.Filter.Policy = filterText
			v.state.Unlock()
		}

		v.state.RLock()
		data := v.filterPolicies()
		v.state.RUnlock()

		v.components.PolicyTable.Props.Data = data
		v.components.PolicyTable.Render()
		v.Draw()
		v.components.PolicyTable.Table.ScrollToTop()
	}

	var debounce *time.Timer
	search.Props.ChangedFunc = func(text string) {
		v.mutex.Lock()
		v.FilterText = text
		v.mutex.Unlock()
		if debounce != nil {
			debounce.Stop()
		}
		debounce = time.AfterFunc(150*time.Millisecond, func() {
			v.Layout.Container.QueueUpdateDraw(update)
		})
	}

	v.Watcher.SubscribeToPolicies(func() { v.Layout.Container.QueueUpdateDraw(update) })
	update()

	v.state.Elements.TableMain = v.components.PolicyTable.Table.Primitive().(*tview.Table)
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
		rx, err := regexp.Compile(filter)
		if err != nil {
			return data
		}
		result := []string{}
		for _, p := range data {
			if rx.MatchString(p) {
				result = append(result, p)
			}
		}
		return result
	}
	return data
}
