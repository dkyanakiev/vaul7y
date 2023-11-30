package component

import (
	"fmt"
	"sort"
	"time"

	"github.com/dkyanakiev/vaulty/models"
	primitive "github.com/dkyanakiev/vaulty/primitives"
	"github.com/dkyanakiev/vaulty/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

const (
	TableTitleMounts = "Secret Engines"
)

var (
	TableHeaderJobs = []string{
		MountPath,
		MountName,
		MountDescription,
		MountType,
	}
)

type SelectMountPathFunc func(mountPath string)

type MountsTable struct {
	Table  Table
	Props  *MountsTableProps
	Logger *zerolog.Logger

	slot *tview.Flex
}

type MountsTableProps struct {
	SelectedMount     string
	SelectPath        SelectMountPathFunc
	HandleNoResources models.HandlerFunc
	// Data              map[string]*models.MountOutput

	Data      map[string]*models.MountOutput
	Namespace string
}

func NewMountsTable() *MountsTable {
	t := primitive.NewTable()

	jt := &MountsTable{
		Table: t,
		Props: &MountsTableProps{},
	}

	return jt
}

func (m *MountsTable) Bind(slot *tview.Flex) {
	m.slot = slot
}

func (m *MountsTable) Render() error {

	m.reset()
	m.Table.SetTitle("%s (%s)", TableTitleMounts, "Default")

	if len(m.Props.Data) == 0 {
		m.Props.HandleNoResources(
			"%sno mounts available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)

		return nil
	}

	m.Table.SetSelectedFunc(m.mountSelected)
	m.Table.RenderHeader(TableHeaderJobs)
	m.renderRows()

	m.slot.AddItem(m.Table.Primitive(), 0, 1, false)
	return nil
}

func (m *MountsTable) GetIDForSelection() string {
	row, _ := m.Table.GetSelection()
	return m.Table.GetCellContent(row, 0)
}

func (m *MountsTable) validate() error {
	// TODO: Revisid validation
	if m.Props.SelectedMount == "" || m.Props.HandleNoResources == nil {
		m.Logger.Err(ErrComponentPropsNotSet).Msgf("Random error: %s", ErrComponentPropsNotSet)
	}

	if m.slot == nil {
		return ErrComponentNotBound
	}

	return nil
}

func (m *MountsTable) reset() {
	m.slot.Clear()
	m.Table.Clear()
}

func (m *MountsTable) mountSelected(row, _ int) {
	mountPath := m.Table.GetCellContent(row, 0)
	m.Props.SelectedMount = mountPath
}

func (m *MountsTable) renderRows() {
	index := 0

	keys := make([]string, 0, len(m.Props.Data))
	for k := range m.Props.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		mount := m.Props.Data[k]
		if mount.Type == "kv" {
			row := []string{
				k,
				mount.Type,
				mount.Description,
				mount.RunningVersion,
			}
			index = index + 1

			c := m.cellColor(mount.Type)

			m.Table.RenderRow(row, index, c)
		}
	}

}

func (m *MountsTable) cellColor(mountType string) tcell.Color {
	c := tcell.ColorWhite
	// Setup splits based on type
	switch mountType {
	case models.MountTypeSystem:
		c = styles.TcellColorAttention
	case models.MountTypeCubbyhole:
		c = tcell.ColorYellow
	case models.MountTypeIdentity:
		c = tcell.ColorRed
	case models.MountTypeKV:
		c = tcell.ColorGreenYellow
	case models.MountTypePki:
		c = tcell.ColorBlue
	}

	return c
}

func formatTimeSince(since time.Duration) string {
	if since.Seconds() < 60 {
		return fmt.Sprintf("%.0fs", since.Seconds())
	}

	if since.Minutes() < 60 {
		return fmt.Sprintf("%.0fm", since.Minutes())
	}

	if since.Hours() < 60 {
		return fmt.Sprintf("%.0fh", since.Hours())
	}

	return fmt.Sprintf("%.0fd", (since.Hours() / 24))
}
