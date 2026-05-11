package component

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/dkyanakiev/vaul7y/internal/models"
	primitive "github.com/dkyanakiev/vaul7y/tui/primitives"
	"github.com/dkyanakiev/vaul7y/tui/styles"
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
	SecretObjTableHeaderVersions = []string{"Version", "Status", "Created"}
)

type SelectSecretPathFunc func(jsonPath string)

type SecretObjTable struct {
	Table               Table
	MetadataTable       Table
	CustomMetadataTable Table
	VersionHistoryTable Table
	TextView            TextView
	TextArea            TextArea
	Props               *SecretObjTableProps
	Logger              *zerolog.Logger
	ShowJson            bool
	ShowMetadata        bool
	Editable            bool
	editMu              sync.RWMutex
	CursorPosition      int
	slot                *tview.Flex
}

func (s *SecretObjTable) IsEditable() bool {
	s.editMu.RLock()
	defer s.editMu.RUnlock()
	return s.Editable
}

func (s *SecretObjTable) SetEditable(val bool) {
	s.editMu.Lock()
	defer s.editMu.Unlock()
	s.Editable = val
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
	mt := primitive.NewTable()
	cmtt := primitive.NewTable()
	vht := primitive.NewTable()
	tv := primitive.NewTextView(1)
	tv.SetTextAlign(tview.AlignLeft)
	tv.SetBorderColor(styles.TcellColorStandard)
	ta := primitive.NewTextArea()

	mt.SetSelectable(false, false)
	cmtt.SetSelectable(false, false)
	vht.SetSelectable(false, false)
	nop := zerolog.Nop()
	jt := &SecretObjTable{
		Table:               t,
		MetadataTable:       mt,
		CustomMetadataTable: cmtt,
		VersionHistoryTable: vht,
		TextView:            tv,
		TextArea:            ta,
		Props:               &SecretObjTableProps{},
		Logger:              &nop,
		ShowJson:            false,
		ShowMetadata:        false,
		Editable:            false,
		slot:                tview.NewFlex(),
	}
	jt.slot.AddItem(jt.TextView.Primitive(), 0, 1, false)
	return jt
}

func (s *SecretObjTable) Bind(slot *tview.Flex) {
	s.slot = slot
}

func (s *SecretObjTable) reset() {
	s.slot.Clear()
	s.Table.Clear()
	s.MetadataTable.Clear()
	s.CustomMetadataTable.Clear()
	s.VersionHistoryTable.Clear()
	s.TextView.Clear()
}

// clearDataTables clears only data rows without touching the slot layout.
// Used to refresh a split view in-place.
func (s *SecretObjTable) clearDataTables() {
	s.Table.Clear()
	s.MetadataTable.Clear()
	s.CustomMetadataTable.Clear()
	s.VersionHistoryTable.Clear()
	s.TextView.Clear()
}

