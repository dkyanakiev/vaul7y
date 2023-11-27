package primitives

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Form struct {
	primitive *tview.Form
	container *tview.Flex
}

func NewForm(title string, c tcell.Color) *Form {
	f := tview.NewForm()
	f.SetTitle(title)
	f.SetTitleAlign(tview.AlignCenter)
	f.SetBackgroundColor(c)
	f.SetFieldTextColor(tcell.ColorBlack)

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(f, 10, 1, true).
			AddItem(nil, 0, 1, false), 80, 1, false).
		AddItem(nil, 0, 1, false)

	return &Form{
		primitive: f,
		container: flex,
	}
}

func (f *Form) Container() tview.Primitive {
	return f.container
}

func (f *Form) Primitive() tview.Primitive {
	return f.primitive
}
