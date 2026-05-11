package component

import (
	"sort"

	"github.com/dkyanakiev/vaul7y/internal/models"
	primitive "github.com/dkyanakiev/vaul7y/tui/primitives"
	"github.com/dkyanakiev/vaul7y/tui/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const TableTitleAuth = "Auth Methods"

var AuthTableHeaders = []string{"Path", "Type", "Description"}

type AuthTable struct {
	Table Table
	Props *AuthTableProps

	slot *tview.Flex
}

type AuthTableProps struct {
	Data              map[string]*models.AuthMethod
	HandleNoResources models.HandlerFunc
}

func NewAuthTable() *AuthTable {
	t := primitive.NewTable()
	return &AuthTable{
		Table: t,
		Props: &AuthTableProps{},
	}
}

func (a *AuthTable) Bind(slot *tview.Flex) {
	a.slot = slot
}

func (a *AuthTable) Render() error {
	a.reset()
	a.Table.SetTitle("%s", TableTitleAuth)

	if len(a.Props.Data) == 0 {
		a.Props.HandleNoResources(
			"%sno auth methods available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)
		return nil
	}

	a.Table.RenderHeader(AuthTableHeaders)
	a.renderRows()
	a.slot.AddItem(a.Table.Primitive(), 0, 1, false)
	return nil
}

func (a *AuthTable) reset() {
	a.slot.Clear()
	a.Table.Clear()
}

func (a *AuthTable) renderRows() {
	keys := make([]string, 0, len(a.Props.Data))
	for k := range a.Props.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		m := a.Props.Data[k]
		row := []string{k, m.Type, m.Description}
		a.Table.RenderRow(row, i+1, tcell.ColorWhite)
	}
}
