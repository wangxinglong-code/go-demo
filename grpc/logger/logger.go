package rpcLogger

import (
	"encoding/json"
	"go-demo/utils/common"
	"go-demo/utils/config"
	"go-demo/utils/logger"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerForRpc *zap.Logger

type ZapLogger struct {
	logger *zap.Logger
}

func InitRcpLog() *zap.Logger {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		//CallerKey:    "",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	var level zapcore.Level
	switch config.Config.Mode {
	case "release":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, os.Stdout, level),
	)

	loggerForRpc = zap.New(core, zap.AddCaller())
	return loggerForRpc
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{
		logger: logger,
	}
}

func (zl *ZapLogger) Info(args ...interface{}) {
	zl.logger.Sugar().Info(args...)
}

func (zl *ZapLogger) Infoln(args ...interface{}) {
	zl.logger.Sugar().Info(args...)
}
func (zl *ZapLogger) Infof(format string, args ...interface{}) {
	zl.logger.Sugar().Infof(format, args...)
}

func (zl *ZapLogger) Warning(args ...interface{}) {
	zl.logger.Sugar().Warn(args...)
}

func (zl *ZapLogger) Warningln(args ...interface{}) {
	zl.logger.Sugar().Warn(args...)
}

func (zl *ZapLogger) Warningf(format string, args ...interface{}) {
	zl.logger.Sugar().Warnf(format, args...)
}

func (zl *ZapLogger) Error(args ...interface{}) {
	zl.logger.Sugar().Error(args...)
}

func (zl *ZapLogger) Errorln(args ...interface{}) {
	zl.logger.Sugar().Error(args...)
}

func (zl *ZapLogger) Errorf(format string, args ...interface{}) {
	zl.logger.Sugar().Errorf(format, args...)
}

func (zl *ZapLogger) Fatal(args ...interface{}) {
	zl.logger.Sugar().Fatal(args...)
}

func (zl *ZapLogger) Fatalln(args ...interface{}) {
	zl.logger.Sugar().Fatal(args...)
}

// Fatalf logs to fatal level
func (zl *ZapLogger) Fatalf(format string, args ...interface{}) {
	zl.logger.Sugar().Fatalf(format, args...)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (zl *ZapLogger) V(v int) bool {

	return false
}

func RpcInfow(reqSource, reqId, reqMethod string, reqTime time.Time, req, res interface{}) {
	reqBytes, _ := json.Marshal(req)
	resBytes, _ := json.Marshal(res)

	if len(resBytes) > common.LogLenDefault {
		resBytes = []byte("log length too long")
	}

	logger.Infow("",
		"log_type", "rpc_request_log",
		"request_time", reqTime.UnixNano()/1e6,
		"request_method", reqMethod,
		"request_body", string(reqBytes),
		"response_time", time.Now().UnixNano()/1e6,
		"response_body", string(resBytes),
		"cost_time", common.End(reqTime),
		"req_id", reqId,
		"req_source", reqSource,
	)
}
