package view

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dkyanakiev/vaul7y/internal/models"
	"github.com/dkyanakiev/vaul7y/tui/component"
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
	v.Layout.Container.SetFocus(v.components.SecretsTable.Table.Primitive())
	v.Layout.Container.SetInputCapture(v.InputSecrets)
	v.components.Commands.Update(component.SecretsCommands)
	v.logger.Debug().Msgf("Selected path for secret is: %v", path)

	search := v.components.Search
	v.state.Toggle.Search = false
	v.state.Filter.Object = ""

	v.components.SecretsTable.Props.SelectedMount = v.state.SelectedMount
	if path != "" {
		v.components.SecretsTable.Props.SelectedPath = fmt.Sprintf("%s%s", v.state.SelectedPath, v.state.SelectedObject)
	}

	update := func() {
		v.state.RLock()
		searching := v.state.Toggle.Search
		selectedMount := v.state.SelectedMount
		v.state.RUnlock()

		v.mutex.RLock()
		filterText := v.FilterText
		v.mutex.RUnlock()

		if searching {
			v.state.Lock()
			v.state.Filter.Object = filterText
			v.state.Unlock()
			v.components.TogglesInfo.Props.FilterText = filterText
		}

		v.state.RLock()
		data := v.filterSecrets()
		v.state.RUnlock()

		v.components.SecretsTable.Props.Data = data
		v.components.SecretsTable.Props.SelectedMount = selectedMount
		v.components.SecretsTable.Render()
		v.Draw()
		v.components.SecretsTable.Table.ScrollToTop()
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

	v.Watcher.SubscribeToSecrets(v.components.SecretsTable.Props.SelectedMount,
		v.components.SecretsTable.Props.SelectedPath, func() { v.Layout.Container.QueueUpdateDraw(update) })
	update()

	v.state.Elements.TableMain = v.components.SecretsTable.Table.Primitive().(*tview.Table)

}

func (v *View) inputSecrets(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyEsc:
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
	case tcell.KeyEnter:
		if v.components.SecretsTable.Table.Primitive().HasFocus() {
			path, secretBool := v.components.SecretsTable.GetIDForSelection()
			v.state.SelectedPath = fmt.Sprintf("%s%s", v.state.SelectedPath, path)
			if secretBool == "true" {
				v.SecretObject(v.state.SelectedMount, v.state.SelectedPath)
			} else {
				v.logger.Debug().Msgf("Running Secrets view with : %v", path)
				v.Secrets(path, secretBool)
			}
			return nil
		}
	case tcell.KeyCtrlF:
		if !v.Layout.Footer.HasFocus() {
			v.state.RLock()
			textInput := v.state.Toggle.TextInput
			v.state.RUnlock()
			if !textInput {
				v.state.Lock()
				v.state.Toggle.TextInput = true
				v.state.Toggle.RecursiveSearch = true
				v.state.Unlock()
				v.components.TextInfoInput.InputField.SetLabel("search: ")
				v.components.TextInfoInput.InputField.SetText("")
				v.TextInput()
			} else {
				v.Layout.Container.SetFocus(v.components.TextInfoInput.InputField.Primitive())
			}
			return nil
		}
	case tcell.KeyCtrlN:
		v.logger.Debug().Msgf("Running New Secret view with : %v", v.state.SelectedPath)
		if !v.Layout.Footer.HasFocus() {
			v.state.RLock()
			textInput := v.state.Toggle.TextInput
			v.state.RUnlock()
			if !textInput {
				v.state.Lock()
				v.state.Toggle.TextInput = true
				v.state.Unlock()
				v.components.TextInfoInput.InputField.SetText("")
				v.TextInput()
			} else {
				v.Layout.Container.SetFocus(v.components.TextInfoInput.InputField.Primitive())
			}
			return nil
		}
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
		case '/':
			if !v.Layout.Footer.HasFocus() {
				v.state.RLock()
				searching := v.state.Toggle.Search
				v.state.RUnlock()
				if !searching {
					v.state.Lock()
					v.state.Toggle.Search = true
					v.state.Unlock()
					v.components.Search.InputField.SetText("")
					v.Search()
				} else {
					v.Layout.Container.SetFocus(v.components.Search.InputField.Primitive())
				}
				return nil
			}
		case 'J':
			if !v.Layout.Footer.HasFocus() {
				v.state.RLock()
				textInput := v.state.Toggle.TextInput
				v.state.RUnlock()
				if !textInput {
					v.state.Lock()
					v.state.Toggle.TextInput = true
					v.state.Toggle.JumpToPath = true
					v.state.Unlock()
					v.components.TextInfoInput.InputField.SetText("")
					v.TextInput()
				} else {
					v.state.Lock()
					v.state.Toggle.JumpToPath = true
					v.state.Unlock()
					v.Layout.Container.SetFocus(v.components.TextInfoInput.InputField.Primitive())
				}
				return nil
			}
		}
	}

	return event
}

func (v *View) filterSecrets() []models.SecretPath {
	data := v.state.SecretsData
	filter := v.state.Filter.Object
	if filter != "" {
		rx, err := regexp.Compile(filter)
		if err != nil {
			return data
		}
		var result []models.SecretPath
		for _, p := range data {
			if rx.MatchString(p.PathName) {
				result = append(result, p)
			}
		}
		return result
	}

	return data
}

func trimLastElement(s string) string {
	dir, _ := filepath.Split(s)
	return strings.TrimSuffix(dir, string(filepath.Separator)) + string(filepath.Separator)
}

func (v *View) CreateNewSecretObject(newObj string) {
	v.logger.Info().Msgf("Creating new secret object for path: %v", v.components.SecretsTable.Props.SelectedPath)
	v.logger.Info().Msgf("Creating new secret object for mount: %v", v.state.SelectedMount)
	if v.state.NewSecretName != "" {
		v.logger.Debug().Msgf("New secret name is: %v", v.state.NewSecretName)
		newPath := fmt.Sprintf("%s%s", v.components.SecretsTable.Props.SelectedPath, v.state.NewSecretName)
		err := v.Client.CreateNewSecret(v.state.SelectedMount, newPath)
		if err != nil {
			v.handleError(string(err.Error()))
		} else {
			v.handleInfo("Secret path created successfully")
		}

	} else {
		v.logger.Debug().Msg("New secret name is empty")
	}
	//v.Layout.Container.SetFocus(v.components.SecretObjTable.Table.Primitive())
}
