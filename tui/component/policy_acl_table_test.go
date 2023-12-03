package component_test

import (
	"testing"

	"github.com/dkyanakiev/vaulty/tui/component"
	"github.com/dkyanakiev/vaulty/tui/component/componentfakes"
	"github.com/dkyanakiev/vaulty/tui/styles"
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
