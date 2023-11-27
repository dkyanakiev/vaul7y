package view

import (
	"fmt"
	"strings"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// base path - v.state.SelectedMount
// latest path - path
// full path - path + path

func (v *View) Secrets(path string, secretBool string) {
	v.viewSwitch()
	v.Layout.Body.Clear()
	v.Layout.Body.SetTitle(fmt.Sprintf("Secrets: %s", path))
	v.Layout.Container.SetInputCapture(v.InputSecrets)
	v.components.Commands.Update(component.SecretsCommands)

	v.components.SecretsTable.Props.SelectedMount = v.state.SelectedMount
	if path != "" {
		v.components.SecretsTable.Props.SelectedPath = fmt.Sprintf("%s%s", v.state.SelectedPath, v.state.SelectedObject)
	}

	update := func() {
		v.components.SecretsTable.Props.Data = v.state.SecretsData
		v.components.SecretsTable.Props.SelectedMount = v.state.SelectedMount

		v.components.SecretsTable.Render()
		v.Draw()
	}

	v.Watcher.SubscribeToSecrets(v.components.SecretsTable.Props.SelectedMount,
		v.components.SecretsTable.Props.SelectedPath, update)
	update()

	v.state.Elements.TableMain = v.components.SecretsTable.Table.Primitive().(*tview.Table)
	v.Layout.Container.SetFocus(v.components.SecretsTable.Table.Primitive())

}

func (v *View) inputSecrets(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'e':
			if v.components.SecretsTable.Table.Primitive().HasFocus() {
				path, secretBool := v.components.SecretsTable.GetIDForSelection()
				v.state.SelectedPath = fmt.Sprintf("%s%s", v.state.SelectedPath, path)
				if secretBool == "true" {
					v.SecretObject(v.state.SelectedMount, v.state.SelectedPath)
				} else {
					v.Secrets(path, secretBool)
				}
				return nil
			}

		//TODO: Need to clean this up
		case 'b':
			if v.components.SecretsTable.Table.Primitive().HasFocus() {
				v.state.SelectedPath = strings.TrimSuffix(v.state.SelectedPath, "/") // Remove trailing slash
				lastSlashIndex := strings.LastIndex(v.state.SelectedPath, "/")
				if lastSlashIndex != -1 {
					v.state.SelectedPath = v.state.SelectedPath[:lastSlashIndex+1] // Keep the slash
				} else if v.state.SelectedPath != "" {
					v.state.SelectedPath = "" // If no slash left and it's not empty, set to empty
					v.components.SecretsTable.Props.SelectedPath = ""
				}
				v.Secrets(v.state.SelectedPath, "false")
			}
		}
	}

	return event
}
