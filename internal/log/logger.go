package log

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ConsoleEncoderName ...
const ConsoleEncoderName = "custom_console"

var (
	ll Logger
)

// Logger wraps zap.Logger
type Logger struct {
	*zap.Logger
	S *zap.SugaredLogger
}

// PrintError prints all error with all meta data and line number.
// It's prefered to be used at top level function.
//
//     func DoSomething() (_err error) {
//         defer ll.PrintError("DoSomething", &_err)
//
func (logger Logger) PrintError(msg string, err *error) {
	if *err != nil {
		ll.S.Errorf("%v: %+v", msg, *err)
	}
}

// Short-hand functions for logging.
var (
	String     = zap.String
)

// DefaultConsoleEncoderConfig ...
var DefaultConsoleEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   ShortColorCallerEncoder,
}

func trimPath(c zapcore.EntryCaller) string {
	return c.TrimmedPath()
}

// ShortColorCallerEncoder encodes caller information with sort path filename and enable color.
func ShortColorCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	const gray, resetColor = "\x1b[90m", "\x1b[0m"
	callerStr := gray + "→ " + trimPath(caller) + ":" + strconv.Itoa(caller.Line) + resetColor
	enc.AppendString(callerStr)
}

func newWithName(name string, opts ...zap.Option) Logger {
	if name == "" {
		_, filename, _, _ := runtime.Caller(1)
		name = filepath.Dir(filename)
	}

	var enabler zap.AtomicLevel
	if e, ok := enablers[name]; ok {
		enabler = e
	} else {
		enabler = zap.NewAtomicLevel()
		enablers[name] = enabler
	}

	setLogLevelFromEnv(name, enabler)
	loggerConfig := zap.Config{
		Level:            enabler,
		Development:      false,
		Encoding:         ConsoleEncoderName,
		EncoderConfig:    DefaultConsoleEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	stacktraceLevel := zap.NewAtomicLevelAt(zapcore.PanicLevel)

	opts = append(opts, zap.AddStacktrace(stacktraceLevel))
	logger, err := loggerConfig.Build(opts...)
	if err != nil {
		panic(err)
	}
	return Logger{logger, logger.Sugar()}
}

// New returns new zap.Logger
func New(opts ...zap.Option) Logger {
	return newWithName("", opts...)
}

var (
	enablers = make(map[string]zap.AtomicLevel)
)

var envPatterns []*regexp.Regexp

func init() {
	err := zap.RegisterEncoder(ConsoleEncoderName, func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return NewConsoleEncoder(cfg), nil
	})
	if err != nil {
		panic(err)
	}

	ll = New()

	envLog := os.Getenv("LOG_LEVEL")
	if envLog == "" {
		return
	}

	var lv zapcore.Level
	err = lv.UnmarshalText([]byte(envLog))
	if err != nil {
		panic(err)
	}

	for _, enabler := range enablers {
		enabler.SetLevel(lv)
	}

	var errPattern string
	envPatterns, errPattern = initPatterns(envLog)
	if errPattern != "" {
		ll.Fatal("Unable to parse LOG_LEVEL. Please set it to a proper value.", String("invalid", errPattern))
	}

	ll.Info("Enable debug log", String("LOG_LEVEL", envLog))
}

func initPatterns(envLog string) ([]*regexp.Regexp, string) {
	patterns := strings.Split(envLog, ",")
	result := make([]*regexp.Regexp, len(patterns))
	for i, p := range patterns {
		r, err := parsePattern(p)
		if err != nil {
			return nil, p
		}

		result[i] = r
	}
	return result, ""
}

func parsePattern(p string) (*regexp.Regexp, error) {
	p = strings.Replace(strings.Trim(p, " "), "*", ".*", -1)
	return regexp.Compile(p)
}

func setLogLevelFromEnv(name string, enabler zap.AtomicLevel) {
	for _, r := range envPatterns {
		if r.MatchString(name) {
			enabler.SetLevel(zap.DebugLevel)
		}
	}
}
