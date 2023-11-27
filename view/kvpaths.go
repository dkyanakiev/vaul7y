package view

import "github.com/rivo/tview"

func (v *View) KVPaths() {
	v.viewSwitch()

	v.Layout.Body.SetTitle("Secret Mounts")
	// TODO
	table := v.components.MountsTable

	table.Render()
	v.Draw()

	v.state.Elements.TableMain = v.components.MountsTable.Table.Primitive().(*tview.Table)
}
