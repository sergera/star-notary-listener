package logger

import (
	"log"
	"time"

	"github.com/sergera/star-notary-listener/internal/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	setLogger()
}

func setLogger() {
	loggerInstance, err := new()
	if err != nil {
		log.Panicf("Could not instance logger: %+v\n", err)
	}

	logger = loggerInstance
}

func new() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		env.LogPath + "starnotary.log",
		"stderr",
	}
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}
	/* remove sampling CPU & I/O load cap so no logs are lost in concurrent use */
	cfg.Sampling = nil
	/* build with caller skip to log caller from client code */
	return cfg.Build(zap.AddCallerSkip(1))
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
