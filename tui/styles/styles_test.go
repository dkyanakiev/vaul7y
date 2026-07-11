package styles

import (
	"testing"

	"github.com/dkyanakiev/vaul7y/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"
)

func TestApplyTheme(t *testing.T) {
	ApplyTheme(config.ThemeConfig{
		Background:         "#202428",
		HighlightPrimary:   "#111111",
		HighlightSecondary: "#222222",
		Standard:           "#333333",
		Active:             "#444444",
		White:              "#555555",
		LightGrey:          "#666666",
		ModalInfo:          "#777777",
		Attention:          "#888888",
	})

	require.True(t, themeConfigured)
	require.Equal(t, "#111111", HighlightPrimaryHex)
	require.Equal(t, "#222222", HighlightSecondaryHex)
	require.Equal(t, "#333333", StandardColorHex)
	require.Equal(t, "#444444", ColorActiveHex)
	require.Equal(t, "#555555", ColorWhiteHex)
	require.Equal(t, "#666666", ColorLightGreyHex)
	require.Equal(t, "#777777", ColorModalInfoHex)
	require.Equal(t, "#888888", ColorAttentionHex)
	require.Equal(t, tcell.GetColor("#202428"), TcellBackgroundColor)
	require.Equal(t, "[#333333]", StandardColorTag)
}

func TestApplyTheme_NormalizesHexColors(t *testing.T) {
	ApplyTheme(config.ThemeConfig{Standard: "333333"})

	require.Equal(t, "#333333", StandardColorHex)
	require.Equal(t, "[#333333]", StandardColorTag)
}
