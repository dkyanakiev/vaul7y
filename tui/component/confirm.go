package component

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	primitive "github.com/dkyanakiev/vaul7y/tui/primitives"
)

const PageNameConfirm = "confirm"

type Confirm struct {
	Modal Modal
	Props *ConfirmProps
	pages *tview.Pages
}

type ConfirmProps struct {
	Done DoneModalFunc
}

func NewConfirm() *Confirm {
	buttons := []string{"Yes", "Cancel"}
	modal := primitive.NewModal("Confirm", buttons, tcell.ColorDarkOliveGreen)
	return &Confirm{
		Modal: modal,
		Props: &ConfirmProps{},
	}
}

func (c *Confirm) Render(msg string) error {
	if c.Props.Done == nil {
		return ErrComponentPropsNotSet
	}
	if c.pages == nil {
		return ErrComponentNotBound
	}
	c.Modal.SetDoneFunc(c.Props.Done)
	c.Modal.SetText(msg)
	c.pages.AddPage(PageNameConfirm, c.Modal.Container(), true, true)
	return nil
}

// FocusPrimitive returns the actual modal primitive that must receive focus.
// Call v.Layout.Container.SetFocus(confirm.FocusPrimitive()) immediately after Render.
func (c *Confirm) FocusPrimitive() tview.Primitive {
	return c.Modal.Primitive()
}

func (c *Confirm) Bind(pages *tview.Pages) {
	c.pages = pages
}
