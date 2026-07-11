package component

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	primitive "github.com/dkyanakiev/vaul7y/tui/primitives"
	"github.com/dkyanakiev/vaul7y/tui/styles"
)

var (
	MainCommands = []string{
		fmt.Sprintf("%sMain Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%sctrl-b%s to display Secret Engines", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sctrl-p%s to display ACL Policies", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sctrl-t%s to display Namespaces", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sctrl-a%s to display Auth Methods", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sctrl-c%s to Quit", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	MountsCommands = []string{
		fmt.Sprintf("\n%s Secret Engines Command List:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%se or Enter%s to explore mount", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	NoViewCommands = []string{}
	PolicyCommands = []string{
		fmt.Sprintf("\n%s ACL Policy Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%si or <Enter> %s to inspect policy", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s/%s Filter policies ", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	PolicyACLCommands = []string{
		fmt.Sprintf("\n%s ACL Policy Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%sj/k%s scroll line down/up", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sd/u%s scroll half-page down/up", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sg/G%s jump to top/bottom", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sE%s to edit policy", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sc%s copy policy to clipboard", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sw%s toggle word wrap", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sesc%s to go back", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	PolicyACLEditCommands = []string{
		fmt.Sprintf("\n%s ACL Policy Edit:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%sctrl-w%s to save policy", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sesc%s to cancel edit", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	SecretsCommands = []string{
		fmt.Sprintf("\n%s Secrets Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%se or enter%s to navigate to selected the path", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sb or esc%s to go back to the previous path", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sctrl-n%s to Create a new secret ", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sJ%s to jump to a specific path", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s/%s Filter objects ", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sctrl-f%s recursive search from current path", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	SearchResultsCommands = []string{
		fmt.Sprintf("\n%s Search Results Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%se or enter%s to open the selected secret", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sv%s toggle matching secret values (re-runs search)", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sb or esc%s to go back to the secrets browser", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	SecretObjectCommands = []string{
		fmt.Sprintf("\n%s Secret Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%sh%s toggle display for secrets", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%st%s toggle metadata panel (v2 only)", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sc%s copy secret to clipboard", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sj%s toggle json view for secret", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sP%s to PATCH secret", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sU%s to UPDATE secret", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sD%s delete current version (v2: soft, v1: permanent)", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sR%s rollback to a previous version (v2 only)", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sX%s permanently destroy a version (v2 only)", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sM%s edit secret metadata (v2 only)", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sb or esc%s to go back to the previous path", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	SecretsObjectPatchCommands = []string{
		fmt.Sprintf("\n%s Secret Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%sctrl-v%s to paste from clipboard", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sctrl-w%s to submit your PATCH/UPDATE request", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sesc%s to go back to the previous path", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	NamespaceObjectCommands = []string{
		fmt.Sprintf("\n%s Namespace Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%sctrl-d%s to back to default namespace", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%sctrl-w%s to back to root namespace", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}
	AuthCommands = []string{
		fmt.Sprintf("\n%s Auth Method Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%sesc%s to go back", styles.HighlightPrimaryTag, styles.StandardColorTag),
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
	// Easy way to handle long list of commands for views
	if len(c.Props.ViewCommands) > 6 {
		commands = c.Props.ViewCommands
	}
	cmds := strings.Join(commands, "\n")
	c.TextView.SetText(cmds)
}

func (c *Commands) Bind(slot *tview.Flex) {
	c.slot = slot
}
