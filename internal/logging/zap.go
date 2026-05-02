package logging

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Mode controls the output encoding of the logger.
type Mode string

const (
	ModeConsole Mode = "console"
	ModeJSON    Mode = "json"
)

// NewZapLogger constructs a *zap.Logger configured for structured output.
// In ModeConsole it emits coloured, human-readable lines suitable for
// a terminal. In ModeJSON it emits newline-delimited JSON suitable for
// log aggregation pipelines (Loki, Elastic, etc.).
//
// level is the minimum severity to emit (e.g. zap.DebugLevel).
func NewZapLogger(mode Mode, level zapcore.Level) *zap.Logger {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // overridden for JSON below
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var enc zapcore.Encoder

	switch mode {
	case ModeJSON:
		encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder // no ANSI in JSON
		enc = zapcore.NewJSONEncoder(encCfg)
	default: // ModeConsole
		encCfg.EncodeCaller = nil
		encCfg.EncodeName = func(name string, enc zapcore.PrimitiveArrayEncoder) {
			if len(name) < 20 {
				name += strings.Repeat(" ", 20-len(name))
			} else {
				name = name[:20]
			}
			enc.AppendString(name)
		}
		enc = zapcore.NewConsoleEncoder(encCfg)
	}

	sink := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(enc, sink, zap.NewAtomicLevelAt(level))

	return zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}
