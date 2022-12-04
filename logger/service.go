package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var infoLogger *zap.SugaredLogger
var errorLogger *zap.SugaredLogger
var warningLogger *zap.SugaredLogger

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

func buildBaseLogger(level zapcore.Level, fileName string) (*zap.SugaredLogger, error) {
	logFile := fmt.Sprintf("%s/%s", os.Getenv("LOG_DIRECTORY"), fileName)

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

func buildInfoLogger() error {
	logger, err := buildBaseLogger(zap.InfoLevel, "info.log")

	infoLogger = logger

	return err
}

func buildErrorLogger() error {
	logger, err := buildBaseLogger(zap.ErrorLevel, "error.log")

	errorLogger = logger

	return err
}

func buildWarningLogger() error {
	logger, err := buildBaseLogger(zap.WarnLevel, "warn.log")

	warningLogger = logger

	return err
}

func BuildLoggers() error {
	envLogDir := os.Getenv("LOG_DIRECTORY")

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	logDir := fmt.Sprintf("%s/%s", path, envLogDir)
	
	_, err = os.Stat(logDir)

	if os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)

		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	if err := buildInfoLogger(); err != nil {
		return err
	}

	if err := buildErrorLogger(); err != nil {
		return err
	}

	if err := buildWarningLogger(); err != nil {
		return err
	}

	return nil
}

func Info(msg ...interface{}) {
	if os.Getenv("APP_ENV") != "prod" && os.Getenv("APP_ENV") != "staging" {
		fmt.Println(fmt.Sprintf("INFO: %v", msg))
	}

	infoLogger.Info(msg)
}

func Error(msg ...interface{}) {
	if os.Getenv("APP_ENV") != "prod" && os.Getenv("APP_ENV") != "staging" {
		fmt.Println(fmt.Sprintf("ERROR: %v", msg))
	}

	errorLogger.Error(msg)
}

func Warn(msg ...interface{}) {
	if os.Getenv("APP_ENV") != "prod" && os.Getenv("APP_ENV") != "staging" {
		fmt.Println(fmt.Sprintf("WARNING: %v", msg))
	}

	warningLogger.Warn(msg)
}
