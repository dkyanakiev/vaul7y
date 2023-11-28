package view

import (
	"github.com/atotto/clipboard"
	"github.com/dkyanakiev/vaulty/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) SecretObject(mount, path string) {
	v.viewSwitch()
	v.Layout.Body.SetTitle("Secret obejects")
	v.Layout.Container.SetInputCapture(v.InputSecret)
	v.components.Commands.Update(component.SecretObjectCommands)
	v.Layout.Container.SetFocus(v.components.SecretObjTable.Table.Primitive())

	v.state.Elements.TableMain = v.components.SecretObjTable.Table.Primitive().(*tview.Table)
	v.components.SecretObjTable.Logger = v.logger
	v.components.SecretObjTable.Props.SelectedPath = path
	v.components.SecretObjTable.Props.ObscureSecrets = true

	update := func() {

		v.logger.Debug().Msgf("Selected path is: %v", v.state.SelectedPath)
		if !v.components.SecretObjTable.Editable {
			v.components.SecretObjTable.Render()

			v.components.SecretObjTable.Props.Data = v.state.SelectedSecret
			v.logger.Debug().Msgf("Updated secret object is: %v", v.components.SecretObjTable.Props.UpdatedData)
			v.Draw()
		}
	}

	v.Watcher.SubscribeToSecret(mount, path, update)
	update()

}

func (v *View) inputSecret(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyRune:
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
		case 'j':
			v.components.SecretObjTable.ShowJson = !v.components.SecretObjTable.ShowJson
			v.components.SecretObjTable.ToggleView()
		case 'U':
			v.components.SecretObjTable.Editable = true
			v.components.SecretObjTable.ToggleView()
		case 'o':
			v.Layout.Pages.RemovePage(component.PageNameInfo)
			v.Layout.Pages.RemovePage(component.PageNameError)
		default:
			if v.components.SecretObjTable.Editable {
				v.components.SecretObjTable.TextView.SetText(v.components.SecretObjTable.TextView.GetText(true) + string(event.Rune()))
			}
		}
	case tcell.KeyCtrlD:
		// Delete all text from the view
		v.components.SecretObjTable.TextView.SetText("")
		return nil
	case tcell.KeyCtrlS:
		// Save the text and make the text view uneditable again
		ok := v.components.SecretObjTable.SaveData(v.components.SecretObjTable.TextView.GetText(true))
		if ok != "" {
			v.handleInfo(ok)
		}
		//TMP: Disabled for now
		//v.Client.UpdateSecretObject("kv0FF76557", "data/secret1", true, v.components.SecretObjTable.Props.UpdatedData)
		//err := v.Client.UpdateSecretObject("kv0FF76557", "data/secret1", true, v.components.SecretObjTable.Props.UpdatedData)
		// if err != nil {
		// 	v.handleError(string(err.Error()))
		// }
		v.components.SecretObjTable.Editable = false
		return nil
	case tcell.KeyEsc:
		v.components.SecretObjTable.Editable = false
		v.components.SecretObjTable.ToggleView()
		return nil
	case tcell.KeyBackspace2, tcell.KeyBackspace:
		// If text editing is enabled, handle backspace
		if v.components.SecretObjTable.Editable {
			text := v.components.SecretObjTable.TextView.GetText(true)
			if len(text) > 0 {
				v.components.SecretObjTable.TextView.SetText(text[:len(text)-1])
			}
		}
		return nil

	}

	return event
}
