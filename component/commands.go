package component

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	primitive "github.com/dkyanakiev/vaulty/primitives"
	"github.com/dkyanakiev/vaulty/styles"
)

var (
	MainCommands = []string{
		fmt.Sprintf("%sMain Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<ctrl-m>%s to display System Mounts", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-p>%s to display ACL Policies", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-c>%s to Quit", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	MountsCommands = []string{
		fmt.Sprintf("\n%s Secret Mounts Command List:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<e>%s to explore mount", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	NoViewCommands = []string{}
	PolicyCommands = []string{
		fmt.Sprintf("\n%s ACL Policy Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<i>%s to inspect policy", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s</>%s apply filter", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	PolicyACLCommands = []string{
		fmt.Sprintf("\n%s ACL Policy Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<Esc>%s to go back", styles.HighlightPrimaryTag, styles.StandardColorTag),
		//fmt.Sprintf("%s</>%s apply filter", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	SecretsCommands = []string{
		fmt.Sprintf("\n%s Secrets Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<e>%s to navigate to selected the path", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<b>%s to go back to the previous path", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	SecretObjectCommands = []string{
		fmt.Sprintf("\n%s Secret Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<h>%s toggle display for secrets", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<c>%s copy secret to clipboard", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<j>%s toggle json view for secret", styles.HighlightPrimaryTag, styles.StandardColorTag),
		//TODO: Work in progress
		//fmt.Sprintf("%s<p>%s patch secret", styles.HighlightPrimaryTag, styles.StandardColorTag),
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