package view

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/dkyanakiev/vaul7y/tui/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) SecretObject(mount, path string) {
	v.viewSwitch()
	v.components.SecretObjTable.SetEditable(false)
	v.components.SecretObjTable.ShowMetadata = false
	v.components.SecretObjTable.ShowJson = false
	v.components.TogglesInfo.Props.Editable = false
	v.components.TogglesInfo.Render()
	v.Layout.Body.SetTitle("Secret object")
	v.Layout.Container.SetInputCapture(v.InputSecret)
	v.Layout.Container.SetFocus(v.components.SecretObjTable.Table.Primitive())
	v.components.Commands.Update(component.SecretObjectCommands)

	v.logger.Debug().Msgf("Selected mount is: %v", mount)
	v.logger.Debug().Msgf("Selected path is: %v", path)
	v.state.Elements.TableMain = v.components.SecretObjTable.Table.Primitive().(*tview.Table)
	v.components.SecretObjTable.Logger = v.logger
	v.components.SecretObjTable.Props.SelectedPath = path
	v.components.SecretObjTable.Props.ObscureSecrets = true

	update := func() {
		v.state.RLock()
		secret := v.state.SelectedSecret
		meta := v.state.SelectedSecretMeta
		selectedPath := v.state.SelectedPath
		v.state.RUnlock()

		v.logger.Debug().Msgf("Selected path is: %v", selectedPath)
		if !v.components.SecretObjTable.IsEditable() {
			v.components.SecretObjTable.Props.Data = secret
			v.components.SecretObjTable.Props.Metadata = meta
			v.components.SecretObjTable.Render()
			v.Draw()
		}
	}

	v.Watcher.SubscribeToSecret(mount, path, func() { v.Layout.Container.QueueUpdateDraw(update) })
	update()

	v.state.Elements.TableMain = v.components.SecretObjTable.Table.Primitive().(*tview.Table)

}

func (v *View) inputSecret(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyRune:
		if !v.components.SecretObjTable.IsEditable() {
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
			case 't':
				v.state.RLock()
				meta := v.state.SelectedSecretMeta
				v.state.RUnlock()
				if meta == nil {
					return nil // v1 engine — no metadata endpoint
				}
				v.components.SecretObjTable.ShowMetadata = !v.components.SecretObjTable.ShowMetadata
				v.logger.Debug().Msgf("Toggle metadata panel: %v", v.components.SecretObjTable.ShowMetadata)
				v.components.SecretObjTable.ToggleMetaView()
			case 'j':
				v.components.SecretObjTable.ShowJson = !v.components.SecretObjTable.ShowJson
				v.components.SecretObjTable.ToggleView()
			case 'P':
				v.promptUpdateFormat("PATCH")
				return nil
			case 'U':
				v.promptUpdateFormat("UPDATE")
				return nil
			case 'D':
				v.confirmDeleteVersion()
				return nil
			case 'R':
				v.promptRollback()
				return nil
			case 'X':
				v.promptDestroyVersion()
				return nil
			case 'M':
				v.editSecretMetadata()
				return nil
			}
		}
	case tcell.KeyCtrlW:
		v.confirmSecretUpdate()
		return nil
	case tcell.KeyEsc:
		if v.components.SecretObjTable.IsEditable() {
			v.components.SecretObjTable.SetEditable(false)
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
	v.state.Lock()
	v.state.SelectedPath = strings.TrimSuffix(v.state.SelectedPath, "/")
	lastSlashIndex := strings.LastIndex(v.state.SelectedPath, "/")
	if lastSlashIndex != -1 {
		v.state.SelectedPath = v.state.SelectedPath[:lastSlashIndex+1]
	} else if v.state.SelectedPath != "" {
		v.state.SelectedPath = ""
		v.components.SecretsTable.Props.SelectedPath = ""
	}
	path := v.state.SelectedPath
	v.state.Unlock()
	v.Secrets(path, "false")
}

// patchSecret validates, saves, and submits the edited JSON. Returns true on
// success so the caller knows whether to navigate away.
func (v *View) patchSecret() bool {
	patch := v.components.SecretObjTable.Props.Update == "PATCH"
	v.logger.Debug().Msgf("patchSecret: patch=%v", patch)

	errMsg := v.components.SecretObjTable.SaveData(v.components.SecretObjTable.TextArea.GetText())
	if errMsg != "" {
		// Show the error inline in the TextArea title so the user stays in
		// edit mode and can correct the JSON without losing their work.
		v.components.SecretObjTable.SetEditStatus(errMsg)
		v.Draw()
		return false
	}

	v.state.RLock()
	enterprise := v.state.Enterprise
	ns := v.state.SelectedNamespace
	mount := v.state.SelectedMount
	v.state.RUnlock()

	if enterprise {
		v.logger.Debug().Msgf("Enterprise version detected, setting namespace to %v", ns)
		v.Client.ChangeNamespace(ns)
	}
	v.logger.Debug().Msgf("Updated secret object is: %v", v.components.SecretObjTable.Props.UpdatedData)
	err := v.Client.UpdateSecretObjectKV2(mount, v.components.SecretObjTable.Props.SelectedPath, patch, v.components.SecretObjTable.Props.UpdatedData)
	if err != nil {
		v.handleError(err.Error())
		return false
	}

	v.components.SecretObjTable.SetEditable(false)
	v.components.SecretObjTable.ToggleView()
	return true
}

