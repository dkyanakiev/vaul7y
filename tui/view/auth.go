package view

import (
	"github.com/dkyanakiev/vaul7y/tui/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) AuthMethods() {
	v.viewSwitch()
	v.Layout.Body.Clear()
	v.Layout.Body.SetTitle("Auth Methods")
	v.Layout.Container.SetFocus(v.components.AuthTable.Table.Primitive())
	v.Layout.Container.SetInputCapture(v.InputAuthMethods)
	v.components.Commands.Update(component.AuthCommands)

	update := func() {
		v.state.RLock()
		data := v.state.AuthMethods
		v.state.RUnlock()

		v.components.AuthTable.Props.Data = data
		v.components.AuthTable.Render()
		v.Draw()
		v.components.AuthTable.Table.ScrollToTop()
	}

	v.Watcher.SubscribeToAuthMethods(func() { v.Layout.Container.QueueUpdateDraw(update) })
	update()

	v.state.Elements.TableMain = v.components.AuthTable.Table.Primitive().(*tview.Table)
}

func (v *View) inputAuthMethods(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}
	switch event.Key() {
	case tcell.KeyEsc:
		v.Mounts()
	}
	return event
}
