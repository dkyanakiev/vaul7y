package component

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"

	"github.com/dkyanakiev/vaulty/internal/models"
	primitive "github.com/dkyanakiev/vaulty/tui/primitives"
	"github.com/dkyanakiev/vaulty/tui/styles"
)

const TableTitleNamespaces = "Namespaces"

var (
	TableHeaderNamespaces = []string{
		LabelName,
	}
)

type SelectNsPathFunc func(ns string)

type NamespaceTable struct {
	Table  Table
	Props  *NamespacesProps
	Logger *zerolog.Logger
	slot   *tview.Flex
}

type NamespacesProps struct {
	SelectedNamespace string
	SelectNs          SelectNsPathFunc
	HandleNoResources models.HandlerFunc
	Data              []string
}

func NewNamespaceTable() *NamespaceTable {
	t := primitive.NewTable()

	return &NamespaceTable{
		Table: t,
		Props: &NamespacesProps{},
	}
}

func (n *NamespaceTable) Bind(slot *tview.Flex) {
	slot.SetTitle("Namespaces")
	n.slot = slot
}

func (n *NamespaceTable) Render() error {
	if n.slot == nil {
		return ErrComponentNotBound
	}

	if n.Props.HandleNoResources == nil {
		return ErrComponentPropsNotSet
	}

	n.reset()
	n.Logger.Debug().Msgf("rendering namespaces: %v", n.Props.Data)
	if len(n.Props.Data) == 0 {
		n.Props.HandleNoResources(
			"%sno namespaces available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)

		return nil
	}

	n.Table.SetTitle(TableTitleNamespaces)
	n.Table.RenderHeader(TableHeaderNamespaces)
	n.Table.SetSelectedFunc(n.namespaceSelected)
	n.renderRows()

	n.slot.AddItem(n.Table.Primitive(), 0, 1, false)
	return nil
}

func (n *NamespaceTable) reset() {
	n.Table.Clear()
	n.slot.Clear()
}

func (n *NamespaceTable) renderRows() {
	index := 0
	for i, ns := range n.Props.Data {
		row := []string{
			ns,
		}

		index = i + 1
		n.Table.RenderRow(row, index, tcell.ColorWhite)
	}
}

func (n *NamespaceTable) GetIDForSelection() string {
	row, _ := n.Table.GetSelection()
	return n.Table.GetCellContent(row, 0)
}

func (n *NamespaceTable) namespaceSelected(row, _ int) {
	ns := n.Table.GetCellContent(row, 0)
	n.Props.SelectedNamespace = ns
}
