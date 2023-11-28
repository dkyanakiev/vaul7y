package component

import (
	"github.com/dkyanakiev/vaulty/primitives"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Not currently used. Might need modification to reuse this for both paths and policies.

const jumpToPlaceholder = "(hit enter or esc to leave)"

type SetDoneFunc func(key tcell.Key)

type JumpToPolicy struct {
	InputField InputField
	Props      *JumpToPolicyProps
	slot       *tview.Flex
}

type JumpToPolicyProps struct {
	DoneFunc SetDoneFunc
}

func NewJumpToPolicy() *JumpToPolicy {
	jj := &JumpToPolicy{}
	jj.Props = &JumpToPolicyProps{}

	in := primitives.NewInputField("jump: ", jumpToPlaceholder)

	in.SetAutocompleteFunc(func(currentText string) (entries []string) {
		return jj.find(currentText)
	})

	jj.InputField = in
	return jj
}

func (jj *JumpToPolicy) Render() error {
	if err := jj.validate(); err != nil {
		return err
	}

	jj.InputField.SetDoneFunc(jj.Props.DoneFunc)
	jj.slot.AddItem(jj.InputField.Primitive(), 0, 2, false)
	return nil
}

func (jj *JumpToPolicy) validate() error {
	if jj.Props.DoneFunc == nil {
		return ErrComponentPropsNotSet
	}

	if jj.slot == nil {
		return ErrComponentNotBound
	}

	return nil
}

func (jj *JumpToPolicy) Bind(slot *tview.Flex) {
	jj.slot = slot
}

func (jj *JumpToPolicy) find(text string) []string {
	result := []string{}
	if text == "" {
		return result
	}

	// for _, j := range jj.Props.Jobs {
	// 	ok := strings.Contains(j.ID, text)
	// 	if ok {
	// 		result = append(result, j.ID)
	// 	}
	// }

	return result
}
