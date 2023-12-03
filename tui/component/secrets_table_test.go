package component_test

import (
	"testing"

	"github.com/dkyanakiev/vaulty/internal/models"
	"github.com/dkyanakiev/vaulty/tui/component"
	"github.com/dkyanakiev/vaulty/tui/component/componentfakes"
	"github.com/dkyanakiev/vaulty/tui/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
)

func TestSecretsTable_Pass(t *testing.T) {
	r := require.New(t)
	t.Run("When there is data to render", func(t *testing.T) {

		fakeTable := &componentfakes.FakeTable{}
		st := component.NewSecretsTable()

		st.Table = fakeTable
		st.Props.Namespace = "default"
		mockData := []models.SecretPath{
			{
				PathName: "mockPathName1",
				IsSecret: true,
			},
			{
				PathName: "mockPathName2",
				IsSecret: false,
			},
		}

		st.Props.Data = mockData
		st.Props.SelectPath = func(id string) {}
		st.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		st.Bind(slot)
		// It doesn't error
		err := st.Render()
		r.NoError(err)

		// Render header rows
		renderHeaderCount := fakeTable.RenderHeaderCallCount()
		r.Equal(renderHeaderCount, 1)

		// Correct headers
		header := fakeTable.RenderHeaderArgsForCall(0)
		r.Equal(component.SecretsTableHeaderJobs, header)

		// It renders the correct number of rows
		renderRowCallCount := fakeTable.RenderRowCallCount()
		r.Equal(renderRowCallCount, 2)

		row1, index1, c1 := fakeTable.RenderRowArgsForCall(0)
		row2, index2, c2 := fakeTable.RenderRowArgsForCall(1)
		expectedRow1 := []string{"mockPathName1", "true"}
		expectedRow2 := []string{"mockPathName2", "false"}

		r.Equal(expectedRow1, row1)
		r.Equal(expectedRow2, row2)

		r.Equal(index1, 1)
		r.Equal(index2, 2)
		r.Equal(c1, tcell.ColorYellow)
		r.Equal(c2, tcell.ColorYellow)

	})

	t.Run("No data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		st := component.NewSecretsTable()

		st.Table = fakeTable
		st.Props.Namespace = "default"

		st.Props.Data = nil
		st.Props.SelectPath = func(id string) {}
		st.Props.HandleNoResources = func(format string, args ...interface{}) {}

		var NoResourcesCalled bool
		st.Props.HandleNoResources = func(format string, args ...interface{}) {
			NoResourcesCalled = true

			r.Equal("%sno secrets available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯", format)
			r.Len(args, 2)
			r.Equal(args[0], styles.HighlightPrimaryTag)
			r.Equal(args[1], styles.HighlightSecondaryTag)
		}
		slot := tview.NewFlex()
		st.Bind(slot)
		// It doesn't error
		err := st.Render()
		r.NoError(err)
		r.True(NoResourcesCalled)

	})

}
