package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func SetupLogger(logLevel string, logFileName string) (*os.File, *zerolog.Logger) {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var logger zerolog.Logger

	// Default level for this example is info, unless debug flag is present
	if strings.EqualFold(logLevel, "debug") {
		level, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			log.Fatal().Err(err).Msg("Invalid log level")
		}
		zerolog.SetGlobalLevel(level)
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		// If debugOn is false, discard all log messages
		logger = zerolog.Nop()
	}

	var logFile *os.File

	// Check if file for logging is set

	if logFileName != "" {
		logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Panic().Err(err).Msg("Error opening log file")
		}
		logger = logger.Output(zerolog.ConsoleWriter{Out: logFile, TimeFormat: zerolog.TimeFieldFormat})
	}

	return logFile, &logger
}

type Config struct {
	VaultAddr      string `yaml:"vault_addr"`
	VaultNamespace string `yaml:"vault_namespace"`
	VaultToken     string `yaml:"vault_token"`
	VaultyLogFile  string `yaml:"vaulty_log_file"`
	VaultyLogLevel string `yaml:"vaulty_log_level"`
}

func checkForVaultAddress() {
	if os.Getenv("VAULT_ADDR") == "" {
		fmt.Println("VAULT_ADDR is not set. Please set it and try again.")
		os.Exit(1)
	}

	if os.Getenv("VAULT_TOKEN") == "" {
		fmt.Println("VAULT_TOKEN is not set. Please set it and try again.")
		os.Exit(1)
	}

}

func LoadConfig() Config {
	// Load the config from the YAML file
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory")
		os.Exit(1)
	}

	yamlFilePath := filepath.Join(home, ".vaul7y.yaml")

	data, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %v\n", err)
		os.Exit(1)
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
	if vaultyLogFile := os.Getenv("VAULTY_LOG_FILE"); vaultyLogFile != "" {
		config.VaultyLogFile = vaultyLogFile
	}
	if vaultyLogLevel := os.Getenv("VAULTY_LOG_LEVEL"); vaultyLogLevel != "" {
		config.VaultyLogLevel = vaultyLogLevel
	}

	if config.VaultAddr == "" {
		fmt.Println("VAULT_ADDR is not set. Please set it and try again.")
		os.Exit(1)
	}

	if config.VaultToken == "" {
		fmt.Println("VAULT_TOKEN is not set. Please set it and try again.")
		os.Exit(1)
	}

	return config
}
