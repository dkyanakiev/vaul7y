package view

import (
	"github.com/dkyanakiev/vaulty/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (v *View) PolicyACL(policyName string) {

	v.viewSwitch()
	v.Layout.Body.SetTitle(policyName)
	v.components.PolicyAclTable.TextView.Clear().ScrollToBeginning()
	v.components.Commands.Update(component.PolicyACLCommands)
	v.Layout.Container.SetInputCapture(v.inputPolicyACL)

	v.state.SelectedPolicyName = policyName
	v.components.PolicyAclTable.Props.SelectedPolicyName = policyName

	update := func() {
		v.components.PolicyAclTable.Props.SelectedPolicyACL = v.state.PolicyACL
		v.components.PolicyAclTable.Render()
		v.Draw()
	}
	v.Watcher.SubscribeToPoliciesACL(update)
	update()

	v.state.Elements.TextMain = v.components.PolicyAclTable.TextView.Primitive().(*tview.TextView)
	v.Layout.Container.SetFocus(v.components.PolicyAclTable.TextView.Primitive())

}

func (v *View) inputPolicyACL(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyEsc:
		v.VPolicy()
	}
	return event
}
