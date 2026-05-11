package view

import (
	"github.com/atotto/clipboard"
	"github.com/dkyanakiev/vaul7y/tui/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) PolicyACL(policyName string) {

	v.viewSwitch()
	v.Layout.Body.SetTitle(policyName)
	v.Layout.Container.SetFocus(v.components.PolicyAclTable.TextView.Primitive())
	v.components.PolicyAclTable.TextView.Clear().ScrollToBeginning()
	v.components.Commands.Update(component.PolicyACLCommands)
	v.Layout.Container.SetInputCapture(v.InputPolicyACL)

	v.state.SelectedPolicyName = policyName
	v.components.PolicyAclTable.Props.SelectedPolicyName = policyName

	update := func() {
		v.state.RLock()
		acl := v.state.PolicyACL
		v.state.RUnlock()

		v.components.PolicyAclTable.Props.SelectedPolicyACL = acl
		v.components.PolicyAclTable.Render()
		v.Draw()
	}
	v.Watcher.SubscribeToPoliciesACL(func() { v.Layout.Container.QueueUpdateDraw(update) })
	update()

	v.state.Elements.TextMain = v.components.PolicyAclTable.TextView.Primitive().(*tview.TextView)

}

func (v *View) InputPolicyACL(event *tcell.EventKey) *tcell.EventKey {
	event = v.InputMainCommands(event)
	return v.inputPolicyACL(event)
}

func (v *View) inputPolicyACL(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	editing := v.components.PolicyAclTable.IsEditable()

	switch event.Key() {
	case tcell.KeyCtrlW:
		if editing {
			v.savePolicyEdit()
			return nil
		}
	case tcell.KeyEsc:
		if editing {
			v.components.PolicyAclTable.SetEditable(false)
			v.components.PolicyAclTable.Render()
			v.components.Commands.Update(component.PolicyACLCommands)
			v.Layout.Container.SetFocus(v.components.PolicyAclTable.TextView.Primitive())
			return nil
		}
		v.VPolicy()
	case tcell.KeyRune:
		if !editing {
			switch event.Rune() {
			case 'E', 'e':
				v.components.PolicyAclTable.SetEditable(true)
				v.components.PolicyAclTable.RenderEditArea()
				v.components.Commands.Update(component.PolicyACLEditCommands)
				v.Layout.Container.SetFocus(v.components.PolicyAclTable.TextArea.Primitive())
				return nil
			case 'c':
				content := v.components.PolicyAclTable.Props.SelectedPolicyACL
				clipboard.WriteAll(content) //nolint:errcheck
				return nil
			case 'w':
				v.components.PolicyAclTable.ToggleWrap()
				v.Draw()
				return nil
			case 'j':
				v.components.PolicyAclTable.ScrollDown(1)
				v.Draw()
				return nil
			case 'k':
				v.components.PolicyAclTable.ScrollUp(1)
				v.Draw()
				return nil
			case 'd':
				v.components.PolicyAclTable.ScrollDown(component.ScrollHalfPage)
				v.Draw()
				return nil
			case 'u':
				v.components.PolicyAclTable.ScrollUp(component.ScrollHalfPage)
				v.Draw()
				return nil
			case 'g':
				v.components.PolicyAclTable.TextView.ScrollToBeginning()
				v.Draw()
				return nil
			case 'G':
				v.components.PolicyAclTable.TextView.ScrollToEnd()
				v.Draw()
				return nil
			}
		}
	}
	return event
}

func (v *View) savePolicyEdit() {
	name := v.components.PolicyAclTable.Props.SelectedPolicyName
	rules := v.components.PolicyAclTable.TextArea.GetText()

	if err := v.Client.UpdatePolicy(name, rules); err != nil {
		v.handleError(err.Error())
		return
	}

	v.components.PolicyAclTable.SetEditable(false)
	v.components.PolicyAclTable.Render()
	v.components.Commands.Update(component.PolicyACLCommands)
	v.Layout.Container.SetFocus(v.components.PolicyAclTable.TextView.Primitive())
	v.handleInfo("Policy updated successfully")
}
