package component_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/dkyanakiev/vaulty/tui/component"
	"github.com/dkyanakiev/vaulty/tui/component/componentfakes"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
)

func TestLogo_Pass(t *testing.T) {
	r := require.New(t)

	textView := &componentfakes.FakeTextView{}
	logo := component.NewLogo("0.0.0")
	logo.TextView = textView

	logo.Bind(tview.NewFlex())

	err := logo.Render()
	r.NoError(err)

	text := textView.SetTextArgsForCall(0)
	versionText := fmt.Sprintf("[#26ffe6]version: %s", "0.0.0")
	expectedLogo := strings.Join(component.LogoASCII, "\n")
	expectedLogo = fmt.Sprintf("%s\n%s", expectedLogo, versionText)
	r.Equal(text, expectedLogo)
}

func TestLogo_Fail(t *testing.T) {
	r := require.New(t)
	logo := component.NewLogo("0.0.0")

	err := logo.Render()
	r.Error(err)

	r.True(errors.Is(err, component.ErrComponentNotBound))
	r.EqualError(err, "component not bound")
}
