package component

import (
	"fmt"

	"github.com/dkyanakiev/vaul7y/internal/models"
	primitive "github.com/dkyanakiev/vaul7y/tui/primitives"
	"github.com/dkyanakiev/vaul7y/tui/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	SearchResultsTableHeader = []string{
		SearchMatchPath,
		SearchMatchKey,
		SearchMatchType,
	}
)

type SearchResultsTable struct {
	Table Table
	Props *SearchResultsTableProps

	slot *tview.Flex
}

type SearchResultsTableProps struct {
	SelectedMount     string
	Term              string
	Running           bool
	Stats             models.SearchStats
	HandleNoResources models.HandlerFunc

	Data []models.SearchMatch
}

func NewSearchResultsTable() *SearchResultsTable {
	t := primitive.NewTable()

	return &SearchResultsTable{
		Table: t,
		Props: &SearchResultsTableProps{},
	}
}

func (s *SearchResultsTable) Bind(slot *tview.Flex) {
	s.slot = slot
}

func (s *SearchResultsTable) reset() {
	s.slot.Clear()
	s.Table.Clear()
}

// GetPathForSelection returns the secret path of the currently selected row.
func (s *SearchResultsTable) GetPathForSelection() string {
	row, _ := s.Table.GetSelection()
	return s.Table.GetCellContent(row, 0)
}

func (s *SearchResultsTable) Render() error {
	s.reset()
	s.Table.SetTitle("%s", s.titleText())

	if len(s.Props.Data) == 0 {
		if s.Props.Running {
			s.Props.HandleNoResources(
				"%ssearching…",
				styles.HighlightPrimaryTag,
			)
		} else {
			s.Props.HandleNoResources(
				"%sno matches found\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
				styles.HighlightPrimaryTag,
				styles.HighlightSecondaryTag,
			)
		}
		return nil
	}

	s.Table.RenderHeader(SearchResultsTableHeader)
	s.renderRows()

	s.slot.AddItem(s.Table.Primitive(), 0, 1, false)
	return nil
}

func (s *SearchResultsTable) titleText() string {
	status := "done"
	if s.Props.Running {
		status = "searching…"
	}

	title := fmt.Sprintf("Search %q in %s — %d matches (%s)",
		s.Props.Term, s.Props.SelectedMount, len(s.Props.Data), status)

	skipped := s.Props.Stats.ListDenied + s.Props.Stats.ReadDenied
	if skipped > 0 {
		title += fmt.Sprintf(" — %d skipped (no access)", skipped)
	}
	if s.Props.Stats.Truncated {
		title += " — result limit reached"
	}
	return title
}

func (s *SearchResultsTable) renderRows() {
	for i, m := range s.Props.Data {
		key := m.Key
		if key == "" {
			key = "-"
		}
		row := []string{
			m.Path,
			key,
			m.Type,
		}
		index := i + 1
		c := tcell.ColorYellow

		s.Table.RenderRow(row, index, c)
	}
}
