package config

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
