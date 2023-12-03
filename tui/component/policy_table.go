package component

import (
	"github.com/dkyanakiev/vaulty/internal/models"
	primitive "github.com/dkyanakiev/vaulty/tui/primitives"
	"github.com/dkyanakiev/vaulty/tui/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	PolicyTableTitle = "Vault ACL Policy List"
)

var (
	PolicyTableHeaderJobs = []string{
		"PolicyName",
	}
)

type SelectPolicyFunc func(policyName string)

type PolicyTable struct {
	Table Table
	Props *PolicyTableProps

	slot *tview.Flex
}

type PolicyTableProps struct {
	SelectedPolicyName string
	SelectPath         SelectPolicyFunc
	HandleNoResources  models.HandlerFunc

	Data      []string
	Namespace string
}

func NewPolicyTable() *PolicyTable {
	t := primitive.NewTable()
	pt := &PolicyTable{
		Table: t,
		Props: &PolicyTableProps{},
	}

	return pt
}

func (p *PolicyTable) Bind(slot *tview.Flex) {
	p.slot = slot
}

func (p *PolicyTable) reset() {
	p.slot.Clear()
	p.Table.Clear()
}

func (p *PolicyTable) Render() error {
	p.reset()

	p.Table.SetTitle("Vault ACL Policies")

	if len(p.Props.Data) == 0 {
		p.Props.HandleNoResources(
			"%sNo policy found\n%s\\(╯°□°)╯︵ ┻━┻",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)
		return nil
	}

	p.Table.SetSelectedFunc(p.policySelected)
	p.Table.RenderHeader(PolicyTableHeaderJobs)
	p.renderRows()

	p.slot.AddItem(p.Table.Primitive(), 0, 1, false)

	return nil

}

func (p *PolicyTable) GetIDForSelection() string {
	row, _ := p.Table.GetSelection()
	return p.Table.GetCellContent(row, 0)
}

func (p *PolicyTable) policySelected(row, _ int) {
	path := p.Table.GetCellContent(row, 0)
	p.Props.SelectedPolicyName = path
}

func (p *PolicyTable) renderRows() {

	for i, policy := range p.Props.Data {
		row := []string{
			policy,
		}
		index := i + 1
		c := tcell.ColorYellow

		p.Table.RenderRow(row, index, c)
	}
}
