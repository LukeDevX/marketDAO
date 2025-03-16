package logs

import (
	"RESTful-API/utils/config"
	"RESTful-API/utils/json"
	"bytes"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	log = logrus.New()
)

func InitLog() {
	//日志路径
	logPath := config.GetConfig("log.filePath").String()
	logName := fmt.Sprintf("%s.", logPath)
	//日志最大生存时间
	logAgeDays, _ := config.GetConfig("log.maxAge").Int()
	//日志级别
	logLevel, _ := logrus.ParseLevel(config.GetConfig("log.level").String())
	//日志模式
	logMode := config.GetConfig("log.mode").String()
	//日志格式
	logFormat := config.GetConfig("log.format").MustString("text")
	log.SetLevel(logLevel)
	r, _ := rotatelogs.New(logName + "%Y-%m-%d" + ".log")

	var mw io.Writer
	switch logMode {
	case "console":
		mw = io.MultiWriter(os.Stdout)
	case "file":
		mw = io.MultiWriter(r)
	case "file|console":
		fallthrough
	case "console|file":
		mw = io.MultiWriter(os.Stdout, r)
	default:
		mw = io.MultiWriter(os.Stdout)
	}
	//日志保存时间
	logAge := time.Duration(logAgeDays) * time.Second * 86400
	rotatelogs.WithMaxAge(logAge)
	log.SetOutput(mw)
	log.SetReportCaller(true)
	if logFormat == "json" {
		log.SetFormatter(&MyJsonTextFormatter{})
	} else {
		log.SetFormatter(&MyTextFormatter{})
	}
}

func Info(f interface{}, v ...any) {
	log.Info(formatLog(f, v...))
}

func Error(f interface{}, v ...any) {
	log.Error(formatLog(f, v...))
}

func Warn(f interface{}, v ...any) {
	log.Warn(formatLog(f, v...))
}

func Trace(f interface{}, v ...any) {
	log.Trace(formatLog(f, v...))
}

func Panic(f interface{}, v ...any) {
	log.Panic(formatLog(f, v...))
}

type MyTextFormatter struct{}

func (m *MyTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")
	var (
		logContent string
	)

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		_, fileName, line, _ := runtime.Caller(7)
		fName := filepath.Base(fileName)
		logContent = fmt.Sprintf("[%s] [%s] [%s:%d] %s\n", timestamp, ColorLevel(entry.Level.String()), fName, line, entry.Message)
	} else {
		logContent = fmt.Sprintf("[%s] [%s] %s\n", timestamp, ColorLevel(entry.Level.String()), entry.Message)
	}

	b.WriteString(logContent)
	return b.Bytes(), nil
}

// json格式
type MyJsonTextFormatter struct{}

func (m *MyJsonTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")
	var (
		logContent = make(map[string]string)
	)

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		_, fileName, line, _ := runtime.Caller(7)
		fName := filepath.Base(fileName)
		logContent["time"] = timestamp
		logContent["level"] = entry.Level.String()
		logContent["file"] = fmt.Sprintf("%s:%d", fName, line)
		logContent["msg"] = entry.Message
	} else {
		logContent["time"] = timestamp
		logContent["level"] = entry.Level.String()
		logContent["msg"] = entry.Message
	}

	jsonText, _ := json.Marshal(logContent)

	b.WriteString(string(jsonText) + "\n")
	return b.Bytes(), nil
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}
