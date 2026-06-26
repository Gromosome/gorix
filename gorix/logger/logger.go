package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorReset  = "\033[0m"
)

type logType int

const (
	errorLog logType = 0
	infoLog  logType = 1
	warnLog  logType = 2
	debugLog logType = 3
)

type CallerLevel int

func (level CallerLevel) Int() int {
	return int(level)
}

const (
	CallerLevelDirect CallerLevel = 0
	CallerLevelParent CallerLevel = 1
	CallerLevelCaller CallerLevel = 2
	CallerLevelOrigin CallerLevel = 3
)

type Logger struct {
	logType     logType
	callerLevel int
}

func NewLogger(callerLevel CallerLevel) *Logger {
	return &Logger{
		logType:     errorLog,
		callerLevel: callerLevel.Int(),
	}

}
func (l Logger) CallerLevel(callerLevel CallerLevel) *Logger {
	l.callerLevel = callerLevel.Int()
	return &l
}
func (l Logger) TypeInfo() *Logger {
	l.logType = infoLog
	return &l
}
func (l Logger) TypeError() *Logger {
	l.logType = errorLog
	return &l
}
func (l Logger) TypeWarn() *Logger {
	l.logType = warnLog
	return &l
}
func (l Logger) TypeDebug() *Logger {
	l.logType = debugLog
	return &l
}

var mu sync.Mutex

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

func callerInfo(skip int) (file string, line int, fn string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0, "unknown"
	}

	runtimeFunc := runtime.FuncForPC(pc)
	if runtimeFunc == nil {
		return file, line, "unknown"
	}

	return file, line, runtimeFunc.Name()
}

func writeLog(
	out io.Writer,
	color string,
	level string,
	skip int,
	format string,
	args ...any,
) {
	mu.Lock()
	defer mu.Unlock()

	file, line, fn := callerInfo(skip)
	msg := fmt.Sprintf(format, args...)

	fmt.Fprintf(
		out,
		"%s[%s] %s %s:%d %s: %s%s\n",
		color,
		timestamp(),
		level,
		file,
		line,
		fn,
		msg,
		colorReset,
	)
}
func (l *Logger) Log(format string, args ...any) {
	skip := l.callerLevel + 3
	switch l.logType {
	case errorLog:
		writeLog(os.Stderr, colorRed, "ERROR", skip, format, args...)
	case infoLog:
		writeLog(os.Stdout, colorGreen, "INFO", skip, format, args...)
	case warnLog:
		writeLog(os.Stderr, colorYellow, "WARN", skip, format, args...)
	case debugLog:
		writeLog(os.Stdout, colorCyan, "DEBUG", skip, format, args...)
	default:
		writeLog(os.Stdout, colorRed, "ERROR", skip, format, args...)
	}
}
