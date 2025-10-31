package logger

import (
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createLogger() *zap.SugaredLogger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.StacktraceKey = "stacktrace"
	encoderCfg.LineEnding = zapcore.DefaultLineEnding

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stderr),
		zap.InfoLevel,
	)

	base := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	runtime.SetFinalizer(base, func(l *zap.Logger) {
		_ = l.Sync()
	})

	return base.Sugar()
}

func Info(msg string, args ...interface{}) {
	logger.Infow(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Warnw(msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger.Errorw(msg, args...)
}

var logger = createLogger()
