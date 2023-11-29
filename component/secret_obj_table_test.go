package component_test

import (
	"encoding/json"
	"testing"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/dkyanakiev/vaulty/component/componentfakes"
	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/vault/api"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
)

func TestSecretObjTable_Pass(t *testing.T) {
	r := require.New(t)

	t.Run("Render data as table", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		fakeTextView := &componentfakes.FakeTextView{}
		st := component.NewSecretObjTable()

		st.Table = fakeTable
		st.TextView = fakeTextView
		st.ShowJson = false
		st.Editable = false

		mockSecret := &api.Secret{
			RequestID:     "mockRequestID",
			LeaseID:       "mockLeaseID",
			LeaseDuration: 3600,
			Renewable:     true,
			Data: map[string]interface{}{
				"data": map[string]interface{}{
					"key1": "dZpT6XnlnktMXaYF",
					"key2": "10mNsYOLfd1OfohW",
				},
			},
		}

		st.Props.Data = mockSecret

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
		headers := fakeTable.RenderHeaderArgsForCall(0)
		r.Equal(headers, component.SecretObjTableHeaderJobs)

		// Render rows
		renderRowCallCount := fakeTable.RenderRowCallCount()
		r.Equal(renderRowCallCount, 2)

		row1, index1, c1 := fakeTable.RenderRowArgsForCall(0)
		row2, index2, c2 := fakeTable.RenderRowArgsForCall(1)

		expectedRow1 := []string{"key1", "dZpT6XnlnktMXaYF"}
		expectedRow2 := []string{"key2", "10mNsYOLfd1OfohW"}

		r.Equal(expectedRow1, row1)
		r.Equal(expectedRow2, row2)
		r.Equal(index1, 1)
		r.Equal(index2, 2)
		r.Equal(c1, tcell.ColorYellow)
		r.Equal(c2, tcell.ColorYellow)

	})

	t.Run("Render data as json", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		fakeTextView := &componentfakes.FakeTextView{}
		st := component.NewSecretObjTable()

		st.Table = fakeTable
		st.TextView = fakeTextView
		st.ShowJson = true
		st.Editable = false

		mockSecret := &api.Secret{
			RequestID:     "mockRequestID",
			LeaseID:       "mockLeaseID",
			LeaseDuration: 3600,
			Renewable:     true,
			Data: map[string]interface{}{
				"data": map[string]interface{}{
					"key1": "dZpT6XnlnktMXaYF",
					"key2": "10mNsYOLfd1OfohW",
				},
			},
		}

		st.Props.Data = mockSecret
		correctText, _ := json.Marshal(mockSecret.Data["data"])

		st.Props.SelectPath = func(id string) {}
		st.Props.HandleNoResources = func(format string, args ...interface{}) {}
		slot := tview.NewFlex()
		st.Bind(slot)
		// It doesn't error
		err := st.Render()
		r.NoError(err)

		// Renders correct text

		fakeTextView.GetTextReturns(string(correctText))
		renderedText := fakeTextView.GetText(true)

		r.Equal(string(correctText), renderedText)
	})

	t.Run("No data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		fakeTextView := &componentfakes.FakeTextView{}
		st := component.NewSecretObjTable()

		st.Table = fakeTable
		st.TextView = fakeTextView
		st.ShowJson = false
		st.Editable = false

		st.Props.Data = nil

		var NoResourcesCalled bool
		st.Props.HandleNoResources = func(format string, args ...interface{}) {
			NoResourcesCalled = true
		}

		slot := tview.NewFlex()
		st.Bind(slot)
		// It doesn't error
		err := st.Render()
		r.NoError(err)

		r.True(NoResourcesCalled)
	})
}
