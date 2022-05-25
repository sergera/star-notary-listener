package logger

import (
	"os"

	"github.com/sergera/star-notary-listener/internal/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

/* function variables for zap field types */
var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	Int8p       = zap.Int8p
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Object      = zap.Object
	Inline      = zap.Inline
	Any         = zap.Any
)

/* type variables for zap encoder types */
type (
	ObjectEncoder         = zapcore.ObjectEncoder
	ArrayEncoder          = zapcore.ArrayEncoder
	PrimitiveArrayEncoder = zapcore.PrimitiveArrayEncoder
)

func Init() {
	logger = newLogger()
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newLogger() *zap.Logger {
	consoleEncoder := zapcore.NewConsoleEncoder(newEncoderConfig())
	fileEncoder := zapcore.NewJSONEncoder(newEncoderConfig())

	consoleSyncer := zapcore.AddSync(os.Stderr)
	fileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   env.LogPath + "starnotary.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
	})

	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleSyncer, levelEnabler),
		zapcore.NewCore(fileEncoder, fileSyncer, levelEnabler),
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}

func Debug(msg string, fields ...zapcore.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

func Panic(msg string, fields ...zapcore.Field) {
	logger.Panic(msg, fields...)
}

func DPanic(msg string, fields ...zapcore.Field) {
	logger.DPanic(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	logger.Fatal(msg, fields...)
}
