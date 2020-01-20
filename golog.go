package golog

import (
	"encoding/json"
	"errors"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type GoLog struct {
	zapLog     *zap.Logger
	Path       string `json:"path"`        //日志文件保存路径 stdout | .log
	Level      string `json:"level"`       //打印的日志级别
	MaxSize    int    `json:"max_size"`    //切割大小
	MaxAge     int    `json:"max_age"`     //保留天数
	MaxBackups int    `json:"max_backups"` //最大备份数
	Caller     bool   `json:"caller"`      //打印行号
	Marshal    bool   `json:"marshal"`     //是否json格式化
	Compress   bool   `json:"compress"`    //是否压缩
}

func LoadConfig(config ...string) (*GoLog, error) {
	log := &GoLog{
		Path:       "stdout",
		Level:      "info",
		MaxSize:    128,
		MaxAge:     30,
		MaxBackups: 7,
		Caller:     true,
		Marshal:    false,
	}
	if len(config) > 0 {
		if err := json.Unmarshal([]byte(config[0]), log); err != nil {
			return nil, err
		}
		if log.Path == "" {
			return nil, errors.New("日志文件路径未配置")
		}
	}
	zapConf := zap.NewProductionEncoderConfig()
	zapConf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	var enc zapcore.Encoder
	var ws zapcore.WriteSyncer
	var enable zapcore.LevelEnabler
	if log.Marshal {
		enc = zapcore.NewJSONEncoder(zapConf)
	} else {
		enc = zapcore.NewConsoleEncoder(zapConf)
	}
	switch log.Path {
	case "stdout":
		ws = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	case "stderr":
		ws = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr))
	case "":
		ws = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	default:
		ws = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&lumberjack.Logger{
			Filename:   log.Path,
			MaxAge:     log.MaxAge,
			MaxSize:    log.MaxSize,
			MaxBackups: log.MaxBackups,
			Compress:   log.Compress,
		}))
	}
	switch log.Level {
	case "debug":
		enable = zapcore.DebugLevel
	case "info":
		enable = zapcore.InfoLevel
	case "warn":
		enable = zapcore.WarnLevel
	case "error":
		enable = zapcore.ErrorLevel
	default:
		enable = zapcore.InfoLevel
	}
	core := zapcore.NewCore(enc, ws, enable)
	if log.Caller {
		log.zapLog = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		log.zapLog = zap.New(core)
	}

	return log, nil
}
