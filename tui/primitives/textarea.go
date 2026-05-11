package primitives

import (
	"github.com/atotto/clipboard"
	"github.com/dkyanakiev/vaul7y/tui/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TextArea struct {
	primitive *tview.TextArea
}

func NewTextArea() *TextArea {
	t := tview.NewTextArea()
	box := tview.NewBox()
	t.Box = box
	t.Box.SetTitleColor(styles.TcellColorHighlighPrimary)
	t.Box.SetBorderColor(styles.TcellColorStandard)
	t.Box.SetBorder(true)

	// Wire Ctrl+Q (copy) and Ctrl+V (paste) to the system clipboard so that
	// paste inserts the whole string at once instead of character by character.
	t.SetClipboard(
		func(text string) { clipboard.WriteAll(text) },
		func() string {
			text, _ := clipboard.ReadAll()
			return text
		},
	)

	return &TextArea{
		primitive: t,
	}
}

func (t *TextArea) Primitive() tview.Primitive {
	return t.primitive
}

func (t *TextArea) SetText(text string, cursorAtEnd bool) *tview.TextArea {
	return t.primitive.SetText(text, cursorAtEnd)
}

func (t *TextArea) SetChangedFunc(handler func()) *tview.TextArea {
	return t.primitive.SetChangedFunc(handler)
}

func (t *TextArea) SetBorder(wrap bool) {
	t.primitive.Box.SetBorder(wrap)

}

func (t *TextArea) SetTitle(title string) {
	t.primitive.Box.SetTitle(title)
}

func (t *TextArea) SetBorderColor(color tcell.Color) {
	t.primitive.Box.SetBorderColor(color)
}

func (t *TextArea) GetText() string {
	return t.primitive.GetText()
}
