package log

import (
	"bytes"
	"fmt"
	"github.com/cjburchell/reefstatus-go/common"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"runtime"
	"time"
)

type Level struct {
	Text     string
	Severity int
}

var (
	DEBUG   = Level{Text: "Debug", Severity: 0}
	INFO    = Level{Text: "Info", Severity: 1}
	WARNING = Level{Text: "Warning", Severity: 2}
	ERROR   = Level{Text: "Error", Severity: 3}
	FATAL   = Level{Text: "Fatal", Severity: 4}
)

func getStackTrace() string {

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Stacktrace:\n"))
	var i int = 2
	for i < 40 {
		if function1, file1, line1, ok := runtime.Caller(i); ok {
			buffer.WriteString(fmt.Sprintf("      at %s (%s:%d)\n", runtime.FuncForPC(function1).Name(), file1, line1))
		} else {
			break
		}
		i++
	}

	return buffer.String()
}

func Warnf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), WARNING)
}

func Warn(v ...interface{}) {
	printLog(fmt.Sprint(v...), WARNING)
}

func Error(err error, v ...interface{}) {
	msg := fmt.Sprint(v...)
	if msg == "" {
		msg = fmt.Sprintf("Error: %s\n%s", err.Error(), getStackTrace())
	} else {
		msg = fmt.Sprintf("%s\nError: %s\n%s", msg, err.Error(), getStackTrace())
	}

	printLog(msg, ERROR)
}

func Errorf(err error, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if msg == "" {
		msg = fmt.Sprintf("Error: %s\n%s", err.Error(), getStackTrace())
	} else {
		msg = fmt.Sprintf("%s\nError: %s\n%s", msg, err.Error(), getStackTrace())
	}

	printLog(msg, ERROR)
}

func Fatal(v ...interface{}) {
	printLog(fmt.Sprint(v...), FATAL)
	log.Panic(v...)
}

func Fatalf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), FATAL)
	log.Panicf(format, v...)
}

func Debug(v ...interface{}) {
	printLog(fmt.Sprint(v...), DEBUG)
}

func Debugf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), DEBUG)
}

func Print(v ...interface{}) {
	printLog(fmt.Sprint(v...), INFO)
}

func Printf(format string, v ...interface{}) {
	printLog(fmt.Sprintf(format, v...), INFO)
}

var minLogLevel int
var logToConsole bool
var logToFile bool
var logger = lumberjack.Logger{
	MaxAge:     1,
	MaxBackups: 20,
}

func init() {
	minLogLevel = common.GetEnvInt("LOG_LEVEL", 1)
	logToConsole = common.GetEnvBool("LOG_CONSOLE", true)
	logToFile = common.GetEnvBool("LOG_FILE", true)
	logger.Filename = common.GetEnv("LOG_FILE_PATH", "e:\\data\\log\\server.log")
}

func printLog(text string, level Level) {
	if level.Severity >= minLogLevel {
		formattedText := fmt.Sprintf("[%s] %d - %s", level.Text, time.Now().UnixNano()/1000000, text)
		if logToConsole {
			fmt.Println(formattedText)
		}

		if logToFile {
			logger.Write([]byte(formattedText + "\n"))
		}
	}
}
