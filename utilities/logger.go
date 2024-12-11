package utilities

import (
	"fmt"
	"sync"
	"time"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapSession is a wrapper for zap.Logger to provide a single instance of the logger.
type zapSession struct {
	zap     	*zap.Logger // zap logger instance
	once        sync.Once   // once is used to make sure that the session is created only once
	environment string      // environment is the environment the logger is running in
	serviceName string      // serviceName is the name of the service the logger is running in
}

// Session returns a zap logger instance.
func (ctx *zapSession) Session(environment string, serviceName string) *zap.Logger {

	ctx.environment = environment
	ctx.serviceName = serviceName

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("failed to create logger: %v\n", r)
		}
		ctx.close()
	}()

	ctx.once.Do(func() {
		config, _, err := ctx.loadConfiguration()
		if err != nil {
			fmt.Printf("failed to create logger: %v\n", err)
			return
		}
		ctx.zap = zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(*config),
			zapcore.AddSync(colorable.NewColorableStdout()),
			ctx.logLeveler(),
		))
		println("")
		// Log the time it took to initialize the logger.
		ctx.zap.Sugar().Debugf("Logger Initialized ==< %v ms >== | Service ==< %v >== ", zap.Duration("took", time.Since(time.Now())).Integer, ctx.serviceName)
	})
	return ctx.zap
}

// config returns the encoder config and atomic level for the logger.
func (ctx *zapSession) loadConfiguration() (*zapcore.EncoderConfig, *zap.AtomicLevel, error) {
	atom := zap.NewAtomicLevel()
	config := zap.NewDevelopmentEncoderConfig()
	config.TimeKey = "timestamp"
	config.FunctionKey = "func"
	config.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return &config, &atom, nil
}

// close closes the logger.
func (ctx *zapSession) close() error {
	return ctx.zap.Sync()
}

// logLeveler returns the log level for the logger.
func (ctx *zapSession) logLeveler() zapcore.LevelEnabler {
	if ctx.environment == "Dev" || ctx.environment == "dev" {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

// Logger is a global instance of zapSession.
var Logger = &zapSession{}
