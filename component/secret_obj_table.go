package component

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/dkyanakiev/vaulty/models"
	primitive "github.com/dkyanakiev/vaulty/primitives"
	"github.com/dkyanakiev/vaulty/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/vault/api"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

const (
	SecretObjTableTitle = "JSON Explorer"
)

var (
	SecretObjTableHeaderJobs = []string{
		"Key",
		"Value",
	}
	SecretObjTableHeaderJson = []string{
		"Json",
	}
)

type SelectSecretPathFunc func(jsonPath string)

type SecretObjTable struct {
	Table    Table
	TextView TextView
	Form     Form

	Props    *SecretObjTableProps
	Logger   *zerolog.Logger
	slot     *tview.Flex
	ShowJson bool
	Editable bool
}

type SecretObjTableProps struct {
	SelectedKey       string
	SelectedValue     string
	SelectedPath      string
	SelectPath        SelectSecretPathFunc
	HandleNoResources models.HandlerFunc

	Namespace      string
	Data           *api.Secret
	UpdatedData    map[string]interface{}
	ObscureSecrets bool
	ChangedFunc    func(text string)
}

func NewSecretObjTable() *SecretObjTable {
	t := primitive.NewTable()
	tv := primitive.NewTextView(1)
	tv.SetTextAlign(tview.AlignLeft)
	tv.SetBorder(true)
	tv.SetDynamicColors(true)
	tv.SetRegions(true)
	tv.SetBorderPadding(0, 0, 1, 1)
	tv.SetBorderColor(styles.TcellColorStandard)

	jt := &SecretObjTable{
		Table:    t,
		TextView: tv,
		Props:    &SecretObjTableProps{},
		ShowJson: false,
		Editable: false,
		slot:     tview.NewFlex(),
	}
	//TODO: Revisit
	jt.slot.AddItem(jt.TextView.Primitive(), 0, 1, false)
	return jt
}

func (s *SecretObjTable) Bind(slot *tview.Flex) {
	s.slot = slot
}

func (s *SecretObjTable) reset() {
	s.slot.Clear()
	s.Table.Clear()
	s.TextView.Clear()
}

func (s *SecretObjTable) ToggleView() {
	s.slot.Clear()
	if !s.Editable {
		if s.ShowJson {
			s.slot.AddItem(s.TextView.Primitive(), 0, 1, true)
			s.renderJson()
		} else {
			s.slot.AddItem(s.Table.Primitive(), 0, 1, true)
			s.renderRows()
		}
	} else {
		s.Props.UpdatedData = s.Props.Data.Data["data"].(map[string]interface{})
		s.slot.AddItem(s.TextView.Primitive(), 0, 1, true)
		s.renderEditField()
	}
}

func (s *SecretObjTable) pathSelected(row, _ int) {
	jsonPath := s.Table.GetCellContent(row, 0)
	s.Props.SelectedKey = fmt.Sprintf("%s%s", s.Props.SelectedKey, jsonPath)
}

func (s *SecretObjTable) GetIDForSelection() (string, string) {
	row, _ := s.Table.GetSelection()
	key := s.Table.GetCellContent(row, 0)
	value := s.Table.GetCellContent(row, 1)
	return key, value
}

func (s *SecretObjTable) Render() error {

	s.reset()

	s.Table.SetTitle("%s (%s)", SecretObjTableTitle, s.Props.SelectedPath)

	if s.Props.Data == nil {
		s.Props.HandleNoResources(
			"%sno Secret Object data available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)

		return nil
	}

	s.Table.SetSelectedFunc(s.pathSelected)
	s.Table.RenderHeader(SecretObjTableHeaderJobs)

	s.ToggleView()
	// s.renderRows()
	// s.slot.AddItem(s.Table.Primitive(), 0, 1, false)
	return nil
}

func (s *SecretObjTable) renderRows() {

	keys := make([]string, 0, len(s.Props.Data.Data["data"].(map[string]interface{})))
	for key := range s.Props.Data.Data["data"].(map[string]interface{}) {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i, key := range keys {
		value := s.Props.Data.Data["data"].(map[string]interface{})[key]

		if s.Props.ObscureSecrets {
			value = "********"
		}

		row := []string{
			key,
			value.(string),
		}
		index := i + 1
		c := tcell.ColorYellow

		s.Table.RenderRow(row, index, c)
	}
}

func (s *SecretObjTable) renderJson() {
	data := s.Props.Data.Data["data"].(map[string]interface{})
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		s.Logger.Err(err).Msgf("error: %s", err)
	}
	s.TextView.SetText(string(jsonData))
}

func (s *SecretObjTable) renderEditField() {
	data := s.Props.UpdatedData
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		s.Logger.Err(err).Msgf("error: %s", err)
	}
	s.TextView.SetText(string(jsonData))

}

func (s *SecretObjTable) SaveData(text string) string {
	var data map[string]interface{}

	s.Logger.Debug().Msg("Saving data")
	s.Logger.Debug().Msg(text)

	err := json.Unmarshal([]byte(text), &data)
	if err != nil {
		s.Logger.Err(err).Msgf("Failed to validate json:: %s", err)
		return "Failed to validate json:"
	}
	s.Props.UpdatedData = data
	return ""
}
