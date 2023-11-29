package component_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/dkyanakiev/vaulty/component/componentfakes"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
)

func TestLogo_Pass(t *testing.T) {
	r := require.New(t)

	textView := &componentfakes.FakeTextView{}
	logo := component.NewLogo()
	logo.TextView = textView

	logo.Bind(tview.NewFlex())

	err := logo.Render()
	r.NoError(err)

	text := textView.SetTextArgsForCall(0)
	expectedLogo := strings.Join(component.LogoASCII, "\n")
	r.Equal(text, expectedLogo)
}

func TestLogo_Fail(t *testing.T) {
	r := require.New(t)
	logo := component.NewLogo()

	err := logo.Render()
	r.Error(err)

	r.True(errors.Is(err, component.ErrComponentNotBound))
	r.EqualError(err, "component not bound")
}
