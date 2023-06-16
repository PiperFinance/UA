package conf

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.SugaredLogger

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func LoadLogger() {
	var zp *zap.Logger
	var err error
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	errFile := fmt.Sprintf("%s%s", Config.LogDir, "/err.log")
	debugFile := fmt.Sprintf("%s%s", Config.LogDir, "/debug.log")
	errLog, _, err := zap.Open(errFile)
	if err != nil {
		panic(err)
	}
	debugLog, _, err := zap.Open(debugFile)
	if err != nil {
		panic(err)
	}

	topicDebugging := debugLog
	topicErrors := errLog

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	fileEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	var core zapcore.Core
	if strings.ToLower(Config.LogLevel) == "debug" {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, topicErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(fileEncoder, topicDebugging, lowPriority),
			zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, topicErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		)
	}
	if Config.DEV {
		zp = zap.New(core, zap.Development(), zap.AddStacktrace(zap.WarnLevel))
	} else {
		zp = zap.New(core, zap.Development())
	}
	defer zp.Sync()
	errLJ := lumberjack.Logger{
		Filename:   errFile,
		MaxSize:    10, // MB
		MaxBackups: 1,
		MaxAge:     3, // days
		Compress:   true,
	}
	debugLJ := lumberjack.Logger{
		Filename:   debugFile,
		MaxSize:    50, // MB
		MaxBackups: 1,
		MaxAge:     3, // days
		Compress:   true,
	}
	zap.RegisterSink("ErrLumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &errLJ,
		}, nil
	})
	zap.RegisterSink("DebugLumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &debugLJ,
		}, nil
	})
	sugar := zp.Sugar()
	Logger = sugar
}
