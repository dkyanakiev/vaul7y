package primitives

import "github.com/rivo/tview"

type TextArea struct {
	*tview.TextArea
}

func NewTextArea() *TextArea {
	t := tview.NewTextArea()

	return &TextArea{TextArea: t}
}

func (t *TextArea) Primitive() tview.Primitive {
	return t.TextArea
}
