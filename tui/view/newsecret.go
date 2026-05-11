package view

import (
	"fmt"

	"github.com/dkyanakiev/vaul7y/tui/component"
	primitive "github.com/dkyanakiev/vaul7y/tui/primitives"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	pageFormatSel = "formatsel"
	pageKVForm    = "kvform"
)

// showFormatSelector shows a "JSON / Key-Value / Cancel" modal.
// restoreCapture is reinstated when the modal is dismissed so the caller's
// view input handler resumes correctly regardless of which button is chosen.
func (v *View) showFormatSelector(
	title, msg string,
	jsonFn, kvFn func(),
	restoreCapture func(*tcell.EventKey) *tcell.EventKey,
) {
	modal := primitive.NewModal(title, []string{"JSON", "Key-Value", "Cancel"}, tcell.ColorDarkOliveGreen)
	modal.SetText(msg)
	modal.SetDoneFunc(func(_ int, label string) {
		v.Layout.Pages.RemovePage(pageFormatSel)
		v.Layout.Container.SetInputCapture(restoreCapture)
		switch label {
		case "JSON":
			jsonFn()
		case "Key-Value":
			kvFn()
		default:
			v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		}
	})

	// Neutral capture while the modal is open — prevents the underlying view's
	// shortcuts ('P', 'U', 'b', '/', etc.) from firing through the modal.
	v.Layout.Container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	v.Layout.Pages.AddPage(pageFormatSel, modal.Container(), true, true)
	v.Layout.Container.SetFocus(modal.Primitive())
}

// promptSecretFormat is called after the user types a new secret name (ctrl-n flow).
func (v *View) promptSecretFormat(name string) {
	basePath := v.components.SecretsTable.Props.SelectedPath
	fullPath := fmt.Sprintf("%s%s", basePath, name)
	mount := v.state.SelectedMount

	v.showFormatSelector(
		"New Secret",
		fmt.Sprintf("Choose format for: %s", fullPath),
		func() { v.openNewSecretJSONEditor(mount, fullPath) },
		func() {
			v.showKVForm(mount, fullPath, nil, false,
				func() {
					v.state.Lock()
					v.state.SelectedPath = fullPath
					v.state.Unlock()
					v.SecretObject(mount, fullPath)
				},
				v.InputSecrets,
			)
		},
		v.InputSecrets,
	)
}

// openNewSecretJSONEditor creates an empty secret and opens the JSON editor.
func (v *View) openNewSecretJSONEditor(mount, path string) {
	go func() {
		if err := v.Client.CreateNewSecret(mount, path); err != nil {
			v.Layout.Container.QueueUpdateDraw(func() { v.handleError(err.Error()) })
			return
		}
		v.Layout.Container.QueueUpdateDraw(func() {
			v.state.Lock()
			v.state.SelectedPath = path
			v.state.Unlock()
			v.SecretObject(mount, path)
			v.components.Commands.Update(component.SecretsObjectPatchCommands)
			v.components.SecretObjTable.SetEditable(true)
			v.components.TogglesInfo.Props.Editable = true
			v.components.SecretObjTable.Props.Update = "UPDATE"
			v.components.TogglesInfo.Render()
			v.components.SecretObjTable.ToggleView()
			v.components.SecretObjTable.TextView.ScrollToBeginning()
			v.Layout.Container.SetFocus(v.components.SecretObjTable.TextArea.Primitive())
		})
	}()
}

// showKVForm displays a centered form for entering one or more key-value pairs.
//
//   - initialData pre-populates the form; nil or empty starts with one blank pair.
//   - patch selects PATCH vs full-write semantics when calling the API.
//   - onSuccess is called on the event goroutine after the write succeeds.
//   - restoreCapture is reinstated when the form is dismissed (cancel or save).
func (v *View) showKVForm(
	mount, path string,
	initialData map[string]interface{},
	patch bool,
	onSuccess func(),
	restoreCapture func(*tcell.EventKey) *tcell.EventKey,
) {
	type kvPair struct {
		key   *tview.InputField
		value *tview.InputField
	}

	var (
		pairs []*kvPair
		form  *tview.Form
	)

	addPair := func(k, val string) {
		idx := len(pairs) + 1
		p := &kvPair{
			key: tview.NewInputField().
				SetLabel(fmt.Sprintf("Key %d:   ", idx)).
				SetFieldWidth(28).
				SetText(k),
			value: tview.NewInputField().
				SetLabel(fmt.Sprintf("Value %d: ", idx)).
				SetFieldWidth(28).
				SetText(val),
		}
		pairs = append(pairs, p)
		if form != nil {
			form.AddFormItem(p.key)
			form.AddFormItem(p.value)
		}
	}

	dismiss := func() {
		v.Layout.Pages.RemovePage(pageKVForm)
		v.Layout.Container.SetInputCapture(restoreCapture)
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
	}

	saveForm := func() {
		data := make(map[string]interface{})
		for _, p := range pairs {
			k := p.key.GetText()
			if k != "" {
				data[k] = p.value.GetText()
			}
		}
		if len(data) == 0 {
			return
		}
		v.Layout.Pages.RemovePage(pageKVForm)
		v.Layout.Container.SetInputCapture(restoreCapture)
		go func() {
			if err := v.Client.UpdateSecretObjectKV2(mount, path, patch, data); err != nil {
				v.Layout.Container.QueueUpdateDraw(func() { v.handleError(err.Error()) })
				return
			}
			v.Layout.Container.QueueUpdateDraw(onSuccess)
		}()
	}

	form = tview.NewForm()

	if len(initialData) > 0 {
		for k, val := range initialData {
			addPair(k, fmt.Sprintf("%v", val))
		}
	} else {
		addPair("", "")
	}

	form.AddButton("Add Key", func() {
		addPair("", "")
		v.Draw()
	})
	form.AddButton("Save", saveForm)
	form.AddButton("Cancel", dismiss)

	form.SetTitle(" Key-Value Secret ").SetBorder(true)
	form.SetBorderColor(tcell.ColorDarkOliveGreen)
	form.SetButtonsAlign(tview.AlignCenter)

	centered := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 0, 2, true).
			AddItem(nil, 0, 1, false), 65, 0, true).
		AddItem(nil, 0, 1, false)

	v.Layout.Container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event == nil {
			return nil
		}
		if event.Key() == tcell.KeyEsc {
			dismiss()
			return nil
		}
		return event
	})

	v.Layout.Pages.AddPage(pageKVForm, centered, true, true)
	v.Layout.Container.SetFocus(form)
}
