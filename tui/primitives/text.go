package primitives

import (
	"github.com/dkyanakiev/vaulty/tui/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TextView struct {
	primitive *tview.TextView
}

func NewTextView(align int) *TextView {
	t := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(align)
	box := tview.NewBox()
	t.Box = box
	t.Box.SetBorderPadding(0, 0, 1, 1)
	t.Box.SetTitleColor(styles.TcellColorHighlighPrimary)
	t.Box.SetBorderColor(styles.TcellColorStandard)
	t.Box.SetBorder(false)
	return &TextView{primitive: t}
}

func (t *TextView) Primitive() tview.Primitive {
	return t.primitive
}

func (t *TextView) ModifyPrimitive(f func(t *tview.TextView)) {
	f(t.primitive)
}

func (t *TextView) SetText(text string) *tview.TextView {
	return t.primitive.SetText(text)
}

func (t *TextView) GetText(wrap bool) string {
	return t.primitive.GetText(wrap)
}

func (t *TextView) SetTitle(title string) {
	t.primitive.Box.SetTitle(title)
}

func (t *TextView) SetBorderColor(color tcell.Color) {
	t.primitive.Box.SetBorderColor(color)
}

func (t *TextView) SetBorder(wrap bool) {
	t.primitive.Box.SetBorder(wrap)
}

func (t *TextView) ScrollToBeginning() *tview.TextView {
	return t.primitive.ScrollToBeginning()
}

func (t *TextView) ScrollToEnd() *tview.TextView {
	return t.primitive.ScrollToEnd()
}

func (t *TextView) Clear() *tview.TextView {
	return t.primitive.Clear()
}

func (t *TextView) Highlight(regionIDs ...string) *tview.TextView {
	return t.primitive.Highlight(regionIDs...)
}

func (t *TextView) SetTextAlign(align int) *tview.TextView {
	return t.primitive.SetTextAlign(align)
}

func (t *TextView) Blur() {
	t.primitive.Blur()
}
