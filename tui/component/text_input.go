package component

import (
	"github.com/dkyanakiev/vaulty/tui/primitives"
	"github.com/rivo/tview"
)

const textPlaceholder = "(hit enter or esc to leave)"

type TextInfoInput struct {
	InputField InputField
	Props      *TextInfoInputProps
	slot       *tview.Flex
}

type TextInfoInputProps struct {
	DoneFunc SetDoneFunc
}

func NewTextInfoInput() *TextInfoInput {
	ti := &TextInfoInput{}
	ti.Props = &TextInfoInputProps{}

	in := primitives.NewInputField("name: ", textPlaceholder)

	ti.InputField = in

	return ti
}

func (ti *TextInfoInput) Render() error {
	if ti.Props.DoneFunc == nil {
		return ErrComponentPropsNotSet
	}

	if ti.slot == nil {
		return ErrComponentNotBound
	}
	ti.InputField.SetDoneFunc(ti.Props.DoneFunc)
	ti.InputField.SetDoneFunc(ti.Props.DoneFunc)
	ti.slot.AddItem(ti.InputField.Primitive(), 0, 2, false)
	return nil
}

func (ti *TextInfoInput) Bind(slot *tview.Flex) {
	ti.slot = slot
}
