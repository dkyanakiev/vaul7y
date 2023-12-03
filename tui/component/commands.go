package component

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	primitive "github.com/dkyanakiev/vaulty/tui/primitives"
	"github.com/dkyanakiev/vaulty/tui/styles"
)

var (
	MainCommands = []string{
		fmt.Sprintf("%sMain Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<ctrl-b>%s to display Secret Engines", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-p>%s to display ACL Policies", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-c>%s to Quit", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	MountsCommands = []string{
		fmt.Sprintf("\n%s Secret Engines Command List:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<e> or <Enter>%s to explore mount", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	NoViewCommands = []string{}
	PolicyCommands = []string{
		fmt.Sprintf("\n%s ACL Policy Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<i> or <Enter> %s to inspect policy", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s</>%s Filter policies ", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	PolicyACLCommands = []string{
		fmt.Sprintf("\n%s ACL Policy Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<esc>%s to go back", styles.HighlightPrimaryTag, styles.StandardColorTag),
		// fmt.Sprintf("%s</>%s apply filter", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	SecretsCommands = []string{
		fmt.Sprintf("\n%s Secrets Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<e> or <enter>%s to navigate to selected the path", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<b> or <esc>%s to go back to the previous path", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-n>%s to Create a new secret ", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s</>%s Filter objects ", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	SecretObjectCommands = []string{
		fmt.Sprintf("\n%s Secret Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<h>%s toggle display for secrets", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<c>%s copy secret to clipboard", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<j>%s toggle json view for secret", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<P>%s to PATCH secret", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<U>%s to UPDATE secret", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<b> or <esc>%s to go back to the previous path", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	SecretsObjectPatchCommands = []string{
		fmt.Sprintf("\n%s Secret Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<ctrl-w>%s to submit your PATCH/UPDATE request", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<esc>%s to go back to the previous path", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
)

type Commands struct {
	TextView TextView
	Props    *CommandsProps
	slot     *tview.Flex
}

type CommandsProps struct {
	MainCommands []string
	ViewCommands []string
}

func NewCommands() *Commands {
	return &Commands{
		TextView: primitive.NewTextView(tview.AlignLeft),
		Props: &CommandsProps{
			MainCommands: MainCommands,
			// ViewCommands: MainCommands,
		},
	}
}

func (c *Commands) Update(commands []string) {
	c.Props.ViewCommands = commands

	c.updateText()
}

func (c *Commands) Render() error {
	if c.slot == nil {
		return ErrComponentNotBound
	}

	c.updateText()

	c.slot.AddItem(c.TextView.Primitive(), 0, 1, false)
	return nil
}

func (c *Commands) updateText() {
	commands := append(c.Props.MainCommands, c.Props.ViewCommands...)
	cmds := strings.Join(commands, "\n")
	c.TextView.SetText(cmds)
}

func (c *Commands) Bind(slot *tview.Flex) {
	c.slot = slot
}
