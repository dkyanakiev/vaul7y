package view

import (
	"strings"

	"github.com/atotto/clipboard"
	"github.com/dkyanakiev/vaulty/tui/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) SecretObject(mount, path string) {
	v.viewSwitch()
	v.Layout.Body.SetTitle("Secret object")
	v.Layout.Container.SetInputCapture(v.InputSecret)
	v.components.Commands.Update(component.SecretObjectCommands)

	v.logger.Debug().Msgf("Selected mount is: %v", mount)
	v.logger.Debug().Msgf("Selected path is: %v", path)
	v.state.Elements.TableMain = v.components.SecretObjTable.Table.Primitive().(*tview.Table)
	v.components.SecretObjTable.Logger = v.logger
	v.components.SecretObjTable.Props.SelectedPath = path
	v.components.SecretObjTable.Props.ObscureSecrets = true

	update := func() {
		v.logger.Debug().Msgf("Focus set to %s", v.state.Elements.TableMain.GetTitle())
		v.logger.Debug().Msgf("Selected path is: %v", v.state.SelectedPath)
		if !v.components.SecretObjTable.Editable {
			v.components.SecretObjTable.Render()
			v.components.SecretObjTable.Props.Data = v.state.SelectedSecret
			v.Draw()
		}
	}

	v.Watcher.SubscribeToSecret(mount, path, update)
	update()

	v.state.Elements.TableMain = v.components.SecretObjTable.Table.Primitive().(*tview.Table)
	v.Layout.Container.SetFocus(v.components.SecretObjTable.Table.Primitive())

	// v.addToHistory(v.state.SelectedNamespace, "secret", func() {
	// 	v.SecretObject(mount, path)
	// })

}

func (v *View) inputSecret(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyRune:
		if !v.components.SecretObjTable.Editable {
			switch event.Rune() {
			case 'h':
				v.components.SecretObjTable.Props.ObscureSecrets = !v.components.SecretObjTable.Props.ObscureSecrets
				v.components.SecretObjTable.Render()
				return nil
			case 'c':
				if v.components.SecretObjTable.ShowJson {
					content := v.components.SecretObjTable.TextView.GetText(true)
					clipboard.WriteAll(content)
				} else {
					row, _ := v.components.SecretObjTable.Table.GetSelection()
					if row > 0 { // Ignore the header row
						// Get the content of the row
						content := v.components.SecretObjTable.Table.GetCellContent(row, 1)
						// Copy the content to the clipboard
						clipboard.WriteAll(content)
					}
				}
				return nil
			case 'b':
				v.goBack()
			case 'j':
				v.components.SecretObjTable.ShowJson = !v.components.SecretObjTable.ShowJson
				v.components.SecretObjTable.ToggleView()
			case 'P':
				v.components.Commands.Update(component.SecretsObjectPatchCommands)
				v.components.SecretObjTable.Editable = true
				v.components.TogglesInfo.Props.Editable = true
				v.components.SecretObjTable.Props.Update = "PATCH"
				v.components.TogglesInfo.Render()
				v.components.SecretObjTable.ToggleView()
				v.components.SecretObjTable.TextView.ScrollToBeginning()
				v.Layout.Container.SetFocus(v.components.SecretObjTable.TextArea.Primitive())
				return nil
			case 'U':
				v.components.Commands.Update(component.SecretsObjectPatchCommands)
				v.components.SecretObjTable.Editable = true
				v.components.TogglesInfo.Props.Editable = true
				v.components.SecretObjTable.Props.Update = "UPDATE"
				v.components.TogglesInfo.Render()
				v.components.SecretObjTable.ToggleView()
				v.components.SecretObjTable.TextView.ScrollToBeginning()
				v.Layout.Container.SetFocus(v.components.SecretObjTable.TextArea.Primitive())
				return nil
			}
		}
	case tcell.KeyCtrlW:
		v.patchSecret()
		v.goBack()
		return nil
	case tcell.KeyEsc:
		if v.components.SecretObjTable.Editable {
			v.components.SecretObjTable.Editable = false
			v.components.TogglesInfo.Props.Editable = false
			v.components.TogglesInfo.Render()
			v.components.SecretObjTable.ToggleView()
			v.Layout.Container.SetFocus(v.components.SecretObjTable.Table.Primitive())
		} else {
			v.goBack()
		}

	}

	return event
}

func (v *View) goBack() {
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

func (v *View) patchSecret() {
	var patch bool
	if v.components.SecretObjTable.Props.Update == "PATCH" {
		v.logger.Debug().Msg("PATCH secret")
		patch = true
	} else {
		v.logger.Debug().Msg("UPDATE secret")
		patch = false
	}
	ok := v.components.SecretObjTable.SaveData(v.components.SecretObjTable.TextArea.GetText())
	if ok != "" {
		v.handleError(ok)
	} else {
		if v.state.Enterprise {
			v.logger.Debug().Msgf("Enterprise version detected, setting namespace to %v", v.state.SelectedNamespace)
			v.Client.ChangeNamespace(v.state.SelectedNamespace)
		}
		v.logger.Debug().Msgf("Updated secret object is: %v", v.components.SecretObjTable.Props.UpdatedData)
		err := v.Client.UpdateSecretObjectKV2(v.state.SelectedMount, v.components.SecretObjTable.Props.SelectedPath, patch, v.components.SecretObjTable.Props.UpdatedData)

		if err != nil {
			v.handleError(string(err.Error()))
		}
		v.components.SecretObjTable.Editable = false
		v.components.SecretObjTable.ToggleView()
	}
}
