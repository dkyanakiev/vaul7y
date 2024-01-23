package component

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/dkyanakiev/vaulty/internal/state"
	"github.com/dkyanakiev/vaulty/tui/primitives"
	"github.com/dkyanakiev/vaulty/tui/styles"
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

func (s *Selections) Init() error {
	if s.slot == nil {
		return ErrComponentNotBound
	}

	// if len(s.state.Namespaces) != 0 {
	s.Namespace.SetOptions(s.state.Namespaces, s.Selected)
	s.Namespace.SetCurrentOption(len(s.state.SelectedNamespace) - 1)
	s.Namespace.SetSelectedFunc(s.rerender)

	// }
	s.state.Elements.DropDownNamespace = s.Namespace.Primitive().(*tview.DropDown)
	s.slot.AddItem(s.Namespace.Primitive(), 0, 1, true)

	return nil
}
func (s *Selections) Render() error {
	if s.slot == nil {
		return ErrComponentNotBound
	}

	// if len(s.state.Namespaces) != 0 {
	s.Namespace.SetOptions(s.state.Namespaces, s.Selected)

	// }
	s.state.Elements.DropDownNamespace = s.Namespace.Primitive().(*tview.DropDown)
	// s.slot.AddItem(s.Namespace.Primitive(), 0, 1, true)

	return nil
}

func (s *Selections) Refresh() {
	s.Render()
}

// func (s *Selections) Update() error {
// 	s.Namespace.SetOptions(s.state.Namespaces, s.selected)
// 	s.Namespace.SetCurrentOption(len(s.state.Namespace) - 1)
// 	s.Namespace.SetSelectedFunc(s.rerender)

// }

func (s *Selections) Bind(slot *tview.Flex) {
	s.slot = slot
}

func (s *Selections) Selected(text string, index int) {
	s.state.SelectedNamespace = text
}

func (s *Selections) rerender(text string, index int) {
	s.state.SelectedNamespace = text
}
