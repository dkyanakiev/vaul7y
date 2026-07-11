package styles

import (
	"fmt"
	"strings"

	"github.com/dkyanakiev/vaul7y/internal/config"
	"github.com/gdamore/tcell/v2"
)

func GetBackgroundColor() tcell.Color {
	return TcellBackgroundColor
}

var (
	themeConfigured = false

	TcellBackgroundColor = tcell.NewRGBColor(40, 44, 48)

	HighlightPrimaryHex   = "#26ffe6"
	HighlightSecondaryHex = "#baff26"
	StandardColorHex      = "#00b57c"
	ColorActiveHex        = "#b3f1ff"
	ColorWhiteHex         = "#ffffff"
	ColorLightGreyHex     = "#cccccc"
	ColorModalInfoHex     = "#61877f"
	ColorAttentionHex     = "#d98b6a"

	StandardColorTag      = fmt.Sprintf("[%s]", StandardColorHex)
	HighlightPrimaryTag   = fmt.Sprintf("[%s]", HighlightPrimaryHex)
	HighlightSecondaryTag = fmt.Sprintf("[%s]", HighlightSecondaryHex)
	ColorActiveTag        = fmt.Sprintf("[%s]", ColorActiveHex)
	ColorWhiteTag         = fmt.Sprintf("[%s]", ColorWhiteHex)
	ColorLighGreyTag      = fmt.Sprintf("[%s]", ColorLightGreyHex)

	TcellColorHighlighPrimary   = tcell.GetColor(HighlightPrimaryHex)
	TcellColorHighlighSecondary = tcell.GetColor(HighlightSecondaryHex)
	TcellColorStandard          = tcell.GetColor(StandardColorHex)
	TcellColorActive            = tcell.GetColor(ColorActiveHex)
	TcellColorWhite             = tcell.GetColor(ColorWhiteHex)
	TcellColorLightGrey         = tcell.GetColor(ColorLightGreyHex)
	TcellColorModalInfo         = tcell.GetColor(ColorModalInfoHex)
	TcellColorAttention         = tcell.GetColor(ColorAttentionHex)
)

func ApplyTheme(theme config.ThemeConfig) {
	themeConfigured = hasTheme(theme)
	if theme.Background != "" {
		TcellBackgroundColor = tcell.GetColor(normalizeHexColor(theme.Background))
	}
	if theme.HighlightPrimary != "" {
		HighlightPrimaryHex = normalizeHexColor(theme.HighlightPrimary)
	}
	if theme.HighlightSecondary != "" {
		HighlightSecondaryHex = normalizeHexColor(theme.HighlightSecondary)
	}
	if theme.Standard != "" {
		StandardColorHex = normalizeHexColor(theme.Standard)
	}
	if theme.Active != "" {
		ColorActiveHex = normalizeHexColor(theme.Active)
	}
	if theme.White != "" {
		ColorWhiteHex = normalizeHexColor(theme.White)
	}
	if theme.LightGrey != "" {
		ColorLightGreyHex = normalizeHexColor(theme.LightGrey)
	}
	if theme.ModalInfo != "" {
		ColorModalInfoHex = normalizeHexColor(theme.ModalInfo)
	}
	if theme.Attention != "" {
		ColorAttentionHex = normalizeHexColor(theme.Attention)
	}

	StandardColorTag = fmt.Sprintf("[%s]", StandardColorHex)
	HighlightPrimaryTag = fmt.Sprintf("[%s]", HighlightPrimaryHex)
	HighlightSecondaryTag = fmt.Sprintf("[%s]", HighlightSecondaryHex)
	ColorActiveTag = fmt.Sprintf("[%s]", ColorActiveHex)
	ColorWhiteTag = fmt.Sprintf("[%s]", ColorWhiteHex)
	ColorLighGreyTag = fmt.Sprintf("[%s]", ColorLightGreyHex)

	TcellColorHighlighPrimary = tcell.GetColor(HighlightPrimaryHex)
	TcellColorHighlighSecondary = tcell.GetColor(HighlightSecondaryHex)
	TcellColorStandard = tcell.GetColor(StandardColorHex)
	TcellColorActive = tcell.GetColor(ColorActiveHex)
	TcellColorWhite = tcell.GetColor(ColorWhiteHex)
	TcellColorLightGrey = tcell.GetColor(ColorLightGreyHex)
	TcellColorModalInfo = tcell.GetColor(ColorModalInfoHex)
	TcellColorAttention = tcell.GetColor(ColorAttentionHex)
}

func TableSelectedStyle() (tcell.Style, bool) {
	if !themeConfigured {
		return tcell.StyleDefault, false
	}
	return tcell.StyleDefault.Foreground(TcellColorWhite).Background(TcellColorActive), true
}

func RowColor() tcell.Color {
	if !themeConfigured {
		return tcell.ColorYellow
	}
	return TcellColorHighlighPrimary
}

func NeutralColor() tcell.Color {
	if !themeConfigured {
		return tcell.ColorWhite
	}
	return TcellColorWhite
}

func SuccessColor() tcell.Color {
	if !themeConfigured {
		return tcell.ColorGreen
	}
	return TcellColorHighlighSecondary
}

func WarningColor() tcell.Color {
	if !themeConfigured {
		return tcell.ColorOrange
	}
	return TcellColorHighlighPrimary
}

func ErrorColor() tcell.Color {
	if !themeConfigured {
		return tcell.ColorRed
	}
	return TcellColorAttention
}

func MountCubbyholeColor() tcell.Color {
	if !themeConfigured {
		return tcell.ColorYellow
	}
	return TcellColorHighlighPrimary
}

func MountIdentityColor() tcell.Color {
	if !themeConfigured {
		return tcell.ColorRed
	}
	return TcellColorAttention
}

func MountKVColor() tcell.Color {
	if !themeConfigured {
		return tcell.ColorGreenYellow
	}
	return TcellColorHighlighSecondary
}

func MountPkiColor() tcell.Color {
	if !themeConfigured {
		return tcell.ColorBlue
	}
	return TcellColorModalInfo
}

func hasTheme(theme config.ThemeConfig) bool {
	return theme.Background != "" ||
		theme.HighlightPrimary != "" ||
		theme.HighlightSecondary != "" ||
		theme.Standard != "" ||
		theme.Active != "" ||
		theme.White != "" ||
		theme.LightGrey != "" ||
		theme.ModalInfo != "" ||
		theme.Attention != ""
}

func normalizeHexColor(color string) string {
	if strings.HasPrefix(color, "#") {
		return color
	}
	return fmt.Sprintf("#%s", color)
}
