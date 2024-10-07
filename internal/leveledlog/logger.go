package leveledlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"example_project/internal/config"

	"github.com/fatih/color"
	"github.com/gofrs/uuid"
)

var (
	redBold    = color.New(color.FgHiRed, color.Bold).SprintFunc()
	yellowBold = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	cyanBold   = color.New(color.FgHiCyan, color.Bold).SprintFunc()
	whiteBold  = color.New(color.FgHiWhite, color.Bold).SprintFunc()
	red        = color.New(color.FgHiRed).SprintFunc()
	yellow     = color.New(color.FgHiYellow).SprintFunc()
	cyan       = color.New(color.FgHiCyan).SprintFunc()
	white      = color.New(color.FgHiWhite).SprintFunc()
	idColor    = color.New(color.FgYellow).SprintFunc()
)

type Level int8

const (
	LevelAll Level = iota
	LevelDebug
	LevelInfo
	LevelHttp
	LevelWarning
	LevelError
	LevelFatal
	LevelOff
)

var Logger = NewLogger(os.Stdout, GetLevel(config.ConfigStruct.LogLevel), true)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO "
	case LevelDebug:
		return "DEBUG"
	case LevelWarning:
		return "WARN "
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelHttp:
		return "HTTP "
	default:
		return ""
	}
}

func GetLevel(level string) Level {
	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARNING":
		return LevelWarning
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	default:
		return LevelAll
	}
}

type LogStruct struct {
	out          io.Writer
	minLevel     Level
	useJSON      bool
	colorize     bool
	absolutePath string
	mu           *sync.Mutex
	id           string
	withID       bool
}

func (l *LogStruct) WithID(id string) *LogStruct {
	return &LogStruct{
		out:          l.out,
		minLevel:     l.minLevel,
		useJSON:      l.useJSON,
		colorize:     l.colorize,
		absolutePath: l.absolutePath,
		mu:           l.mu,
		id:           id,
		withID:       true,
	}
}

func (l *LogStruct) SetFilePath(path string) {
	l.absolutePath = path
}

func NewLogger(out io.Writer, minLevel Level, colorize bool) *LogStruct {
	id, _ := uuid.NewV4()
	return &LogStruct{
		out:      out,
		minLevel: minLevel,
		useJSON:  false,
		colorize: colorize,
		id:       fmt.Sprintf("%v", id),
		mu:       &sync.Mutex{},
	}
}

func NewJSONLogger(out io.Writer, minLevel Level) *LogStruct {
	return &LogStruct{
		out:      out,
		minLevel: minLevel,
		useJSON:  true,
		mu:       &sync.Mutex{},
	}
}

func (l *LogStruct) Debug(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.print(LevelDebug, message)
}

func (l *LogStruct) Info(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.print(LevelInfo, message)
}

func (l *LogStruct) Warning(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.print(LevelWarning, message)
}

func (l *LogStruct) Error(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.print(LevelError, message)
}

func (l *LogStruct) Fatal(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.print(LevelFatal, message)
	os.Exit(1)
}

// method for http request log
func (l *LogStruct) Print(v ...interface{}) {
	var line string
	for _, vv := range v {
		line += vv.(string)
	}
	l.print(LevelHttp, line)
}

func (l *LogStruct) KafkaLog(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.out, message)
}

func (l *LogStruct) print(level Level, message string) {
	if level < l.minLevel {
		return
	}
	var line string
	if l.useJSON {
		line = jsonLine(level, message)
	} else {
		line = textLine(level, l.id, message, l.colorize)
	}
	if l.minLevel < LevelError {
		line = strings.Replace(line, "/usr/src/app/", "/", 1)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.out, line)
}

func textLine(level Level, id string, message string, colorize bool) string {
	var fileLoc, line string
	_, fileName, lineNumber, ok := runtime.Caller(3)
	if ok {
		fileLoc = fmt.Sprintf("%v:%v", fileName, lineNumber)
	}

	if level == LevelHttp {
		line = fmt.Sprintf(" | [%v] %v ", time.Now().Format("2006-01-02 15:04:05"), message)
	} else {
		line = fmt.Sprintf(" | [%v] [%v] %v --> %v ", time.Now().Format("2006-01-02 15:04:05"), idColor(id), fileLoc, message)
	}

	if colorize {
		switch level {
		case LevelError, LevelFatal:
			line = fmt.Sprintf("%-6s", redBold(level)) + red(line)
		case LevelWarning:
			line = fmt.Sprintf("%-6s", yellowBold(level)) + yellow(line)
		case LevelDebug:
			line = fmt.Sprintf("%-6s", cyanBold(level)) + cyan(line)
		case LevelInfo:
			line = fmt.Sprintf("%-6s", whiteBold(level)) + white(line)
		case LevelHttp:
			line = fmt.Sprintf("%-6s", yellowBold(level)) + line
		}
	} else {
		line = fmt.Sprintf("%-6s%v", level, line)
	}

	// if level >= LevelError {
	// 	line += fmt.Sprintf("\n%s", string(debug.Stack()))
	// }

	return line
}

func jsonLine(level Level, message string) string {
	aux := struct {
		Level   string `json:"level"`
		Time    string `json:"time"`
		Message string `json:"message"`
		Trace   string `json:"trace,omitempty"`
	}{
		Level:   level.String(),
		Time:    time.Now().UTC().Format(time.RFC3339),
		Message: message,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		return fmt.Sprintf("%s: unable to marshal log message: %s", LevelError.String(), err.Error())
	}

	return string(line)
}
