package config

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"gopkg.in/yaml.v2"
)

type Config struct {
	VaultAddr         string      `yaml:"vault_addr"`
	VaultNamespace    string      `yaml:"vault_namespace"`
	VaultToken        string      `yaml:"vault_token"`
	VaultCaCert       string      `yaml:"vault_cacert"`
	VaultClientCert   string      `yaml:"vault_client_cert"`
	VaultClientKey    string      `yaml:"vault_client_key"`
	VaultyLogFile     string      `yaml:"vaulty_log_file"`
	VaultyLogLevel    string      `yaml:"vaulty_log_level"`
	VaultyRefreshRate int         `yaml:"vaulty_refresh_rate"`
	Theme             ThemeConfig `yaml:"theme"`
}

type ThemeConfig struct {
	Background         string `yaml:"background"`
	HighlightPrimary   string `yaml:"highlight_primary"`
	HighlightSecondary string `yaml:"highlight_secondary"`
	Standard           string `yaml:"standard"`
	Active             string `yaml:"active"`
	White              string `yaml:"white"`
	LightGrey          string `yaml:"light_grey"`
	ModalInfo          string `yaml:"modal_info"`
	Attention          string `yaml:"attention"`
}

func LoadConfig(cfgFile string) Config {
	var config Config
	// Load the config from the YAML file
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting user home directory")
	}

	var data []byte
	if cfgFile == "" {
		yamlFilePath := filepath.Join(home, ".vaul7y.yaml")
		if _, err := os.Stat(yamlFilePath); !os.IsNotExist(err) {
			data, err = os.ReadFile(yamlFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading YAML file: %v\n", err)
			}
		}
	} else {
		data, err = os.ReadFile(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading YAML file: %v\n", err)
		}
	}

	if data != nil {
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing YAML file: %v\n", err)
		}
	}

	// Check for vault cache
	home, err = os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting user home directory")
	} else {
		vaultTokenPath := filepath.Join(home, ".vault-token")
		if _, err := os.Stat(vaultTokenPath); err == nil {
			data, err := os.ReadFile(vaultTokenPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault token file: %v\n", err)
			} else {
				config.VaultToken = string(data)
			}
		}
	}

	// Overwrite with environment variables if they are set
	if vaultAddr := os.Getenv("VAULT_ADDR"); vaultAddr != "" {
		config.VaultAddr = vaultAddr
	}
	if vaultNamespace := os.Getenv("VAULT_NAMESPACE"); vaultNamespace != "" {
		config.VaultNamespace = vaultNamespace
	}
	if vaultToken := os.Getenv("VAULT_TOKEN"); vaultToken != "" {
		config.VaultToken = vaultToken
	}
	if vaultCaCert := os.Getenv("VAULT_CACERT"); vaultCaCert != "" {
		config.VaultCaCert = vaultCaCert
	}
	if vaultClientCert := os.Getenv("VAULT_CLIENT_CERT"); vaultClientCert != "" {
		config.VaultClientCert = vaultClientCert
	}
	if vaultClientKey := os.Getenv("VAULT_CLIENT_KEY"); vaultClientKey != "" {
		config.VaultClientKey = vaultClientKey
	}
	if vaultyLogFile := os.Getenv("VAULTY_LOG_FILE"); vaultyLogFile != "" {
		config.VaultyLogFile = vaultyLogFile
	}
	if vaultyLogLevel := os.Getenv("VAULTY_LOG_LEVEL"); vaultyLogLevel != "" {
		config.VaultyLogLevel = vaultyLogLevel
	}
	if vaultyRefreshRate := os.Getenv("VAULTY_REFRESH_RATE"); vaultyRefreshRate != "" {
		vaultyRefreshRateInt, err := strconv.Atoi(vaultyRefreshRate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting VAULTY_REFRESH_RATE to int: %v", err)
		} else {
			config.VaultyRefreshRate = vaultyRefreshRateInt
		}
	}

	if config.VaultToken == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting user home directory")
		} else {
			vaultTokenPath := filepath.Join(home, ".vault-token")
			if _, err := os.Stat(vaultTokenPath); err == nil {
				data, err := os.ReadFile(vaultTokenPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading vault token file: %v\n", err)
				} else {
					config.VaultToken = string(data)
				}
			}
		}
	}

	if config.VaultAddr == "" {
		fmt.Fprintln(os.Stderr, "VAULT_ADDR is not set. Please set it and try again.")
		os.Exit(1)
	}

	if config.VaultToken == "" {
		fmt.Fprintln(os.Stderr, "VAULT_TOKEN is not set. Please set it and try again.")
		os.Exit(1)
	}

	if config.VaultyRefreshRate == 0 {
		config.VaultyRefreshRate = 30
	}

	if strings.EqualFold(config.VaultyLogLevel, "debug") {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGTERM)

			<-ch
			fmt.Fprintln(os.Stderr, "Dumping goroutines")
			bufsize := int(10 * 1024 * 1024) // 10 MiB
			buf := make([]byte, bufsize)
			n := runtime.Stack(buf, true)
			filename := fmt.Sprintf("%s.dump", config.VaultyLogFile)

			os.WriteFile(filename, buf[:n], 0644)
			os.Exit(1)
		}()

	}

	return config
}
