package component_test

import (
	"errors"
	"testing"

	"github.com/dkyanakiev/vaulty/component"
	"github.com/dkyanakiev/vaulty/component/componentfakes"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommands_Update(t *testing.T) {
	commands := component.NewCommands()

	newCommands := []string{"new command 1", "new command 2"}
	commands.Update(newCommands)

	assert.Equal(t, newCommands, commands.Props.ViewCommands)
}

func TestCommands_OK(t *testing.T) {
	r := require.New(t)
	textView := &componentfakes.FakeTextView{}
	cmds := component.NewCommands()
	cmds.TextView = textView

	cmds.Props.MainCommands = []string{"command1", "command2"}
	cmds.Props.ViewCommands = []string{"subCmd1", "subCmd2"}

	cmds.Bind(tview.NewFlex())
	err := cmds.Render()
	r.NoError(err)

	text := textView.SetTextArgsForCall(0)
	r.Equal(text, "command1\ncommand2\nsubCmd1\nsubCmd2")
}

func TestCommands_Fail(t *testing.T) {
	r := require.New(t)
	textView := &componentfakes.FakeTextView{}
	cmds := component.NewCommands()
	cmds.TextView = textView
	err := cmds.Render()
	r.Error(err)

	r.True(errors.Is(err, component.ErrComponentNotBound))
	r.EqualError(err, "component not bound")
}
