package component_test

import (
	"testing"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/dkyanakiev/vaulty/component/componentfakes"
	"github.com/dkyanakiev/vaulty/models"
	"github.com/dkyanakiev/vaulty/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
)

func TestMountsTable_Pass(t *testing.T) {
	r := require.New(t)
	t.Run("When the component is bound", func(t *testing.T) {

		fakeTable := &componentfakes.FakeTable{}
		mTable := component.NewMountsTable()

		mTable.Table = fakeTable
		mTable.Props.Namespace = "default"
		mTable.Props.Data = map[string]*models.MountOutput{
			"path-one/": {
				Type:           "kv",
				Description:    "description",
				RunningVersion: "v0.15",
			},
		}

		mTable.Props.SelectPath = func(id string) {}
		mTable.Props.HandleNoResources = func(format string, args ...interface{}) {}
		slot := tview.NewFlex()
		mTable.Bind(slot)
		// It doesn't error
		err := mTable.Render()
		r.NoError(err)

		// Render header rows
		renderHeaderCount := fakeTable.RenderHeaderCallCount()
		r.Equal(renderHeaderCount, 1)

		// Correct headers
		header := fakeTable.RenderHeaderArgsForCall(0)
		r.Equal(component.TableHeaderJobs, header)

		// It renders the correct number of rows
		renderRowCallCount := fakeTable.RenderRowCallCount()
		r.Equal(renderRowCallCount, 1)

		row1, index1, c1 := fakeTable.RenderRowArgsForCall(0)
		expectedRow1 := []string{"path-one/", "kv", "description", "v0.15"}

		r.Equal(expectedRow1, row1)
		r.Equal(index1, 1)
		r.Equal(c1, tcell.ColorGreenYellow)

	})

	t.Run("No data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		mTable := component.NewMountsTable()

		mTable.Table = fakeTable
		mTable.Props.Namespace = "default"
		mTable.Props.Data = map[string]*models.MountOutput{}

		var NoResourcesCalled bool
		mTable.Props.HandleNoResources = func(format string, args ...interface{}) {
			NoResourcesCalled = true

			r.Equal("%sno mounts available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯", format)
			r.Len(args, 2)
			r.Equal(args[0], styles.HighlightPrimaryTag)
			r.Equal(args[1], styles.HighlightSecondaryTag)
		}

		slot := tview.NewFlex()
		mTable.Bind(slot)

		// It doesn't error
		err := mTable.Render()
		r.NoError(err)

		// It handled the case that there are no resources
		r.True(NoResourcesCalled)

		// It didn't returned after handling no resources
		r.Equal(fakeTable.RenderHeaderCallCount(), 0)
		r.Equal(fakeTable.RenderRowCallCount(), 0)

	})

}

//TODO: Revisit after validation logic is fixed for mounts_table
// func TestMountsTable_Fail(t *testing.T) {
// 	r := require.New(t)

// 	t.Run("When SelectDeployment is not set", func(t *testing.T) {
// 		mt := component.NewMountsTable()

// 		mt.Props.HandleNoResources = func(format string, args ...interface{}) {}

// 		slot := tview.NewFlex()
// 		mt.Bind(slot)

// 		// It doesn't error
// 		err := mt.Render()
// 		r.NoError(err)

// 		// It provides the correct error message
// 		r.EqualError(err, "component properties not set")

// 		// It is the correct error
// 		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
// 	})
// }
