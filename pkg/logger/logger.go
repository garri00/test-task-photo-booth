package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"

	"test-task-photo-booth/src/config"
)

const (
	LogLevelDebug   = "DEBUG"
	LogLevelInfo    = "INFO"
	LogLevelTrace   = "TRACE"
	LogLevelPanic   = "PANIC"
	LogLevelNoLevel = "NOLEVEL"
	LogLevelError   = "ERROR"
	LogLevelFatal   = "FATAL"
	LogLevelWarn    = "WARN"
)

const (
	TimeFormat         = "Monday, 02-01-2006 15:04:05 MST"
	FileOpenPermission = 0664

	logPath         = "logs"
	mainLogFileName = "main.log"
)

var Log = zerolog.New(os.Stdout).With().Timestamp().Logger()

//var TestLog = zerolog.New(os.Stdout).With().Timestamp().Logger().Level(zerolog.InfoLevel)

func SetLogger(c config.Configs) error {
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		return fmt.Errorf("os.MkdirAll() failed: %w", err)
	}

	//Set log filePath
	filePath := filepath.Join(logPath, mainLogFileName)

	runLogFile, err := os.OpenFile(
		filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		FileOpenPermission,
	)
	if err != nil {
		return fmt.Errorf("cannot open log file: %w", err)
	}

	//Set global logger time format
	zerolog.TimeFieldFormat = TimeFormat

	multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile)

	Log = zerolog.New(multi).With().Timestamp().Logger()

	//Set log level info
	logLevel := getZerologLogLevel(c.LogLevel)
	Log = Log.Level(logLevel)

	Log.Info().Any("log level", Log.GetLevel().String()).Msg("Log level set")

	return nil
}

// SetServiceLogger creates additional logger
func SetServiceLogger(serviceLogName string, c config.Configs) (*zerolog.Logger, error) {
	var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		return &logger, fmt.Errorf("os.MkdirAll() failed: %w", err)
	}

	loggerPath := filepath.Join(logPath, serviceLogName+".log")

	//set log file
	runLogFile, err := os.OpenFile(
		loggerPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		FileOpenPermission,
	)
	if err != nil {
		return &logger, fmt.Errorf("cannot open log file: %w", err)
	}

	multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile)

	logger = zerolog.New(multi).With().Timestamp().Logger()

	logLevel := getZerologLogLevel(c.LogLevel)
	logger = logger.Level(logLevel)

	logger.Info().Any("log level", logger.GetLevel().String()).Msgf("Service log: %s, set level", serviceLogName)

	return &logger, nil
}

func getZerologLogLevel(logLevel string) zerolog.Level {
	zerologLoglevel := zerolog.InfoLevel

	switch logLevel {
	case LogLevelTrace:
		zerologLoglevel = zerolog.TraceLevel
	case LogLevelDebug:
		zerologLoglevel = zerolog.DebugLevel
	case LogLevelInfo:
		zerologLoglevel = zerolog.InfoLevel
	case LogLevelWarn:
		zerologLoglevel = zerolog.WarnLevel
	case LogLevelError:
		zerologLoglevel = zerolog.ErrorLevel
	case LogLevelFatal:
		zerologLoglevel = zerolog.FatalLevel
	case LogLevelPanic:
		zerologLoglevel = zerolog.PanicLevel
	case LogLevelNoLevel:
		zerologLoglevel = zerolog.NoLevel
	}

	return zerologLoglevel
}
