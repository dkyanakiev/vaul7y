package primitives

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputField struct {
	primitive *tview.InputField
}

func NewInputField(label, placeholder string) *InputField {
	i := tview.NewInputField()
	i.SetLabel(label)
	i.SetFieldWidth(0)
	i.SetAcceptanceFunc(tview.InputFieldMaxLength(40))
	i.SetPlaceholder(placeholder)
	i.SetBorder(true)
	i.SetFieldBackgroundColor(tcell.NewRGBColor(40, 44, 48))
	i.SetBackgroundColor(tcell.NewRGBColor(40, 44, 48))
	i.SetBorderAttributes(tcell.AttrDim)

	return &InputField{i}
}

func (i *InputField) SetDoneFunc(handler func(k tcell.Key)) {
	i.primitive.SetDoneFunc(handler)
}

func (i *InputField) SetChangedFunc(handler func(text string)) {
	i.primitive.SetChangedFunc(handler)
}

func (i *InputField) SetText(text string) {
	i.primitive.SetText(text)
}

func (i *InputField) GetText() string {
	return i.primitive.GetText()
}

func (i *InputField) SetAutocompleteFunc(callback func(currentText string) (entries []string)) {
	i.primitive.SetAutocompleteFunc(callback)
}

func (i *InputField) Primitive() tview.Primitive {
	return i.primitive
}
