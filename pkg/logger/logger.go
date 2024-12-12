package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	encConfig := zap.NewDevelopmentEncoderConfig()

	// Customize level format with colors
	encConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.InfoLevel:
			enc.AppendString("\x1b[36mINFO\x1b[0m") // Cyan
		case zapcore.WarnLevel:
			enc.AppendString("\x1b[33mWARN\x1b[0m") // Yellow
		case zapcore.ErrorLevel:
			enc.AppendString("\x1b[31mERROR\x1b[0m") // Red
		case zapcore.DebugLevel:
			enc.AppendString("\x1b[32mDEBUG\x1b[0m") // Green
		default:
			enc.AppendString(l.CapitalString())
		}
	}

	// Customize time format
	encConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02T15:04:05.000") + "]")
	}

	// Customize caller format
	encConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(caller.TrimmedPath())
	}

	// Customize the final format
	encConfig.ConsoleSeparator = " "

	// Create custom encoder
	consoleEncoder := zapcore.NewConsoleEncoder(encConfig)

	// Create core
	core := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		getLogLevel(),
	)

	// Create logger with caller info
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func Info(msg string, fields ...zapcore.Field) {
	log.Info(msg, fields...)
}

func Error(msg string, err error, fields ...zapcore.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	log.Error(msg, fields...)
}

func Debug(msg string, fields ...zapcore.Field) {
	log.Debug(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	log.Warn(msg, fields...)
}

// Field constructors for convenience
func String(key, value string) zapcore.Field {
	return zap.String(key, value)
}

func Int(key string, value int) zapcore.Field {
	return zap.Int(key, value)
}

func Any(key string, value interface{}) zapcore.Field {
	return zap.Any(key, value)
}

func getLogLevel() zapcore.Level {
	if os.Getenv("ENV") == "production" {
		return zap.InfoLevel
	}
	return zap.DebugLevel
}
