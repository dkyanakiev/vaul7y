package primitives

import "github.com/rivo/tview"

type Box struct {
	*tview.Box
}

func NewBox() *Box {
	b := tview.NewBox()

	return &Box{Box: b}
}

func (b *Box) Primitive() tview.Primitive {
	return b.Box
}
