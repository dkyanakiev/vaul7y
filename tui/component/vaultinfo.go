package component

import (
	"github.com/dkyanakiev/vaulty/tui/primitives"
	"github.com/rivo/tview"
)

type VaultInfo struct {
	TextView TextView
	Props    *VaultInfoProps

	slot *tview.Flex
}

type VaultInfoProps struct {
	Info string
}

func NewVaultInfo() *VaultInfo {
	return &VaultInfo{
		TextView: primitives.NewTextView(tview.AlignLeft),
		Props:    &VaultInfoProps{},
	}
}

func (c *VaultInfo) Render() error {
	if c.slot == nil {
		return ErrComponentNotBound
	}

	c.TextView.SetText(c.Props.Info)
	c.slot.AddItem(c.TextView.Primitive(), 0, 1, false)

	return nil
}

func (c *VaultInfo) Bind(slot *tview.Flex) {
	c.slot = slot
}
