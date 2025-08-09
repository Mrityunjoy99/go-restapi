package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	WithContext(ctx context.Context) Logger
	Flush() error
}

type zapLogger struct {
	logger *zap.Logger
}

var (
	once sync.Once
)

// NewZapLogger creates a new zap logger instance.
// serviceName will be added as a field to all log messages.
// It's safe to call this function multiple times - it will only initialize the logger once.
// ensureLogsDir ensures the logs directory exists and returns the path to the daily log file
func ensureLogsDir() (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Create logs directory if it doesn't exist
	logsDir := filepath.Join(cwd, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Generate log file name with current date
	now := time.Now()
	logFileName := fmt.Sprintf("log-%s.log", now.Format("2006-01-02"))
	return filepath.Join(logsDir, logFileName), nil
}

func NewZapLogger(serviceName string) (Logger, error) {
	var (
		zl      *zap.Logger
		initErr error
	)

	once.Do(func() {
		// Get log file path and ensure logs directory exists first
		logFilePath, err := ensureLogsDir()
		if err != nil {
			initErr = fmt.Errorf("failed to setup logs directory: %w", err)
			return
		}

		config := zap.NewProductionConfig()

		// Customize the encoder config
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		config.EncoderConfig = encoderConfig

		// Configure output paths
		config.OutputPaths = []string{"stdout", logFilePath}
		config.ErrorOutputPaths = []string{"stderr", logFilePath}

		// Set the log level based on environment variable or use InfoLevel as default
		switch os.Getenv("LOG_LEVEL") {
		case "debug":
			config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		case "warn":
			config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		case "error":
			config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		default:
			config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		}

		zl, initErr = config.Build(zap.AddCallerSkip(1))
		if initErr != nil {
			initErr = fmt.Errorf("failed to build logger: %w", initErr)
			return
		}

		// Add service name to all logs
		zl = zl.With(zap.String("service", serviceName))
		zap.ReplaceGlobals(zl)
	})

	if initErr != nil {
		return nil, initErr
	}

	if zl == nil {
		return nil, fmt.Errorf("failed to initialize logger")
	}

	return &zapLogger{logger: zl}, nil
}

func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}

	if ctx.Value("request_id") == nil {
		return l
	}

	l.logger = l.logger.With(zap.String("request_id", ctx.Value("request_id").(string)))

	return l
}

func (l *zapLogger) Flush() error {
	return l.logger.Sync()
}
