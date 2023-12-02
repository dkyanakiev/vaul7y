package component

import (
	"fmt"
	"strconv"

	"github.com/dkyanakiev/vaulty/internal/models"
	primitive "github.com/dkyanakiev/vaulty/tui/primitives"
	"github.com/dkyanakiev/vaulty/tui/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	SecretsTableMount = "Secrets Explorer"
)

var (
	SecretsTableHeaderJobs = []string{
		SecretPath,
		SecretObject,
	}
)

type SelectSecretsPathFunc func(mountPath string)

type SecretsTable struct {
	Table Table
	Props *SecretsTableProps

	slot *tview.Flex
}

type SecretsTableProps struct {
	SelectedMount     string
	SelectedObject    string
	SelectedPath      string
	SelectPath        SelectMountPathFunc
	HandleNoResources models.HandlerFunc
	// Data              map[string]*models.MountOutput

	Data      []models.SecretPath
	Namespace string
}

func NewSecretsTable() *SecretsTable {
	t := primitive.NewTable()

	st := &SecretsTable{
		Table: t,
		Props: &SecretsTableProps{},
	}

	return st
}

func (s *SecretsTable) Bind(slot *tview.Flex) {
	s.slot = slot
}

func (s *SecretsTable) reset() {
	s.slot.Clear()
	s.Table.Clear()
}

func (s *SecretsTable) pathSelected(row, _ int) {
	mountPath := s.Table.GetCellContent(row, 0)
	s.Props.SelectedMount = fmt.Sprintf("%s%s", s.Props.SelectedMount, mountPath)
}

func (s *SecretsTable) GetIDForSelection() (string, string) {
	row, _ := s.Table.GetSelection()
	name := s.Table.GetCellContent(row, 0)
	secret := s.Table.GetCellContent(row, 1)
	return name, secret
}

func (s *SecretsTable) Render() error {
	s.reset()
	fullPath := fmt.Sprintf("%s%s", s.Props.SelectedMount, s.Props.SelectedPath)
	s.Table.SetTitle("%s (%s)", TableTitleMounts, fullPath)

	if len(s.Props.Data) == 0 {
		s.Props.HandleNoResources(
			"%sno secrets available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)

		return nil
	}

	s.Table.SetSelectedFunc(s.pathSelected)
	s.Table.RenderHeader(SecretsTableHeaderJobs)
	s.renderRows()

	s.slot.AddItem(s.Table.Primitive(), 0, 1, false)
	return nil
}

func (s *SecretsTable) renderRows() {

	for i, obj := range s.Props.Data {
		row := []string{
			obj.PathName,
			strconv.FormatBool(obj.IsSecret),
		}
		index := i + 1
		c := tcell.ColorYellow

		s.Table.RenderRow(row, index, c)
	}
}
