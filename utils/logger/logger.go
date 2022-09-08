package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-demo/utils/config"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger
var loggerForBig *zap.Logger

func InitLoggerConfig() error {

	logService := InitServiceLogger()
	defer logService.Sync()

	InitBigLogger()

	return nil
}

// 初始化服务日志
func InitServiceLogger() *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:  "ts",
		LevelKey: "level",
		NameKey:  "logger",
		//CallerKey:      "caller",
		MessageKey: "msg",
		//StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     ZnTimeEncoder,                 //时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	// 设置日志级别（默认info级别，可以根据需要设置级别）
	var level zapcore.Level
	switch config.Config.Log.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	atom := zap.NewAtomicLevelAt(level)

	//服务日志配置
	configService := zap.Config{
		Level:            atom,               // 日志级别
		Development:      true,               // 开发模式，堆栈跟踪
		Encoding:         "json",             // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,      // 编码器配置
		OutputPaths:      []string{"stdout"}, // 日志写入文件的地址
		ErrorOutputPaths: []string{"stderr"}, // 将系统内的error记录到文件的地址
	}

	// 构建日志
	logger, _ := configService.Build()
	sugarLogger = logger.Sugar()

	return sugarLogger
}

func ZnTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

//func Debug(c *gin.Context, args ...interface{}) {
//	sugarLogger.Debug(args...)
//}

func Debugf(c *gin.Context, template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	sugarLogger.Debugw(msg, "req_id", c.GetString("req_id"), "req_source", c.GetString("req_source"))
}

//func Info(c *gin.Context, args ...interface{}) {
//	sugarLogger.Info(args...)
//}

func Infof(c *gin.Context, template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	sugarLogger.Infow(msg, "req_id", c.GetString("req_id"), "req_source", c.GetString("req_source"))

}

func SystemInfof(template string, args ...interface{}) {
	sugarLogger.Infof(template, args...)
}

func Infow(template string, args ...interface{}) {
	sugarLogger.Infow(template, args...)
}

//func Warn(c *gin.Context, args ...interface{}) {
//	sugarLogger.Warn(args...)
//}

func Warnf(c *gin.Context, template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	go Alert(msg, c.GetString("req_id"))
	sugarLogger.Warnw(msg, "req_id", c.GetString("req_id"), "req_source", c.GetString("req_source"))

}

func SystemWarnf(template string, args ...interface{}) {
	sugarLogger.Warnf(template, args...)
}

//func Error(c *gin.Context, args ...interface{}) {
//	sugarLogger.Error(args...)
//}

func Errorf(c *gin.Context, template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	go Alert(msg, c.GetString("req_id"))
	sugarLogger.Errorw(msg, "req_id", c.GetString("req_id"), "req_source", c.GetString("req_source"))
}

func DPanic(args ...interface{}) {
	sugarLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	sugarLogger.DPanicf(template, args...)
}

//func Panic(args ...interface{}) {
//	sugarLogger.Panic(args...)
//}

func Panicf(template string, args ...interface{}) {
	sugarLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	sugarLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	sugarLogger.Fatalf(template, args...)
}

// 初始化大数据日志
func InitBigLogger() {
	// 设置一些基本日志格式
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "msg",
		//LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	//固定路径
	infoWriter := getWriter("/www/arachnia_log/" + config.Config.Log.AppKey + "/bigdata")
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), zap.InfoLevel),
	)

	loggerForBig = zap.New(core, zap.AddCaller())
}

func getWriter(filename string) io.Writer {
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H"+".log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func BigInfow(template string) {
	loggerForBig.Info(template)
}

type Text struct {
	Content string `json:"content"`
}

func Alert(msg, reqId string) {
	if config.Config.Env != "release" {
		return
	}
	host, _ := os.Hostname()
	msg = "host:" + host + "\n" + "req_id:" + reqId + "\n" + msg
	param := make(map[string]interface{}, 2)
	param["msgtype"] = "text"
	param["text"] = Text{Content: msg}
	var jsonData []byte
	jsonData, _ = json.Marshal(param)
	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + config.Config.WarningKey
	PostJson(url, jsonData)
}

func PostJson(url string, params []byte) ([]byte, error) {
	body := bytes.NewBuffer(params)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	responseByte, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(responseByte))
	}

	return responseByte, nil
}
