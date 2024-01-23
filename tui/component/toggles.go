package component

import (
	"fmt"
	"strconv"

	primitive "github.com/dkyanakiev/vaulty/tui/primitives"
	"github.com/dkyanakiev/vaulty/tui/styles"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

type TogglesInfo struct {
	TextView TextView
	Props    *TogglesInfoProps
	Logger   *zerolog.Logger

	slot *tview.Flex
}

type TogglesInfoProps struct {
	Info       string
	Namespace  string
	FilterText string ""
	Editable   bool
}

func NewTogglesInfo() *TogglesInfo {
	return &TogglesInfo{
		TextView: primitive.NewTextView(tview.AlignLeft),
		Props:    &TogglesInfoProps{},
	}
}

func (t *TogglesInfo) InitialRender(ns string) error {
	if t.slot == nil {
		return ErrComponentNotBound
	}
	editableStr := strconv.FormatBool(t.Props.Editable)

	text := fmt.Sprintf(
		"\n%sNamespace: %s %s \n%sEdit Mode: %s %s  \n%sFilter: %s %s \n",
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		ns,
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		editableStr,
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		t.Props.FilterText,
	)

	t.TextView.SetText(text)
	t.slot.AddItem(t.TextView.Primitive(), 0, 1, false)

	return nil
}

func (t *TogglesInfo) Render() error {
	if t.slot == nil {
		return ErrComponentNotBound
	}
	editableStr := strconv.FormatBool(t.Props.Editable)

	text := fmt.Sprintf(
		"\n%sNamespace: %s %s \n%sEdit Mode: %s %s  \n%sFilter: %s %s \n",
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		t.Props.Namespace,
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		editableStr,
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		t.Props.FilterText,
	)

	t.TextView.SetText(text)
	//t.slot.AddItem(t.TextView.Primitive(), 0, 1, false)

	return nil
}

func (t *TogglesInfo) Bind(slot *tview.Flex) {
	t.slot = slot
}