func (v *View) confirmDeleteVersion() {
	mount := v.state.SelectedMount
	path := v.components.SecretObjTable.Props.SelectedPath
	msg := fmt.Sprintf("Delete current version of:\n%s\n\nContinue?", path)

	v.components.Confirm.Props.Done = func(_ int, label string) {
		v.Layout.Pages.RemovePage(component.PageNameConfirm)
		if label == "Yes" {
			go func() {
				if err := v.Client.DeleteCurrentSecretVersion(mount, path); err != nil {
					v.Layout.Container.QueueUpdateDraw(func() { v.handleError(err.Error()) })
				} else {
					v.Layout.Container.QueueUpdateDraw(v.goBack)
				}
			}()
		} else {
			v.Layout.Container.SetFocus(v.components.SecretObjTable.Table.Primitive())
		}
	}
	v.components.Confirm.Render(msg) //nolint:errcheck
	v.Layout.Container.SetFocus(v.components.Confirm.FocusPrimitive())
}

func (v *View) promptRollback() {
	origDoneFunc := v.components.TextInfoInput.Props.DoneFunc

	v.state.Lock()
	v.state.Toggle.TextInput = true
	v.state.Toggle.JumpToPath = false
	v.state.Unlock()
	v.components.TextInfoInput.InputField.SetLabel("Rollback to version: ")
	v.components.TextInfoInput.InputField.SetText("")

	mount := v.state.SelectedMount
	path := v.components.SecretObjTable.Props.SelectedPath

	v.components.TextInfoInput.Props.DoneFunc = func(key tcell.Key) {
		v.components.TextInfoInput.Props.DoneFunc = origDoneFunc
		text := v.components.TextInfoInput.InputField.GetText()
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.Layout.Footer.RemoveItem(v.components.TextInfoInput.InputField.Primitive())
		v.state.Lock()
		v.state.Toggle.TextInput = false
		v.state.Unlock()
		v.components.TextInfoInput.InputField.SetLabel("")

		version, err := strconv.Atoi(strings.TrimSpace(text))
		if err != nil || version < 1 {
			v.handleError("Invalid version number")
			return
		}
		go func() {
			if err := v.Client.RollbackSecret(mount, path, version); err != nil {
				v.Layout.Container.QueueUpdateDraw(func() { v.handleError(err.Error()) })
			} else {
				v.Layout.Container.QueueUpdateDraw(func() { v.SecretObject(mount, path) })
			}
		}()
	}

	v.TextInput()
}

func (v *View) promptDestroyVersion() {
	origDoneFunc := v.components.TextInfoInput.Props.DoneFunc

	v.state.Lock()
	v.state.Toggle.TextInput = true
	v.state.Toggle.JumpToPath = false
	v.state.Unlock()
	v.components.TextInfoInput.InputField.SetLabel("Destroy version #: ")
	v.components.TextInfoInput.InputField.SetText("")

	mount := v.state.SelectedMount
	path := v.components.SecretObjTable.Props.SelectedPath

	v.components.TextInfoInput.Props.DoneFunc = func(key tcell.Key) {
		v.components.TextInfoInput.Props.DoneFunc = origDoneFunc
		text := v.components.TextInfoInput.InputField.GetText()
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.Layout.Footer.RemoveItem(v.components.TextInfoInput.InputField.Primitive())
		v.state.Lock()
		v.state.Toggle.TextInput = false
		v.state.Unlock()
		v.components.TextInfoInput.InputField.SetLabel("")

		version, err := strconv.Atoi(strings.TrimSpace(text))
		if err != nil || version < 1 {
			v.handleError("Invalid version number")
			return
		}

		msg := fmt.Sprintf("PERMANENTLY DESTROY version %d of:\n%s\n\nThis CANNOT be undone. Continue?", version, path)
		v.components.Confirm.Props.Done = func(_ int, label string) {
			v.Layout.Pages.RemovePage(component.PageNameConfirm)
			if label == "Yes" {
				go func() {
					if err := v.Client.DestroySecretVersions(mount, path, []int{version}); err != nil {
						v.Layout.Container.QueueUpdateDraw(func() { v.handleError(err.Error()) })
					} else {
						v.Layout.Container.QueueUpdateDraw(func() { v.SecretObject(mount, path) })
					}
				}()
			} else {
				v.Layout.Container.SetFocus(v.components.SecretObjTable.Table.Primitive())
			}
		}
		v.components.Confirm.Render(msg) //nolint:errcheck
		v.Layout.Container.SetFocus(v.components.Confirm.FocusPrimitive())
	}

	v.TextInput()
}

