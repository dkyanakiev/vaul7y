package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dkyanakiev/vaul7y/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeConfig writes yaml to a temp file and returns its path.
func writeConfig(t *testing.T, content string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), "vaul7y.yaml")
	require.NoError(t, os.WriteFile(f, []byte(content), 0600))
	return f
}

// baseConfig is a minimal YAML that satisfies required fields so tests that
// use t.Setenv for other assertions don't accidentally hit os.Exit(1).
const baseYAML = `
vault_addr: https://base.vault.example.com
vault_token: base-token
`

func TestLoadConfig_YAML(t *testing.T) {
	t.Run("non-auth fields parsed from yaml", func(t *testing.T) {
		// Use env var for VAULT_TOKEN so the ~/.vault-token file on disk
		// (if present) cannot override what we assert on for other fields.
		t.Setenv("VAULT_ADDR", "")
		t.Setenv("VAULT_NAMESPACE", "")
		t.Setenv("VAULT_TOKEN", "env-token")

		cfg := config.LoadConfig(writeConfig(t, `
vault_addr: https://yaml.vault.example.com
vault_token: yaml-token
vault_namespace: yaml-ns
vaulty_log_file: /tmp/vaul7y.log
vaulty_log_level: info
vaulty_refresh_rate: 45
theme:
  background: "#202428"
  highlight_primary: "#111111"
  highlight_secondary: "#222222"
  standard: "#333333"
  active: "#444444"
  white: "#555555"
  light_grey: "#666666"
  modal_info: "#777777"
  attention: "#888888"
`))
		assert.Equal(t, "https://yaml.vault.example.com", cfg.VaultAddr)
		assert.Equal(t, "yaml-ns", cfg.VaultNamespace)
		assert.Equal(t, "/tmp/vaul7y.log", cfg.VaultyLogFile)
		assert.Equal(t, "info", cfg.VaultyLogLevel)
		assert.Equal(t, 45, cfg.VaultyRefreshRate)
		assert.Equal(t, "#202428", cfg.Theme.Background)
		assert.Equal(t, "#111111", cfg.Theme.HighlightPrimary)
		assert.Equal(t, "#222222", cfg.Theme.HighlightSecondary)
		assert.Equal(t, "#333333", cfg.Theme.Standard)
		assert.Equal(t, "#444444", cfg.Theme.Active)
		assert.Equal(t, "#555555", cfg.Theme.White)
		assert.Equal(t, "#666666", cfg.Theme.LightGrey)
		assert.Equal(t, "#777777", cfg.Theme.ModalInfo)
		assert.Equal(t, "#888888", cfg.Theme.Attention)
		// Token was overridden by env, so we only check it came from env.
		assert.Equal(t, "env-token", cfg.VaultToken)
	})

	t.Run("default refresh rate applied when yaml omits it", func(t *testing.T) {
		t.Setenv("VAULT_TOKEN", "env-token")
		cfg := config.LoadConfig(writeConfig(t, `
vault_addr: https://yaml.vault.example.com
vault_token: yaml-token
`))
		assert.Equal(t, 30, cfg.VaultyRefreshRate)
	})
}

func TestLoadConfig_EnvVars(t *testing.T) {
	base := writeConfig(t, baseYAML)

	t.Run("VAULT_ADDR env var overrides yaml", func(t *testing.T) {
		t.Setenv("VAULT_ADDR", "https://env.vault.example.com")
		t.Setenv("VAULT_TOKEN", "env-token")

		cfg := config.LoadConfig(base)

		assert.Equal(t, "https://env.vault.example.com", cfg.VaultAddr)
	})

	t.Run("VAULT_NAMESPACE env var loaded", func(t *testing.T) {
		t.Setenv("VAULT_TOKEN", "env-token")
		t.Setenv("VAULT_NAMESPACE", "env-ns")

		cfg := config.LoadConfig(base)

		assert.Equal(t, "env-ns", cfg.VaultNamespace)
	})

	t.Run("VAULTY_REFRESH_RATE from env", func(t *testing.T) {
		t.Setenv("VAULT_TOKEN", "env-token")
		t.Setenv("VAULTY_REFRESH_RATE", "60")

		cfg := config.LoadConfig(base)

		assert.Equal(t, 60, cfg.VaultyRefreshRate)
	})

	t.Run("invalid VAULTY_REFRESH_RATE falls back to default 30", func(t *testing.T) {
		t.Setenv("VAULT_TOKEN", "env-token")
		t.Setenv("VAULTY_REFRESH_RATE", "notanumber")

		cfg := config.LoadConfig(base)

		assert.Equal(t, 30, cfg.VaultyRefreshRate)
	})

	t.Run("VAULTY_LOG_FILE and VAULTY_LOG_LEVEL from env", func(t *testing.T) {
		t.Setenv("VAULT_TOKEN", "env-token")
		t.Setenv("VAULTY_LOG_FILE", "/tmp/test.log")
		t.Setenv("VAULTY_LOG_LEVEL", "debug")

		cfg := config.LoadConfig(base)

		assert.Equal(t, "/tmp/test.log", cfg.VaultyLogFile)
		assert.Equal(t, "debug", cfg.VaultyLogLevel)
	})

	t.Run("TLS env vars are loaded", func(t *testing.T) {
		t.Setenv("VAULT_TOKEN", "env-token")
		t.Setenv("VAULT_CACERT", "/etc/ssl/ca.pem")
		t.Setenv("VAULT_CLIENT_CERT", "/etc/ssl/client.pem")
		t.Setenv("VAULT_CLIENT_KEY", "/etc/ssl/client.key")

		cfg := config.LoadConfig(base)

		assert.Equal(t, "/etc/ssl/ca.pem", cfg.VaultCaCert)
		assert.Equal(t, "/etc/ssl/client.pem", cfg.VaultClientCert)
		assert.Equal(t, "/etc/ssl/client.key", cfg.VaultClientKey)
	})
}

func TestLoadConfig_Priority(t *testing.T) {
	t.Run("env var wins over yaml for VAULT_ADDR", func(t *testing.T) {
		cfgFile := writeConfig(t, `
vault_addr: https://yaml.vault.example.com
vault_token: yaml-token
`)
		t.Setenv("VAULT_ADDR", "https://from-env.vault.example.com")
		t.Setenv("VAULT_TOKEN", "env-token")

		cfg := config.LoadConfig(cfgFile)

		assert.Equal(t, "https://from-env.vault.example.com", cfg.VaultAddr)
	})

	t.Run("env var wins over yaml for VAULT_TOKEN", func(t *testing.T) {
		cfgFile := writeConfig(t, `
vault_addr: https://yaml.vault.example.com
vault_token: from-yaml
`)
		t.Setenv("VAULT_TOKEN", "from-env")

		cfg := config.LoadConfig(cfgFile)

		assert.Equal(t, "from-env", cfg.VaultToken)
	})
}
