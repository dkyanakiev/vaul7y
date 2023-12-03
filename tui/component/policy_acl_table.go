package component

import (
	"github.com/dkyanakiev/vaulty/internal/models"
	primitive "github.com/dkyanakiev/vaulty/tui/primitives"
	"github.com/dkyanakiev/vaulty/tui/styles"
	"github.com/rivo/tview"
)

var (
	PolicyAclTableHeaders = []string{
		"ACL",
	}
)

type SelectPolicyACLFunc func(policyName string)

type PolicyAclTable struct {
	TextView TextView
	Props    *PolicyAclTableProps
	Flex     *tview.Flex

	slot *tview.Flex
}

type PolicyAclTableProps struct {
	SelectedPolicyName string
	// TODO: Might use data?
	SelectedPolicyACL string
	SelectPath        SelectPolicyFunc
	HandleNoResources models.HandlerFunc

	Data      []string
	Namespace string
}

func NewPolicyAclTable() *PolicyAclTable {
	t := primitive.NewTextView(1)
	t.SetTextAlign(tview.AlignLeft)
	t.SetBorderColor(styles.TcellColorStandard)
	t.SetBorder(true)

	flex := tview.NewFlex().
		//(t, 0, 1, true).
		AddItem(tview.NewBox(), 0, 1, false)

	pt := &PolicyAclTable{
		Flex:     flex,
		TextView: t,
		Props:    &PolicyAclTableProps{},
	}

	return pt
}

func (p *PolicyAclTable) Bind(slot *tview.Flex) {
	p.slot = slot
}

func (p *PolicyAclTable) reset() {
	p.slot.Clear()
	p.TextView.Clear()
}

func (p *PolicyAclTable) Render() error {
	p.reset()
	//p.Table.RenderHeader(PolicyAclTableHeaders)

	if p.Props.SelectedPolicyACL == "" {
		p.Props.HandleNoResources(
			"%sCant read ACL policy \n%s\\(╯°□°)╯︵ ┻━┻",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)
		return nil
	}

	p.renderACL()
	p.slot.AddItem(p.TextView.Primitive(), 0, 1, false)
	return nil
}

func (p *PolicyAclTable) renderACL() {
	p.TextView.SetTitle(p.Props.SelectedPolicyName)
	p.TextView.SetText(p.Props.SelectedPolicyACL)
}
