package view

import (
	"fmt"

	"github.com/dkyanakiev/vaulty/internal/models"
	"github.com/dkyanakiev/vaulty/tui/component"
	"github.com/dkyanakiev/vaulty/tui/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) Init(version string) {
	// ClusterInfo
	v.state.Version = version
	v.components.VaultInfo.Props.Info = fmt.Sprintf(
		"%sAddress: %s%s %s\n%sVersion:%s %s\n%sNamespace: %s%s %s\n",
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		v.state.VaultAddress,
		styles.StandardColorTag,
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		v.state.VaultVersion,
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		v.state.Namespace,
		styles.StandardColorTag,
	)

	v.components.VaultInfo.Bind(v.Layout.Elements.ClusterInfo)
	v.components.VaultInfo.InitialRender()
	// TogglesInfo
	v.components.TogglesInfo.Bind(v.Layout.Elements.ClusterInfo)
	v.components.TogglesInfo.InitialRender()

	// Logo
	v.components.Logo.Bind(v.Layout.Header.SlotLogo)
	v.components.Logo.Render()

	// Commands
	v.components.Commands.Bind(v.Layout.Header.SlotCmd)
	v.components.Commands.Render()

	// MountsTable
	v.components.MountsTable.Bind(v.Layout.Body)
	v.components.MountsTable.Props.HandleNoResources = v.handleNoResources

	// PolicyView
	v.components.PolicyTable.Bind(v.Layout.Body)
	v.components.PolicyTable.Props.HandleNoResources = v.handleNoResources

	// PolicyAclView
	v.components.PolicyAclTable.Bind(v.Layout.Body)
	v.components.PolicyAclTable.Props.HandleNoResources = v.handleNoResources

	// SecretView
	v.components.SecretsTable.Bind(v.Layout.Body)
	v.components.SecretsTable.Props.HandleNoResources = v.handleNoResources

	// SecretObjectView
	v.components.SecretObjTable.Bind(v.Layout.Body)
	v.components.SecretObjTable.Props.HandleNoResources = v.handleNoResources

	// NamespaceTable
	v.components.NamespaceTable.Bind(v.Layout.Body)
	v.components.NamespaceTable.Props.HandleNoResources = v.handleNoResources

	// Selections
	v.components.Selections.Bind(v.Layout.Elements.Dropdowns)
	// v.components.Selections.Render()
	v.components.Selections.Init()
	// v.components.Selections.Namespace.SetSelectedFunc(func(option string, optionIndex int) {
	// 	v.Layout.Container.SetFocus(v.state.Elements.TableMain)
	// 	_, t := v.components.Selections.Namespace.Primitive().(*tview.DropDown).GetCurrentOption()
	// 	v.state.SelectedNamespace = t
	// 	v.logger.Debug().Msgf("Selected namespace: %s", v.state.SelectedNamespace)
	// 	v.logger.Debug().Msgf("New Namespace: %s", fmt.Sprintf("%s/%s", v.state.Namespace, v.state.SelectedNamespace))
	// 	v.state.Namespaces = v.Client.ChangeNamespace(fmt.Sprintf("%s/%s", v.state.Namespace, v.state.SelectedNamespace))
	// 	v.logger.Debug().Msgf("New list: %v", v.state.Namespaces)
	// 	v.components.Selections.Namespace.SetOptions(v.state.Namespaces, func(text string, index int) {
	// 		// Add your logic here if you want to do something when an option is selected
	// 		v.state.SelectedNamespace = text

	// 	})
	// 	v.UpdateVaultInfo()
	// 	// v.components.Selections.Render(list)
	// })

	// Set the selected function
	v.components.Selections.Namespace.SetSelectedFunc(func(option string, optionIndex int) {
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		_, t := v.components.Selections.Namespace.Primitive().(*tview.DropDown).GetCurrentOption()
		v.state.SelectedNamespace = t
		v.logger.Debug().Msgf("New Namespace: %s", fmt.Sprintf("%s/%s", v.state.RootNamespace, v.state.SelectedNamespace))
		// Call the function to update the namespaces and the dropdown options
		v.updateNamespaces()
		v.UpdateVaultInfo()
	})

	// v.components.Selections.Props.DoneFunc = func(key tcell.Key) {

	// JumpToPolicy
	// v.components.JumpToPolicy.Bind(v.Layout.Footer)
	// v.components.JumpToPolicy.Props.DoneFunc = func(key tcell.Key) {
	// 	v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
	// 	v.Layout.Footer.RemoveItem(v.components.JumpToPolicy.InputField.Primitive())
	// 	v.Layout.Container.SetFocus(v.state.Elements.TableMain)

	// 	id := v.components.JumpToPolicy.InputField.GetText()
	// 	if id != "" {
	// 		//jobID := v.components.JumpToPolicy.InputField.GetText()
	// 		//v.Allocations(jobID)
	// 	}

	// 	v.components.JumpToPolicy.InputField.SetText("")
	// 	v.state.Toggle.JumpToPolicy = false
	// }

	// SearchField
	v.components.Search.Bind(v.Layout.Footer)
	v.components.Search.Props.DoneFunc = func(key tcell.Key) {
		v.state.Toggle.Search = false
		v.components.Search.InputField.SetText("")
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.Layout.Footer.RemoveItem(v.components.Search.InputField.Primitive())
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		v.components.TogglesInfo.Render()
	}

	v.components.Search.Props.ChangedFunc = func(text string) {
		v.FilterText = text
	}

	// TextInput (New secret)
	v.components.TextInfoInput.Bind(v.Layout.Footer)
	v.components.TextInfoInput.Props.DoneFunc = func(key tcell.Key) {
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.Layout.Footer.RemoveItem(v.components.TextInfoInput.InputField.Primitive())
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		v.components.TextInfoInput.Render()

		newText := v.components.TextInfoInput.InputField.GetText()
		v.state.NewSecretName = newText
		v.CreateNewSecretObject(newText)
		v.components.TextInfoInput.InputField.SetText("")
		v.state.Toggle.TextInput = false
	}

	// Error
	v.components.Error.Bind(v.Layout.Pages)
	v.components.Error.Props.Done = func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Quit" {
			v.Layout.Container.Stop()
			return
		}

		v.Layout.Pages.RemovePage(component.PageNameError)
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		///v.GoBack()
	}

	// Info
	v.components.Info.Bind(v.Layout.Pages)
	v.components.Info.Props.Done = func(buttonIndex int, buttonLabel string) {
		v.Layout.Pages.RemovePage(component.PageNameInfo)
		v.logger.Debug().Msgf("Info page removed, Active page is: %s", v.state.Elements.TableMain.GetTitle())
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		v.GoBack()
	}

	// Warn
	v.components.Failure.Bind(v.Layout.Pages)
	v.components.Failure.Props.Done = func(buttonIndex int, buttonLabel string) {
		v.Layout.Pages.RemovePage(component.PageNameInfo)
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		v.GoBack()
	}

	v.Watcher.SubscribeHandler(models.HandleError, v.handleError)
	v.Watcher.SubscribeHandler(models.HandleFatal, v.handleFatal)

	stop := make(chan struct{})

	go v.DrawLoop(stop)
	// v.logger.Debug().Msgf("Active page is: %s", v.state.Elements.TableMain.GetTitle())
	// Set initial view to jobs
	v.Mounts()
}

