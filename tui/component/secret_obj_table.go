package component

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/dkyanakiev/vaulty/internal/models"
	primitive "github.com/dkyanakiev/vaulty/tui/primitives"
	"github.com/dkyanakiev/vaulty/tui/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/vault/api"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

const (
	SecretObjTableTitle = "Secret: "
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
	Table          Table
	TextView       TextView
	TextArea       TextArea
	Props          *SecretObjTableProps
	Logger         *zerolog.Logger
	ShowJson       bool
	ShowMetadata   bool
	Editable       bool
	CursorPosition int
	slot           *tview.Flex
}

type SecretObjTableProps struct {
	SelectedKey       string
	SelectedValue     string
	SelectedPath      string
	MissingSecret     bool
	JsonOnly          bool
	SelectPath        SelectSecretPathFunc
	HandleNoResources models.HandlerFunc

	Namespace      string
	Data           *api.Secret
	Metadata       *models.Metadata
	UpdatedData    map[string]interface{}
	ObscureSecrets bool
	Update         string
	ChangedFunc    func(text string)
}

func NewSecretObjTable() *SecretObjTable {
	t := primitive.NewTable()
	tv := primitive.NewTextView(1)
	tv.SetTextAlign(tview.AlignLeft)
	tv.SetBorderColor(styles.TcellColorStandard)
	ta := primitive.NewTextArea()

	jt := &SecretObjTable{
		Table:        t,
		TextView:     tv,
		TextArea:     ta,
		Props:        &SecretObjTableProps{},
		ShowJson:     false,
		ShowMetadata: false,
		Editable:     false,
		slot:         tview.NewFlex(),
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
	if !s.ShowMetadata {
		if !s.Editable {
			if s.Props.JsonOnly {
				s.slot.AddItem(s.TextView.Primitive(), 0, 1, true)
				s.renderJson()
			} else {
				if s.ShowJson {
					s.slot.AddItem(s.TextView.Primitive(), 0, 1, true)
					s.renderJson()
				} else {
					s.slot.AddItem(s.Table.Primitive(), 0, 1, true)
					s.renderRows()
				}
			}
		} else {
			if !s.Props.MissingSecret {
				s.Props.UpdatedData = s.Props.Data.Data["data"].(map[string]interface{})
			} else {
				s.Props.UpdatedData = make(map[string]interface{})
			}
			s.TextView.SetText(s.TextArea.GetText())
			s.slot.AddItem(s.TextArea.Primitive(), 0, 1, true)
			s.renderEditArea()
		}
	}
}

func (s *SecretObjTable) ToggleMetaView() {
	s.Logger.Debug().Msgf("ShowMetadata: %v", s.ShowMetadata)
	if s.ShowMetadata {
		s.slot.AddItem(s.Table.Primitive(), 0, 1, true)
		// s.slot.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		// 	AddItem(tview.NewTable().SetBorder(true).SetTitle("Metadata"), 10, 1, false).
		// 	AddItem(tview.NewTable().SetBorder(true).SetTitle("Custom MetaData"), 0, 1, false), 0, 2, false)
		s.renderMetadata()
	} else {
		s.Render()
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
	if !s.ShowMetadata {
		s.Props.MissingSecret = false
		s.Props.JsonOnly = false
		s.reset()
		s.Table.SetTitle("%s %s", SecretObjTableTitle, s.Props.SelectedPath)
		s.validationLogic()

		if s.Props.MissingSecret {
			s.Props.HandleNoResources(
				"%sno Secret Object data available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
				styles.HighlightPrimaryTag,
				styles.HighlightSecondaryTag,
			)
			return nil
		}

		s.Table.SetSelectedFunc(s.pathSelected)
		s.Table.RenderHeader(SecretObjTableHeaderJobs)

		if !s.Props.MissingSecret {
			s.ToggleView()
		}

	}
	return nil
}

func (s *SecretObjTable) renderMetadata() error {
	if s.Props.Metadata == nil {
		s.Props.HandleNoResources(
			"%sno Secret Object data available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)
		return nil
	}

	s.Logger.Debug().Msgf("slot: %v", s.slot.GetItemCount())
	// item := s.slot.GetItem(1)
	// // itemType := reflect.TypeOf(item)
	// count1 := item.(*tview.Flex).GetItem(0)
	// count2 := item.(*tview.Flex).GetItem(1)
	// itemType := reflect.TypeOf(count1)
	// itemType2 := reflect.TypeOf(count2)
	// s.Logger.Debug().Msgf("Item: %v", itemType)
	// s.Logger.Debug().Msgf("Item: %v", itemType2)

	// item := s.slot.GetItem(1)
	// table1 := item.(*tview.Flex).GetItem(1)
	// table2 := table1.(*tview.Box).
	// // // Assuming item is of type interface{}
	// table, ok := item.(*tview.Table)
	// if !ok {
	// 	panic(itemType)
	// }

	// Now you can call Table methods on table

	// s.Logger.Debug().Msgf("Item: %v", itemType)
	s.Table.SetSelectedFunc(s.pathSelected)
	s.Table.RenderHeader(SecretObjTableHeaderJobs)
	// table1.(*tview.Table).SetCell(0, 0, tview.NewTableCell("Created Time").SetTextColor(tcell.ColorYellow))
	// .RenderRow([]string{"Created Time", s.Props.Metadata.CreatedTime}, 0, tcell.ColorYellow)
	s.Table.RenderRow([]string{"Created Time", s.Props.Metadata.CreatedTime}, 0, tcell.ColorYellow)
	s.Table.RenderRow([]string{"Current Version", strconv.Itoa(s.Props.Metadata.CurrentVersion)}, 1, tcell.ColorYellow)
	s.Table.RenderRow([]string{"Oldest Version", strconv.Itoa(s.Props.Metadata.CurrentVersion)}, 2, tcell.ColorYellow)
	s.Table.RenderRow([]string{"Update Time", s.Props.Metadata.UpdatedTime}, 3, tcell.ColorYellow)
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
		var strValue string
		if value != nil {
			strValue = value.(string)
		} else {
			strValue = ""
		}

		row := []string{
			key,
			strValue,
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
	s.TextView.SetBorder(true)
	s.TextView.SetTitle(fmt.Sprintf("%s %s", SecretObjTableTitle, s.Props.SelectedPath))
	s.TextView.SetText(string(jsonData))
}

func (s *SecretObjTable) renderEditArea() {
	data := s.Props.UpdatedData
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		s.Logger.Err(err).Msgf("error: %s", err)
	}
	s.TextArea.SetTitle(fmt.Sprintf("%s %s", SecretObjTableTitle, s.Props.SelectedPath))
	s.TextArea.SetText(string(jsonData), true)
}

func (s *SecretObjTable) SaveData(text string) string {
	var data map[string]interface{}

	s.Logger.Debug().Msg("Saving data")
	s.Logger.Debug().Msg(text)

	err := json.Unmarshal([]byte(text), &data)
	if err != nil {
		s.Logger.Err(err).Msgf("Failed to validate json:: %s", err)
		s.Props.UpdatedData = nil
		return "Failed to validate json:"
	} else {
		s.Props.UpdatedData = data
	}
	return ""
}

func isJSONFlat(objmap map[string]interface{}) bool {
	for _, value := range objmap {
		_, ok := value.(map[string]interface{})
		if ok {
			return false
		}
	}

	return true
}

func (s *SecretObjTable) validateData() bool {
	if s.Props.Data != nil && s.Props.Data.Data != nil {

		data, ok := s.Props.Data.Data["data"].(map[string]interface{})
		if !ok {
			return false
		}

		if data == nil || len(data) == 0 {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

func (s *SecretObjTable) validationLogic() {
	validateResult := s.validateData()
	if validateResult {
		if s.Props.Data != nil && s.Props.Data.Data != nil {
			if len(s.Props.Data.Data) > 0 {
				data, ok := s.Props.Data.Data["data"]
				if ok && data != nil {
					s.Props.MissingSecret = false
					val := isJSONFlat(data.(map[string]interface{}))
					if val {
						s.Props.JsonOnly = false
					} else {
						s.Logger.Info().Msgf("Secret data is not flat json, disabing table view: %v", s.Props.JsonOnly)
						s.Props.JsonOnly = true
					}
				} else {
					s.Props.MissingSecret = true
				}
			} else {
				s.Logger.Debug().Msg("Secret data is an empty map")
				s.Props.MissingSecret = true
			}
		} else {
			s.Logger.Debug().Msg("Secret data is nil")
			s.Props.MissingSecret = true
		}
	} else {
		s.Logger.Debug().Msg("Secret data is not valid")
		s.Props.MissingSecret = true
	}
}
