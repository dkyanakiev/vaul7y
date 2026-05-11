package view

import (
	"github.com/dkyanakiev/vaul7y/internal/models"
	"github.com/dkyanakiev/vaul7y/tui/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) Mounts() {
	v.logger.Debug().Msg("Mounts view")
	v.viewSwitch()
	v.Layout.Body.SetTitle("Secret Mounts")
	v.Layout.Container.SetFocus(v.components.MountsTable.Table.Primitive())
	v.Layout.Container.SetInputCapture(v.InputMounts)
	v.components.Commands.Update(component.MountsCommands)
	v.state.Elements.TableMain = v.components.MountsTable.Table.Primitive().(*tview.Table)
	v.components.MountsTable.Logger = v.logger
	v.components.SecretsTable.Props.SelectedPath = ""
	v.state.SelectedMount = ""
	v.state.SelectedPath = ""
	v.state.SelectedObject = ""

	update := func() {
		v.state.RLock()
		mounts := v.state.Mounts
		selectedMount := v.state.SelectedMount
		selectedPath := v.state.SelectedPath
		v.state.RUnlock()

		v.components.MountsTable.Props.Data = mounts
		v.components.MountsTable.Render()
		v.Draw()
		v.logger.Debug().Msg("Updated mounts table")
		v.logger.Debug().Msgf("Selected mount is: %v", selectedMount)
		v.logger.Debug().Msgf("Selected path is: %v", selectedPath)
	}

	v.Watcher.SubscribeToMounts(func() { v.Layout.Container.QueueUpdateDraw(update) })
	// v.Watcher.UpdateMounts()
	update()

	// Add this view to the history
	// v.addToHistory(v.state.SelectedNamespace, "mounts", func() {
	// 	v.Mounts()
	// })
}

func (v *View) parseMounts(data []*models.MountOutput) []*models.MountOutput {
	return nil
}

func (v *View) inputMounts(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}
	//todo
	switch event.Key() {
	case tcell.KeyEnter:
		if v.components.MountsTable.Table.Primitive().HasFocus() {
			v.state.Lock()
			v.state.SelectedMount = v.components.MountsTable.GetIDForSelection()
			v.state.Unlock()
			v.Secrets("", "false")
			return nil
		}
	case tcell.KeyRune:
		switch event.Rune() {
		case 'e':
			if v.components.MountsTable.Table.Primitive().HasFocus() {
				v.state.Lock()
				v.state.SelectedMount = v.components.MountsTable.GetIDForSelection()
				v.state.Unlock()
				v.Secrets("", "false")
				return nil
			}
		}
	}

	return event
}
