package component_test

import (
	"errors"
	"testing"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/dkyanakiev/vaulty/component/componentfakes"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
)

func TestVaultInfo_Pass(t *testing.T) {
	r := require.New(t)

	textView := &componentfakes.FakeTextView{}
	vaultInfo := component.NewVaultInfo()
	vaultInfo.TextView = textView

	vaultInfo.Props.Info = "info"

	vaultInfo.Bind(tview.NewFlex())

	err := vaultInfo.Render()
	r.NoError(err)

	text := textView.SetTextArgsForCall(0)
	r.Equal(text, "info")
}

func TestVaultInfo_Failt(t *testing.T) {
	r := require.New(t)

	textView := &componentfakes.FakeTextView{}
	vaultInfo := component.NewVaultInfo()
	vaultInfo.TextView = textView
	vaultInfo.Props.Info = "info"

	err := vaultInfo.Render()
	r.Error(err)

	r.True(errors.Is(err, component.ErrComponentNotBound))
	r.EqualError(err, "component not bound")
}
