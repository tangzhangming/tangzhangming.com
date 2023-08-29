package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 定义一个全局 logger 实例
// Logger提供快速、分级、结构化的日志记录。所有方法对于并发使用都是安全的。
// Logger是为每一微秒和每一个分配都很重要的上下文设计的，
// 因此它的API有意倾向于性能和类型安全，而不是简便性。
// 对于大多数应用程序，SugaredLogger在性能和人体工程学之间取得了更好的平衡。
var logger *zap.Logger

// SugaredLogger将基本的Logger功能封装在一个较慢但不那么冗长的API中。任何Logger都可以通过其Sugar方法转换为sugardlogger。
//与Logger不同，SugaredLogger并不坚持结构化日志记录。对于每个日志级别，它公开了四个方法:
//   - methods named after the log level for log.Print-style logging
//   - methods ending in "w" for loosely-typed structured logging
//   - methods ending in "f" for log.Printf-style logging
//   - methods ending in "ln" for log.Println-style logging

// For example, the methods for InfoLevel are:
//
//	Info(...any)           Print-style logging
//	Infow(...any)          Structured logging (read as "info with")
//	Infof(string, ...any)  Printf-style logging
//	Infoln(...any)         Println-style logging
var sugarLogger *zap.SugaredLogger

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	// NewCore创建一个向WriteSyncer写入日志的Core。

	// A WriteSyncer is an io.Writer that can also flush any buffered data. Note
	// that *os.File (and thus, os.Stderr and os.Stdout) implement WriteSyncer.

	// LevelEnabler决定在记录消息时是否启用给定的日志级别。
	// Each concrete Level value implements a static LevelEnabler which returns
	// true for itself and all higher logging levels. For example WarnLevel.Enabled()
	// will return true for WarnLevel, ErrorLevel, DPanicLevel, PanicLevel, and
	// FatalLevel, but return false for InfoLevel and DebugLevel.
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	// New constructs a new Logger from the provided zapcore.Core and Options. If
	// the passed zapcore.Core is nil, it falls back to using a no-op
	// implementation.
	logger = zap.New(core)
	// Sugar封装了Logger，以提供更符合人体工程学的API，但速度略慢。糖化一个Logger的成本非常低，
	// 因此一个应用程序同时使用Loggers和SugaredLoggers是合理的，在性能敏感代码的边界上在它们之间进行转换。
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	// NewJSONEncoder创建了一个快速、低分配的JSON编码器。编码器适当地转义所有字段键和值。
	// NewProductionEncoderConfig returns an opinionated EncoderConfig for
	// production environments.
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

func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}
