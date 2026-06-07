package logging

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// SetupLogging configures logging based on environment variables
func SetupLogging() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}

	// Configure zap logger
	config := zap.NewProductionConfig()

	// Set log level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "INFO":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "WARN", "WARNING":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "ERROR":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	// Configure output
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	// Build the logger
	logger, err := config.Build()
	if err != nil {
		// Fallback to simple logger if zap fails
		globalLogger = zap.NewNop()
		return
	}

	globalLogger = logger
	zap.ReplaceGlobals(logger)
}

// GetLogger returns a logger with the given name
func GetLogger(name string) *zap.Logger {
	if globalLogger == nil {
		SetupLogging()
	}
	return globalLogger.With(zap.String("logger", name))
}
