package gotils

//go:generate enumer -type LogLevel loglevel.go

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// LogLevel sets the criticallity of a logging output. It is used to filter logging messages
// depending on their priority. Compare to log4j.
type LogLevel int

// The predefined LogLevels that are used by the logging funktions below.
const (
	Off LogLevel = iota
	Fatal
	Error
	Warn
	Debug
	Info
	All
)

var loggerflags = log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds | log.LUTC
var conveninceLogger *Logger

// A Logger is an onbject the offers several method to write Messages to a stream.
// Atually it is a wrapper aroung the 'log' package, that enhances the LogLevel functionality.
type Logger struct {
	internallogger *log.Logger
	ActiveLoglevel LogLevel
}

// NewLoggerFromFile creates a new Logger. It take a file parameter (io.Writer) output file
// and a LogLevel to filter the messages that are wanted.
// The first created logger wil be set to be the convenience logger (see the convenience
// functions Log...() below). Afterwards created loggers will not overwrie this. The convenience
// logger can be reset by using the method SetConvenienceLogger.
func NewLoggerFromFile(logfile io.Writer, level LogLevel) *Logger {
	l := &Logger{}
	l.internallogger = log.New(logfile, "LOG: ", loggerflags)
	l.ActiveLoglevel = level
	if conveninceLogger == nil {
		conveninceLogger = l
	}
	return l
}

// NewLogger creates a new Logger. It take a string file name as output file
// and a LogLevel to filter the messages that are wanted.
// The logger will use io.StdOut if the log filename string parameter is "STDOUT"
func NewLogger(logfilename string, level LogLevel) *Logger {
	var lfilename string
	var logfile io.Writer
	if logfilename == "" || strings.ToUpper(logfilename) == "STDOUT" {
		lfilename = "STDOUT"
		logfile = os.Stdout
	} else {
		lfilename = logfilename
		logfile, _ = os.OpenFile(lfilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	}
	return NewLoggerFromFile(logfile, level)
}

func (l *Logger) writelog(level LogLevel, format string, args ...interface{}) {
	if l.ActiveLoglevel >= level {
		l.internallogger.SetPrefix(strings.ToUpper(fmt.Sprintf("%5s: ", level)))
		l.internallogger.Output(3, fmt.Sprintf(format, args...))
	}
}

// Info works just as fmt.Printf, but prints into the loggers stream.
// The message is only printed if ActiveLogLevel is set higher or equal to 'Info'
func (l *Logger) Info(format string, args ...interface{}) {
	l.writelog(Info, format, args...)
}

// Debug works just as fmt.Printf, but prints into the loggers stream.
// The message is only printed if ActiveLogLevel is set higher or equal to 'Debug'
func (l *Logger) Debug(format string, args ...interface{}) {
	l.writelog(Debug, format, args...)
}

// Warn works just as fmt.Printf, but prints into the loggers stream.
// The message is only printed if ActiveLogLevel is set higher or equal to 'Warn'
func (l *Logger) Warn(format string, args ...interface{}) {
	l.writelog(Warn, format, args...)
}

// Error works just as fmt.Printf, but prints into the loggers stream.
// The message is only printed if ActiveLogLevel is set higher or equal to 'Error'
func (l *Logger) Error(format string, args ...interface{}) {
	l.writelog(Error, format, args...)
}

// Fatal works just as fmt.Printf, but prints into the loggers stream.
// The message is only printed if ActiveLogLevel is set hogher or equal to 'Fatal'
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.writelog(Fatal, format, args...)
}

// -----------------------------

// SetConvenienceLogger sets a logger as a singleton object. The LogInfo etc.
// functions use this singleton to offer logging function without an object context.
func (l *Logger) SetConvenienceLogger() {
	conveninceLogger = l
}

func outputToStandardLogger(level LogLevel, format string, args ...interface{}) {
	p := log.Prefix()
	f := log.Flags()
	log.SetFlags(loggerflags)
	log.SetPrefix(strings.ToUpper(fmt.Sprintf("%5s: ", level)))
	log.Output(3, fmt.Sprintf(format, args...))
	log.SetPrefix(p)
	log.SetFlags(f)
}

// LogInfo works just as fmt.Printf, but prints into the Convenience loggers stream, as set with
// SetConvenienceLogger(). It usses the standard logger (package log) if te Convenience logger is unset.
// The message is only printed if ActiveLogLevel is set higher or equal to 'Info'
func LogInfo(format string, args ...interface{}) {
	if conveninceLogger != nil {
		conveninceLogger.writelog(Info, format, args...)
	} else {
		outputToStandardLogger(Info, format, args...)
	}
}

// LogDebug works just as fmt.Printf, but prints into the Convenience loggers stream, as set with
// SetConvenienceLogger(). It usses the standard logger (package log) if te Convenience logger is unset.
// The message is only printed if ActiveLogLevel is set higher or equal to 'Debug'
func LogDebug(format string, args ...interface{}) {
	if conveninceLogger != nil {
		conveninceLogger.writelog(Debug, format, args...)
	} else {
		outputToStandardLogger(Debug, format, args...)
	}
}

// LogWarn works just as fmt.Printf, but prints into the Convenience loggers stream, as set with
// SetConvenienceLogger(). It usses the standard logger (package log) if te Convenience logger is unset.
// The message is only printed if ActiveLogLevel is set higher or equal to 'Warn'
func LogWarn(format string, args ...interface{}) {
	if conveninceLogger != nil {
		conveninceLogger.writelog(Warn, format, args...)
	} else {
		outputToStandardLogger(Warn, format, args...)
	}
}

// LogError works just as fmt.Printf, but prints into the Convenience loggers stream, as set with
// SetConvenienceLogger(). It usses the standard logger (package log) if te Convenience logger is unset.
// The message is only printed if ActiveLogLevel is set higher or equal to 'Error'
func LogError(format string, args ...interface{}) {
	if conveninceLogger != nil {
		conveninceLogger.writelog(Error, format, args...)
	} else {
		outputToStandardLogger(Error, format, args...)
	}
}

// LogFatal works just as fmt.Printf, but prints into the Convenience loggers stream, as set with
// SetConvenienceLogger(). It usses the standard logger (package log) if te Convenience logger is unset.
// The message is only printed if ActiveLogLevel is set hogher or equal to 'Fatal'
func LogFatal(format string, args ...interface{}) {
	if conveninceLogger != nil {
		conveninceLogger.writelog(Fatal, format, args...)
	} else {
		outputToStandardLogger(Fatal, format, args...)
	}
}