func (v *View) UpdateVaultInfo() {
	// Update the component's state
	newNS := fmt.Sprintf("%s/%s", v.state.RootNamespace, v.state.SelectedNamespace)

	v.components.VaultInfo.Props.Info = fmt.Sprintf(
		"%sAddress: %s%s %s\n%sVersion:%s %s\n%sNamespace: %s%s %s\n",
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		v.state.VaultAddress,
		styles.StandardColorTag,
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		v.state.VaultVersion,
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		newNS,
		styles.StandardColorTag,
	)

	// Re-render the component
	v.components.VaultInfo.Render()
}

// Function to update the namespaces and the dropdown options
func (v *View) updateNamespaces() {
	var oldNamespace string
	if v.state.SelectedNamespace == "" {
		v.logger.Debug().Msgf("Changing namespace to: %s", fmt.Sprintf("%s/%s", v.state.RootNamespace, v.state.SelectedNamespace))
		v.state.Namespaces = v.Client.ChangeNamespace(v.state.RootNamespace)
		v.logger.Debug().Msgf("New Namespaces: %s", v.state.Namespaces)
	} else {
		oldNamespace = v.state.SelectedNamespace
		v.logger.Debug().Msgf("Changing namespace to: %s", fmt.Sprintf("%s/%s", v.state.RootNamespace, v.state.SelectedNamespace))
		v.state.Namespaces = v.Client.ChangeNamespace(fmt.Sprintf("%s/%s", v.state.RootNamespace, v.state.SelectedNamespace))
		v.logger.Debug().Msgf("New Namespaces: %s", v.state.Namespaces)
	}
	// Check if the namespaces slice is not empty
	if len(v.state.Namespaces) > 0 {
		v.components.Selections.Namespace.SetOptions(v.state.Namespaces, func(text string, index int) {
			if v.state.SelectedNamespace == "" {
				v.state.SelectedNamespace = text
			} else {
				v.state.SelectedNamespace = fmt.Sprintf("%s/%s", v.state.SelectedNamespace, text)
			}
			// Only update namespaces if the selected namespace has changed
			if oldNamespace != v.state.SelectedNamespace {
				// Use a separate goroutine for the recursive call
				go v.updateNamespaces()
			}
			v.UpdateVaultInfo()
		})
	} else {
		// Set a dummy option
		v.components.Selections.Namespace.SetOptions([]string{"No namespaces available"}, nil)
	}
}
