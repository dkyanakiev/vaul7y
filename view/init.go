package view

import (
	"fmt"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/dkyanakiev/vaulty/models"
	"github.com/dkyanakiev/vaulty/styles"
	"github.com/gdamore/tcell/v2"
)

func (v *View) Init(version string) {
	// ClusterInfo
	v.components.VaultInfo.Props.Info = fmt.Sprintf(
		"%sAddress: %s %s\n%sVersion:%s %s\n%sNamespace: %s %s",
		styles.HighlightSecondaryTag,
		v.state.VaultAddress,
		styles.StandardColorTag,
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		version,
		styles.HighlightSecondaryTag,
		v.state.Namespace,
		styles.StandardColorTag,
	)

	v.components.VaultInfo.Bind(v.Layout.Elements.ClusterInfo)
	v.components.VaultInfo.Render()
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
	}

	v.components.Search.Props.ChangedFunc = func(text string) {
		v.FilterText = text
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
		v.logger.Debug().Msgf("Info page removed, Active page is: ", v.state.Elements.TableMain)
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		// v.GoBack()
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

	// Set initial view to jobs
	v.Mounts()
}
