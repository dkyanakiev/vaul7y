package component

import (
	"github.com/dkyanakiev/vaulty/models"
	"github.com/dkyanakiev/vaulty/primitives"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	ErrComponentNotBound    = models.Comp("component not bound")
	ErrComponentPropsNotSet = models.Comp("component properties not set")
)

//go:generate counterfeiter . DoneModalFunc
type DoneModalFunc func(buttonIndex int, buttonLabel string)

type Primitive interface {
	Primitive() tview.Primitive
}

//go:generate counterfeiter . Table
type Table interface {
	Primitive
	SetTitle(format string, args ...interface{})
	GetCellContent(row, column int) string
	GetSelection() (row, column int)
	Clear()
	RenderHeader(data []string)
	RenderRow(data []string, index int, c tcell.Color)
	SetSelectedFunc(fn func(row, column int))
	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey)
}

//go:generate counterfeiter . TextView
type TextView interface {
	Primitive
	GetText(bool) string
	SetText(text string) *tview.TextView
	Write(data []byte) (int, error)
	Highlight(regionIDs ...string) *tview.TextView
	Clear() *tview.TextView
	ModifyPrimitive(f func(t *tview.TextView))
}

//go:generate counterfeiter . Modal
type Modal interface {
	Primitive
	SetDoneFunc(handler func(buttonIndex int, buttonLabel string))
	SetText(text string)
	SetFocus(index int)
	Container() tview.Primitive
}

type Form interface {
	Primitive
	Container() tview.Primitive
}

//go:generate counterfeiter . InputField
type InputField interface {
	Primitive
	SetDoneFunc(handler func(k tcell.Key))
	SetChangedFunc(handler func(text string))
	SetAutocompleteFunc(callback func(currentText string) (entries []string))
	SetText(text string)
	GetText() string
}

//go:generate counterfeiter . DropDown
type DropDown interface {
	Primitive
	SetOptions(options []string, selected func(text string, index int))
	SetCurrentOption(index int)
	SetSelectedFunc(selected func(text string, index int))
}

//go:generate counterfeiter . Selector
type Selector interface {
	Primitive
	GetTable() *primitives.Table
	Container() tview.Primitive
}
