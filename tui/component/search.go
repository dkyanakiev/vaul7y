package component

import (
	"fmt"

	primitive "github.com/dkyanakiev/vaulty/tui/primitives"
	"github.com/rivo/tview"
)

const searchPlaceholder = "(hit esc to leave the filter)"

type SearchField struct {
	InputField InputField
	Props      *SearchFieldProps
	slot       *tview.Flex
}

type SearchFieldProps struct {
	DoneFunc    SetDoneFunc
	ChangedFunc func(text string)
}

func NewSearchField(label string) *SearchField {
	sf := &SearchField{}
	sf.Props = &SearchFieldProps{}
	label = fmt.Sprintf("%s ", label)
	sf.InputField = primitive.NewInputField(label, searchPlaceholder)
	return sf
}

func (s *SearchField) Render() error {
	if s.Props.DoneFunc == nil || s.Props.ChangedFunc == nil {
		return ErrComponentPropsNotSet
	}

	if s.slot == nil {
		return ErrComponentNotBound
	}

	s.InputField.SetDoneFunc(s.Props.DoneFunc)
	s.InputField.SetChangedFunc(s.Props.ChangedFunc)
	s.slot.AddItem(s.InputField.Primitive(), 0, 2, false)
	return nil
}

func (s *SearchField) Bind(slot *tview.Flex) {
	s.slot = slot
}
