package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init() {
	env := os.Getenv("ENV") // dev / prod

	var config zap.Config

	if env == "production" {
		// 🔴 PROD → JSON logs
		config = zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Encoding:    "json",

			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",

				EncodeLevel:  zapcore.LowercaseLevelEncoder,
				EncodeTime:   zapcore.EpochTimeEncoder,
				EncodeCaller: zapcore.ShortCallerEncoder,
			},

			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	} else {
		// 🟢 DEV → Pretty logs
		config = zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
			Development: true,
			Encoding:    "console",

			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:    "ts",
				LevelKey:   "level",
				CallerKey:  "caller",
				MessageKey: "msg",

				EncodeLevel:  zapcore.CapitalColorLevelEncoder,
				EncodeTime:   zapcore.ISO8601TimeEncoder,
				EncodeCaller: zapcore.ShortCallerEncoder,
			},

			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	var err error
	Log, err = config.Build(zap.AddCaller())
	if err != nil {
		panic(err)
	}
}

func Sync() {
	Log.Sync()
}