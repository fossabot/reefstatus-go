package log

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"time"
)

type Level string

const (
	INFO    = "Info"
	DEBUG   = "Debug"
	WARNING = "Warning"
	ERROR   = "Error"
	FATAL   = "Fatal"
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

func printLog(text string, level Level) {
	formattedText := fmt.Sprintf("[%s] %d - %s", level, time.Now().UnixNano()/1000000, text)
	fmt.Println(formattedText)
}
