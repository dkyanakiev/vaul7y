package component_test

import (
	"testing"

	"github.com/dkyanakiev/vaul7y/tui/component"
	"github.com/dkyanakiev/vaul7y/tui/component/componentfakes"
	"github.com/dkyanakiev/vaul7y/tui/styles"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
)

func TestPolicyACLTable_Pass(t *testing.T) {
	r := require.New(t)
	t.Run("Data to be rendered", func(t *testing.T) {
		// fakeTable := &componentfakes.FakeTable{}
		fakeTextView := &componentfakes.FakeTextView{}
		pt := component.NewPolicyAclTable()

		pt.TextView = fakeTextView
		pt.Props.SelectedPolicyName = "policy-one"
		pt.Props.SelectedPolicyACL = "path \"secret/data/path\" { capabilities = [\"read\", \"list\"] }"

		pt.Props.SelectPath = func(id string) {}
		pt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		pt.Bind(slot)

		err := pt.Render()
		r.NoError(err)

		// It renders the correct text
		fakeTextView.GetTextReturns(pt.Props.SelectedPolicyACL)
		fakeTextView.GetTextReturnsOnCall(0, pt.Props.SelectedPolicyACL)
		renderedText := fakeTextView.GetText(true)

		r.Equal("path \"secret/data/path\" { capabilities = [\"read\", \"list\"] }", renderedText)

	})

	t.Run("No data to be rendered", func(t *testing.T) {
		fakeTextView := &componentfakes.FakeTextView{}
		pt := component.NewPolicyAclTable()

		pt.TextView = fakeTextView
		pt.Props.SelectedPolicyName = "policy-one"
		pt.Props.SelectedPolicyACL = ""

		var handleNoResourcesCalled bool
		pt.Props.HandleNoResources = func(format string, args ...interface{}) {
			handleNoResourcesCalled = true

			r.Equal("%sCant read ACL policy \n%s\\(╯°□°)╯︵ ┻━┻", format)
			r.Len(args, 2)
			r.Equal(args[0], styles.HighlightPrimaryTag)
			r.Equal(args[1], styles.HighlightSecondaryTag)
		}

		slot := tview.NewFlex()
		pt.Bind(slot)

		err := pt.Render()
		r.NoError(err)

		r.True(handleNoResourcesCalled)

	})
}

func TestPolicyACLTable_Fail(t *testing.T) {
}

func TestPolicyAclTable_Scroll(t *testing.T) {
	// Helpers wire up a real tview.TextView through the FakeTextView stub so
	// the closure inside ScrollDown/ScrollUp operates on real in-memory state.
	makeTable := func(initialRow int) (*component.PolicyAclTable, *tview.TextView) {
		fakeTV := &componentfakes.FakeTextView{}
		realTV := tview.NewTextView()
		realTV.ScrollTo(initialRow, 0)
		fakeTV.ModifyPrimitiveStub = func(fn func(*tview.TextView)) { fn(realTV) }

		pt := component.NewPolicyAclTable()
		pt.TextView = fakeTV
		return pt, realTV
	}

	t.Run("ScrollDown increases row offset", func(t *testing.T) {
		pt, realTV := makeTable(10)
		pt.ScrollDown(5)
		row, col := realTV.GetScrollOffset()
		require.Equal(t, 15, row)
		require.Equal(t, 0, col)
	})

	t.Run("ScrollDown from zero", func(t *testing.T) {
		pt, realTV := makeTable(0)
		pt.ScrollDown(3)
		row, _ := realTV.GetScrollOffset()
		require.Equal(t, 3, row)
	})

	t.Run("ScrollUp decreases row offset", func(t *testing.T) {
		pt, realTV := makeTable(10)
		pt.ScrollUp(3)
		row, col := realTV.GetScrollOffset()
		require.Equal(t, 7, row)
		require.Equal(t, 0, col)
	})

	t.Run("ScrollUp clamps to zero when delta exceeds offset", func(t *testing.T) {
		pt, realTV := makeTable(3)
		pt.ScrollUp(10)
		row, _ := realTV.GetScrollOffset()
		require.Equal(t, 0, row)
	})

	t.Run("ScrollUp from zero stays at zero", func(t *testing.T) {
		pt, realTV := makeTable(0)
		pt.ScrollUp(5)
		row, _ := realTV.GetScrollOffset()
		require.Equal(t, 0, row)
	})

	t.Run("ScrollUp exact delta lands on zero", func(t *testing.T) {
		pt, realTV := makeTable(5)
		pt.ScrollUp(5)
		row, _ := realTV.GetScrollOffset()
		require.Equal(t, 0, row)
	})
}