func (v *View) editSecretMetadata() {
	v.state.RLock()
	meta := v.state.SelectedSecretMeta
	v.state.RUnlock()
	if meta == nil {
		v.handleError("Metadata editing is only supported for KV v2")
		return
	}

	mount := v.state.SelectedMount
	path := v.components.SecretObjTable.Props.SelectedPath

	initial := map[string]interface{}{
		"max_versions":         meta.MaxVersions,
		"cas_required":         meta.CasRequired,
		"delete_version_after": meta.DeleteVersionAfter,
	}
	jsonBytes, _ := json.MarshalIndent(initial, "", "  ")

	v.components.Commands.Update(component.SecretsObjectPatchCommands)
	v.components.SecretObjTable.SetEditable(true)
	v.components.TogglesInfo.Props.Editable = true
	v.components.SecretObjTable.Props.Update = "METADATA"
	v.components.TogglesInfo.Render()
	v.components.SecretObjTable.ToggleView()
	v.components.SecretObjTable.TextArea.SetText(string(jsonBytes), true)
	v.components.SecretObjTable.TextView.ScrollToBeginning()
	v.Layout.Container.SetFocus(v.components.SecretObjTable.TextArea.Primitive())

	origDone := v.components.Confirm.Props.Done
	v.components.Confirm.Props.Done = func(btnIdx int, label string) {
		v.Layout.Pages.RemovePage(component.PageNameConfirm)
		if label == "Yes" {
			raw := v.components.SecretObjTable.TextArea.GetText()
			var opts map[string]interface{}
			if err := json.Unmarshal([]byte(raw), &opts); err != nil {
				v.components.SecretObjTable.SetEditStatus("Invalid JSON: " + err.Error())
				v.Draw()
				v.Layout.Container.SetFocus(v.components.SecretObjTable.TextArea.Primitive())
				return
			}
			optsToSend := opts
			go func() {
				opts := optsToSend
				if err := v.Client.UpdateSecretMetadata(mount, path, opts); err != nil {
					v.Layout.Container.QueueUpdateDraw(func() {
						v.handleError(err.Error())
						v.components.SecretObjTable.SetEditable(false)
						v.components.TogglesInfo.Props.Editable = false
						v.components.TogglesInfo.Render()
						v.components.SecretObjTable.ToggleView()
					})
				} else {
					v.Layout.Container.QueueUpdateDraw(func() {
						v.components.SecretObjTable.SetEditable(false)
						v.components.TogglesInfo.Props.Editable = false
						v.components.TogglesInfo.Render()
						v.SecretObject(mount, path)
					})
				}
			}()
		} else {
			v.components.SecretObjTable.SetEditable(false)
			v.components.TogglesInfo.Props.Editable = false
			v.components.TogglesInfo.Render()
			v.components.SecretObjTable.ToggleView()
			v.Layout.Container.SetFocus(v.components.SecretObjTable.Table.Primitive())
		}
		v.components.Confirm.Props.Done = origDone
	}
}

// promptUpdateFormat shows the JSON/Key-Value format selector for P/U operations.
func (v *View) promptUpdateFormat(mode string) {
	mount := v.state.SelectedMount
	path := v.components.SecretObjTable.Props.SelectedPath
	patch := mode == "PATCH"

	openJSONEditor := func() {
		v.components.Commands.Update(component.SecretsObjectPatchCommands)
		v.components.SecretObjTable.SetEditable(true)
		v.components.TogglesInfo.Props.Editable = true
		v.components.SecretObjTable.Props.Update = mode
		v.components.TogglesInfo.Render()
		v.components.SecretObjTable.ToggleView()
		v.components.SecretObjTable.TextView.ScrollToBeginning()
		v.Layout.Container.SetFocus(v.components.SecretObjTable.TextArea.Primitive())
	}

	openKVForm := func() {
		var initialData map[string]interface{}
		if d := v.components.SecretObjTable.Props.Data; d != nil && d.Data != nil {
			if kv, ok := d.Data["data"].(map[string]interface{}); ok {
				initialData = kv
			}
		}
		v.showKVForm(mount, path, initialData, patch,
			func() { v.SecretObject(mount, path) },
			v.InputSecret,
		)
	}

	v.showFormatSelector(
		fmt.Sprintf("%s Secret", mode),
		fmt.Sprintf("Choose format for %s:", path),
		openJSONEditor,
		openKVForm,
		v.InputSecret,
	)
}

func (v *View) confirmSecretUpdate() {
	op := v.components.SecretObjTable.Props.Update
	path := v.components.SecretObjTable.Props.SelectedPath
	msg := fmt.Sprintf("%s secret at path:\n%s\n\nContinue?", op, path)

	v.components.Confirm.Props.Done = func(_ int, label string) {
		v.Layout.Pages.RemovePage(component.PageNameConfirm)
		if label == "Yes" {
			if v.patchSecret() {
				v.goBack()
			} else {
				v.Layout.Container.SetFocus(v.components.SecretObjTable.TextArea.Primitive())
			}
		} else {
			v.Layout.Container.SetFocus(v.components.SecretObjTable.TextArea.Primitive())
		}
	}
	v.components.Confirm.Render(msg) //nolint:errcheck
	v.Layout.Container.SetFocus(v.components.Confirm.FocusPrimitive())
}
