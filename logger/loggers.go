package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var loggers = make(map[string]*zap.SugaredLogger)

func wrapLumberjack(level zapcore.Level, fileName string) func(core zapcore.Core) zapcore.Core {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     10,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		level,
	)

	return func(core2 zapcore.Core) zapcore.Core {
		return core
	}
}

func buildBaseLogger(level zapcore.Level, fileName string, country string) (*zap.SugaredLogger, error) {
	logFile := fmt.Sprintf("%s/%s/%s", os.Getenv("LOG_DIRECTORY"), country, fileName)

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logFile}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.DisableStacktrace = false
	cfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := cfg.Build(zap.WrapCore(wrapLumberjack(level, logFile)))

	if err != nil {
		return nil, err
	}

	createdLogger := logger.Sugar()

	err = createdLogger.Sync()

	if err != nil {
		return nil, err
	}

	return createdLogger, nil
}

func buildInfoLogger(country string) error {
	logger, err := buildBaseLogger(zap.InfoLevel, "info.log", country)

	loggers[fmt.Sprintf("info_%s", country)] = logger

	return err
}

func buildErrorLogger(country string) error {
	logger, err := buildBaseLogger(zap.ErrorLevel, "error.log", country)

	loggers[fmt.Sprintf("error_%s", country)] = logger

	return err
}

func buildWarningLogger(country string) error {
	logger, err := buildBaseLogger(zap.WarnLevel, "warn.log", country)

	loggers[fmt.Sprintf("warning_%s", country)] = logger

	return err
}

func BuildLoggers(countries []string) error {
	envLogDir := os.Getenv("LOG_DIRECTORY")

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, c := range countries {
		logDir := fmt.Sprintf("%s/%s/%s", path, envLogDir, c)

		_, err = os.Stat(logDir)

		if err != nil && os.IsNotExist(err) {
			err := os.MkdirAll(logDir, os.ModePerm)

			if err != nil {
				return err
			}
		} else if err != nil && !os.IsExist(err) {
			return err
		}

		if err := buildInfoLogger(c); err != nil {
			return err
		}

		if err := buildErrorLogger(c); err != nil {
			return err
		}

		if err := buildWarningLogger(c); err != nil {
			return err
		}
	}

	return nil
}

func Info(country string, msg ...interface{}) {
	loggers[fmt.Sprintf("info_%s", country)].Info(msg)
}

func Error(country string, msg ...interface{}) {
	loggers[fmt.Sprintf("error_%s", country)].Error(msg)
}

func Warn(country string, msg ...interface{}) {
	loggers[fmt.Sprintf("warning_%s", country)].Warn(msg)
}