func (s *SecretObjTable) ToggleView() {
	s.slot.Clear()
	if !s.ShowMetadata {
		if !s.IsEditable() {
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
		s.slot.Clear()

		// Left panel: normal secret content (table or JSON).
		var mainContent tview.Primitive
		if s.Props.JsonOnly || s.ShowJson {
			s.renderJson()
			mainContent = s.TextView.Primitive()
		} else {
			s.Table.SetTitle("%s %s", SecretObjTableTitle, s.Props.SelectedPath)
			s.Table.RenderHeader(SecretObjTableHeaderJobs)
			s.renderRows()
			mainContent = s.Table.Primitive()
		}

		// Right panel: metadata / custom-metadata / version-history stacked vertically.
		metaPanel := tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(s.MetadataTable.Primitive(), 10, 0, false).
			AddItem(s.CustomMetadataTable.Primitive(), 5, 0, false).
			AddItem(s.VersionHistoryTable.Primitive(), 0, 1, false)

		s.slot.SetDirection(tview.FlexColumn).
			AddItem(mainContent, 0, 2, true).
			AddItem(metaPanel, 0, 1, false)

		s.renderMetadata()
	} else {
		s.slot.SetDirection(tview.FlexRow)
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
	// When the split view is active, refresh data in-place without rebuilding
	// the slot layout (which would destroy the split).
	if s.ShowMetadata {
		s.clearDataTables()
		s.validationLogic()
		if !s.Props.MissingSecret {
			if s.Props.JsonOnly || s.ShowJson {
				s.renderJson()
			} else {
				s.Table.SetTitle("%s %s", SecretObjTableTitle, s.Props.SelectedPath)
				s.Table.RenderHeader(SecretObjTableHeaderJobs)
				s.renderRows()
			}
		}
		s.renderMetadata()
		return nil
	}

	s.Props.MissingSecret = false
	s.Props.JsonOnly = false
	s.reset()
	s.Table.SetTitle("%s %s", SecretObjTableTitle, s.Props.SelectedPath)
	s.validationLogic()

	if s.Props.MissingSecret {
		s.Props.HandleNoResources(
			"%sno Secret Object data available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
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

	return nil
}

func (s *SecretObjTable) renderMetadata() error {
	if s.Props.Metadata == nil {
		return nil
	}

	deleteAfter := s.Props.Metadata.DeleteVersionAfter
	if deleteAfter == "" || deleteAfter == "0s" {
		deleteAfter = "Never"
	}

	s.MetadataTable.SetTitle("Secret Metadata")
	s.MetadataTable.RenderRow([]string{"Updated Time", ConvertTimeFormat(s.Props.Metadata.UpdatedTime)}, 0, tcell.ColorYellow)
	s.MetadataTable.RenderRow([]string{"Created Time", ConvertTimeFormat(s.Props.Metadata.CreatedTime)}, 1, tcell.ColorYellow)
	s.MetadataTable.RenderRow([]string{"Current Version", strconv.Itoa(s.Props.Metadata.CurrentVersion)}, 2, tcell.ColorYellow)
	s.MetadataTable.RenderRow([]string{"Oldest Version", strconv.Itoa(s.Props.Metadata.OldestVersion)}, 3, tcell.ColorYellow)
	s.MetadataTable.RenderRow([]string{"Max Versions", strconv.Itoa(s.Props.Metadata.MaxVersions)}, 4, tcell.ColorYellow)
	s.MetadataTable.RenderRow([]string{"CAS Required", strconv.FormatBool(s.Props.Metadata.CasRequired)}, 5, tcell.ColorYellow)
	s.MetadataTable.RenderRow([]string{"Delete Version After", deleteAfter}, 6, tcell.ColorYellow)

	s.CustomMetadataTable.SetTitle("Custom Metadata")
	i := 0
	for k, v := range s.Props.Metadata.CustomMetadata {
		value, ok := v.(string)
		if !ok {
			continue
		}
		s.CustomMetadataTable.RenderRow([]string{k, value}, i, tcell.ColorYellow)
		i++
	}

	s.renderVersionHistory()
	return nil
}

func (s *SecretObjTable) renderVersionHistory() {
	if s.Props.Metadata == nil || len(s.Props.Metadata.Versions) == 0 {
		return
	}

	s.VersionHistoryTable.SetTitle("Version History")
	s.VersionHistoryTable.RenderHeader(SecretObjTableHeaderVersions)

	type vEntry struct {
		num int
		v   models.Version
	}
	entries := make([]vEntry, 0, len(s.Props.Metadata.Versions))
	for numStr, v := range s.Props.Metadata.Versions {
		n, _ := strconv.Atoi(numStr)
		entries = append(entries, vEntry{n, v})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].num > entries[j].num })

	for i, e := range entries {
		status := "active"
		color := tcell.ColorWhite
		if e.num == s.Props.Metadata.CurrentVersion {
			status = "current"
			color = tcell.ColorGreen
		} else if e.v.Destroyed {
			status = "destroyed"
			color = tcell.ColorRed
		} else if e.v.DeletionTime != "" {
			status = "deleted"
			color = tcell.ColorOrange
		}
		s.VersionHistoryTable.RenderRow([]string{
			fmt.Sprintf("v%d", e.num),
			status,
			ConvertTimeFormat(e.v.CreatedTime),
		}, i+1, color)
	}
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
			strValue = fmt.Sprintf("%v", value)
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
	s.SetEditStatus("")
	s.TextArea.SetText(string(jsonData), false)

	s.TextArea.SetChangedFunc(func() {
		var d map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(s.TextArea.GetText()), &d); jsonErr != nil {
			s.TextArea.SetBorderColor(tcell.ColorRed)
		} else {
			s.TextArea.SetBorderColor(tcell.ColorGreen)
			s.SetEditStatus("")
		}
	})
}

// SetEditStatus updates the TextArea title and border to reflect a validation
// message. An empty msg resets both to their default state.
func (s *SecretObjTable) SetEditStatus(msg string) {
	if msg != "" {
		s.TextArea.SetBorderColor(tcell.ColorRed)
		s.TextArea.SetTitle(msg)
	} else {
		s.TextArea.SetBorderColor(styles.TcellColorStandard)
		s.TextArea.SetTitle(fmt.Sprintf("%s %s", SecretObjTableTitle, s.Props.SelectedPath))
	}
}

func (s *SecretObjTable) SaveData(text string) string {
	var data map[string]interface{}

	s.Logger.Debug().Msg("Saving data")
	s.Logger.Debug().Msg(text)

	err := json.Unmarshal([]byte(text), &data)
	if err != nil {
		s.Logger.Err(err).Msgf("Failed to validate json: %s", err)
		s.Props.UpdatedData = nil
		return fmt.Sprintf("Invalid JSON: %s", err.Error())
	}
	s.Props.UpdatedData = data
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

func ConvertTimeFormat(input string) string {
	t, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		return input
	}
	return t.Format("Monday, 02-Jan-06 15:04:05 MST")
}
