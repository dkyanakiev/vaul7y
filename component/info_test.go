package component_test

import (
	"errors"
	"testing"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/dkyanakiev/vaulty/component/componentfakes"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
)

func TestInfo_Pass(t *testing.T) {
	r := require.New(t)

	e := component.NewInfo()

	modal := &componentfakes.FakeModal{}
	e.Modal = modal
	pages := tview.NewPages()
	e.Bind(pages)

	e.Props.Done = func(buttonIndex int, buttonLabel string) {

	}

	err := e.Render("Info")
	r.NoError(err)

	// actualDone := modal.SetDoneFuncArgsForCall(0)
	text := modal.SetTextArgsForCall(0)

	// actualDone(0, "OK")

	r.Equal(text, "Info")
}

func TestInfo_Fail(t *testing.T) {
	r := require.New(t)

	t.Run("When the component isn't bound", func(t *testing.T) {
		e := component.NewInfo()

		e.Props.Done = func(buttonIndex int, buttonLabel string) {}

		err := e.Render("Info")
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component not bound")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentNotBound))
	})

	t.Run("When DoneFunc is not set", func(t *testing.T) {
		e := component.NewInfo()

		pages := tview.NewPages()
		e.Bind(pages)

		err := e.Render("error")
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})
}
