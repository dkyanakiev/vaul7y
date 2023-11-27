package view

import (
	"github.com/dkyanakiev/vaulty/component"
	"github.com/dkyanakiev/vaulty/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) Mounts() {
	v.viewSwitch()
	v.Layout.Body.SetTitle("Secret Mounts")
	v.Layout.Container.SetInputCapture(v.InputMounts)
	v.components.Commands.Update(component.MountsCommands)
	v.components.MountsTable.Logger = v.logger
	v.components.SecretsTable.Props.SelectedPath = ""
	v.state.SelectedMount = ""
	v.state.SelectedPath = ""
	v.state.SelectedObject = ""

	update := func() {
		v.components.MountsTable.Props.Data = v.state.Mounts
		v.components.MountsTable.Render()
		v.Draw()
		v.logger.Println("Updated mounts table")
		v.logger.Println("Selected mount is: ", v.state.SelectedMount)
		v.logger.Println("Selected path is: ", v.state.SelectedPath)
	}

	//v.Watcher.SubscribeToMounts(update)
	v.Watcher.UpdateMounts()
	update()

	v.state.Elements.TableMain = v.components.MountsTable.Table.Primitive().(*tview.Table)
	v.Layout.Container.SetFocus(v.components.MountsTable.Table.Primitive())

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
	case tcell.KeyRune:
		switch event.Rune() {
		case 'e':
			if v.components.MountsTable.Table.Primitive().HasFocus() {
				v.state.SelectedMount = v.components.MountsTable.GetIDForSelection()
				v.Secrets("", "false")
				return nil
			}
		}
	}

	return event
}
