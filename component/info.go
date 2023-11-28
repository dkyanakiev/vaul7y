package component

import (
	"github.com/rivo/tview"

	primitive "github.com/dkyanakiev/vaulty/primitives"
	"github.com/dkyanakiev/vaulty/styles"
)

const PageNameInfo = "info"

type Info struct {
	Modal Modal
	Props *InfoProps
	pages *tview.Pages
}

type InfoProps struct {
	Done DoneModalFunc
}

func NewInfo() *Info {
	buttons := []string{"OK"}
	modal := primitive.NewModal("Info", buttons, styles.TcellColorModalInfo)

	return &Info{
		Modal: modal,
		Props: &InfoProps{},
	}
}

func (i *Info) Render(msg string) error {
	if i.Props.Done == nil {
		return ErrComponentPropsNotSet
	}

	if i.pages == nil {
		return ErrComponentNotBound
	}

	i.Props.Done = func(buttonIndex int, buttonLabel string) {
		i.pages.RemovePage(PageNameInfo)

	}
	i.Modal.SetDoneFunc(i.Props.Done)
	i.Modal.SetText(msg)
	i.pages.AddPage(PageNameInfo, i.Modal.Container(), true, true)

	return nil
}

func (i *Info) Bind(pages *tview.Pages) {
	i.pages = pages
}
