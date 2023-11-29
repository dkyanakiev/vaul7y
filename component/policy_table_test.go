package component_test

import (
	"testing"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/dkyanakiev/vaulty/component/componentfakes"
	"github.com/dkyanakiev/vaulty/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
)

func TestPolicyTable_Pass(t *testing.T) {
	r := require.New(t)
	t.Run("When the component is bound", func(t *testing.T) {

		fakeTable := &componentfakes.FakeTable{}
		pt := component.NewPolicyTable()

		pt.Table = fakeTable
		pt.Props.Data = []string{"policy1", "policy2", "policy3"}

		pt.Props.SelectPath = func(id string) {}
		pt.Props.HandleNoResources = func(format string, args ...interface{}) {}
		slot := tview.NewFlex()
		pt.Bind(slot)
		// It doesn't error
		err := pt.Render()
		r.NoError(err)

		// Render header rows
		renderHeaderCount := fakeTable.RenderHeaderCallCount()
		r.Equal(renderHeaderCount, 1)

		// Correct headers
		header := fakeTable.RenderHeaderArgsForCall(0)
		r.Equal(component.PolicyTableHeaderJobs, header)

		// It renders the correct number of rows
		renderRowCallCount := fakeTable.RenderRowCallCount()
		r.Equal(renderRowCallCount, 3)

		row1, index1, c1 := fakeTable.RenderRowArgsForCall(0)
		row2, index2, c2 := fakeTable.RenderRowArgsForCall(1)
		row3, index3, c3 := fakeTable.RenderRowArgsForCall(2)

		expectedRow1 := []string{"policy1"}
		expectedRow2 := []string{"policy2"}
		expectedRow3 := []string{"policy3"}

		r.Equal(expectedRow1, row1)
		r.Equal(expectedRow2, row2)
		r.Equal(expectedRow3, row3)
		r.Equal(index1, 1)
		r.Equal(index2, 2)
		r.Equal(index3, 3)

		r.Equal(c1, tcell.ColorYellow)
		r.Equal(c2, tcell.ColorYellow)
		r.Equal(c3, tcell.ColorYellow)

	})

	t.Run("No data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		pt := component.NewPolicyTable()

		pt.Table = fakeTable
		pt.Props.Data = []string{}

		pt.Props.HandleNoResources = func(format string, args ...interface{}) {}
		var NoResourcesCalled bool
		pt.Props.HandleNoResources = func(format string, args ...interface{}) {
			NoResourcesCalled = true

			r.Equal("%sNo policy found\n%s\\(╯°□°)╯︵ ┻━┻", format)
			r.Len(args, 2)
			r.Equal(args[0], styles.HighlightPrimaryTag)
			r.Equal(args[1], styles.HighlightSecondaryTag)
		}

		slot := tview.NewFlex()
		pt.Bind(slot)
		// It doesn't error
		err := pt.Render()
		r.NoError(err)

		// It handled the case that there are no resources
		r.True(NoResourcesCalled)

		// It didn't returned after handling no resources
		r.Equal(fakeTable.RenderHeaderCallCount(), 0)
		r.Equal(fakeTable.RenderRowCallCount(), 0)

	})

}
