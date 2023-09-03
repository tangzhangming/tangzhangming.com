package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


var logger *zap.Logger

type Field = zap.Field

// SugaredLogger将基本的Logger功能封装在一个较慢但不那么冗长的API中。任何Logger都可以通过其Sugar方法转换为sugardlogger。
// var sugarLogger *zap.SugaredLogger

func InitLogger() {
	// NewCore创建一个向WriteSyncer写入日志的Core。
	core := zapcore.NewCore(getEncoder(), getLogWriter(), zapcore.DebugLevel)
	logger = zap.New(core)
}

func getEncoder() zapcore.Encoder {
	// NewJSONEncoder创建了一个快速、低分配的JSON编码器。编码器适当地转义所有字段键和值。
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	// Create创建或截断指定文件。如果文件已经存在，它将被截断。如果该文件不存在，则以模式0666(在umask之前)创建。
	// 如果成功，返回的File上的方法可以用于IO;关联的文件描述符模式为O_RDWR。如果有一个错误，它的类型将是PathError。
	file, _ := os.Create("./logs/log.log")
	// AddSync converts an io.Writer to a WriteSyncer. It attempts to be
	// intelligent: if the concrete type of the io.Writer implements WriteSyncer,
	// we'll use the existing Sync method. If it doesn't, we'll add a no-op Sync.
	return zapcore.AddSync(file)
}

// 快捷调用zapcore
func String(key string, val string) Field { return zap.String(key, val) }
func Int(key string, val int) Field       { return zap.Int(key, val) }

// 快捷调用
func Debug(msg string, fields ...Field)  { logger.Debug(msg, fields...) }
func Info(msg string, fields ...Field)   { logger.Info(msg, fields...) }
func Warn(msg string, fields ...Field)   { logger.Warn(msg, fields...) }
func Error(msg string, fields ...Field)  { logger.Error(msg, fields...) }
func DPanic(msg string, fields ...Field) { logger.DPanic(msg, fields...) }
func Panic(msg string, fields ...Field)  { logger.Panic(msg, fields...) }
