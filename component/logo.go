package component

import (
	"fmt"
	"strings"

	primitive "github.com/dkyanakiev/vaulty/primitives"
	"github.com/rivo/tview"
)

var LogoASCII = []string{
	`[#00b57c]                           `,
	`____   ____            .___________        `,
	`\   \ /   /____   __ __|  \______  \___.__.`,
	` \   Y   /\__  \ |  |  \  |   /    <   |  |`,
	`  \     /  / __ \|  |  /  |__/    / \___  |`,
	`   \___/  (____  /____/|____/____/  / ____|`,
	`			    \/                   \/     `,
	`[#26ffe6]Vaul7y - Terminal Dashboard`,
}

type Logo struct {
	TextView TextView
	slot     *tview.Flex
	Props    *LogoProps
}

type LogoProps struct {
	Version string
}

func NewLogo(version string) *Logo {
	t := primitive.NewTextView(tview.AlignRight)
	return &Logo{
		TextView: t,
		Props: &LogoProps{
			Version: version,
		},
	}
}

func (l *Logo) Render() error {
	if l.slot == nil {
		return ErrComponentNotBound
	}

	versionText := fmt.Sprintf("[#26ffe6]version: %s", l.Props.Version)
	logo := strings.Join(LogoASCII, "\n")
	logo = fmt.Sprintf("%s\n%s", logo, versionText)
	l.TextView.SetText(logo)
	l.slot.AddItem(l.TextView.Primitive(), 0, 1, false)
	return nil
}

func (l *Logo) Bind(slot *tview.Flex) {
	l.slot = slot
}
