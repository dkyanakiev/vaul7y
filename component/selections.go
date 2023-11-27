package component

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/dkyanakiev/vaulty/primitives"
	"github.com/dkyanakiev/vaulty/state"
	"github.com/dkyanakiev/vaulty/styles"
)

var (
	labelNamespaceDropdown = fmt.Sprintf("%sNamespace <s>: ▾ %s",
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
	)
)

type Selections struct {
	Namespace DropDown
	Mounts    DropDown
	//Path
	state *state.State
	slot  *tview.Flex
}

func NewSelections(state *state.State) *Selections {
	return &Selections{
		Namespace: primitives.NewDropDown(labelNamespaceDropdown),
		state:     state,
	}
}

func (s *Selections) Render() error {
	if s.slot == nil {
		return ErrComponentNotBound
	}

	//s.Namespace.SetOptions(convert(s.state.Namespace), s.selected)
	s.Namespace.SetCurrentOption(len(s.state.Namespace) - 1)
	s.Namespace.SetSelectedFunc(s.rerender)

	//s.state.Elements.DropDownNamespace = s.Namespace.Primitive().(*tview.DropDown)
	s.slot.AddItem(s.Namespace.Primitive(), 0, 1, true)

	return nil
}

func (s *Selections) Bind(slot *tview.Flex) {
	s.slot = slot
}

func (s *Selections) selected(text string, index int) {
	s.state.SelectedNamespace = text
}

func (s *Selections) rerender(text string, index int) {
	s.state.SelectedNamespace = text
}

// func convert(list []*models.Namespaces) []string {
// 	var ns []string
// 	for _, n := range list {
// 		ns = append(ns, n.Name)
// 	}
// 	return ns
// }
