package component

import (
	"sync"

	"github.com/dkyanakiev/vaul7y/internal/models"
	primitive "github.com/dkyanakiev/vaul7y/tui/primitives"
	"github.com/dkyanakiev/vaul7y/tui/styles"
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
	TextArea TextArea
	Props    *PolicyAclTableProps
	Flex     *tview.Flex
	editable bool
	wordWrap bool
	editMu   sync.RWMutex

	slot *tview.Flex
}

const ScrollHalfPage = 15

func (p *PolicyAclTable) ToggleWrap() {
	p.wordWrap = !p.wordWrap
	wrap := p.wordWrap
	p.TextView.ModifyPrimitive(func(tv *tview.TextView) {
		tv.SetWrap(wrap)
		tv.SetWordWrap(wrap)
	})
}

func (p *PolicyAclTable) ScrollDown(lines int) {
	p.TextView.ModifyPrimitive(func(tv *tview.TextView) {
		row, col := tv.GetScrollOffset()
		tv.ScrollTo(row+lines, col)
	})
}

func (p *PolicyAclTable) ScrollUp(lines int) {
	p.TextView.ModifyPrimitive(func(tv *tview.TextView) {
		row, col := tv.GetScrollOffset()
		if newRow := row - lines; newRow > 0 {
			tv.ScrollTo(newRow, col)
		} else {
			tv.ScrollTo(0, col)
		}
	})
}

func (p *PolicyAclTable) IsEditable() bool {
	p.editMu.RLock()
	defer p.editMu.RUnlock()
	return p.editable
}

func (p *PolicyAclTable) SetEditable(val bool) {
	p.editMu.Lock()
	defer p.editMu.Unlock()
	p.editable = val
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
	t.ModifyPrimitive(func(tv *tview.TextView) {
		tv.SetWrap(true)
		tv.SetWordWrap(true)
		tv.SetScrollable(true)
	})

	ta := primitive.NewTextArea()
	ta.SetBorder(true)
	ta.SetBorderColor(styles.TcellColorStandard)

	flex := tview.NewFlex().
		AddItem(tview.NewBox(), 0, 1, false)

	pt := &PolicyAclTable{
		Flex:     flex,
		TextView: t,
		TextArea: ta,
		Props:    &PolicyAclTableProps{},
		wordWrap: true,
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

func (p *PolicyAclTable) RenderEditArea() {
	p.reset()
	p.TextArea.SetTitle(p.Props.SelectedPolicyName + " [EDITING]")
	p.TextArea.SetText(p.Props.SelectedPolicyACL, true)
	p.slot.AddItem(p.TextArea.Primitive(), 0, 1, true)
}
