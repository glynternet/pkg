package log

import (
	"bytes"
	"io"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/mattn/go-isatty"
)

const levelKey = "level"

// TerminalDetector reports whether a writer is connected to a terminal.
// It can be supplied in Options to make logger selection deterministic.
type TerminalDetector func(io.Writer) bool

// Options configures a Logger.
type Options struct {
	// IsTerminal overrides the default terminal detection when non-nil.
	IsTerminal TerminalDetector
}

// NewLogger creates a logger which writes coloured, human-readable logfmt to
// terminals and JSON to all other writers.
func NewLogger(w io.Writer) Logger {
	return NewLoggerWithOptions(w, Options{})
}

// NewLoggerWithOptions creates a logger with the supplied options. JSON is
// used unless IsTerminal reports that w is a terminal.
func NewLoggerWithOptions(w io.Writer, options Options) Logger {
	detector := options.IsTerminal
	if detector == nil {
		detector = isTerminal
	}

	var output log.Logger = log.NewJSONLogger(w)
	if detector(w) {
		output = log.NewLogfmtLogger(colorWriter{Writer: w})
	}

	return logger{
		Logger: log.With(
			level.NewFilter(
				level.NewInjector(output, level.InfoValue()),
				level.AllowDebug()),
			"ts", log.DefaultTimestampUTC,
			"caller", log.Caller(4)),
	}
}

type Logger interface {
	Log(...KV) error
}

type logger struct {
	log.Logger
}

func (l logger) Log(kvs ...KV) error {
	return l.Logger.Log(kvsAsParams(kvs)...)
}

func Debug(logger Logger, kvs ...KV) error {
	return logger.Log(append(
		kvs,
		KV{K: levelKey, V: "debug"},
	)...)
}

func Info(logger Logger, kvs ...KV) error {
	return logger.Log(append(
		kvs,
		KV{K: levelKey, V: "info"},
	)...)
}

func Warn(logger Logger, kvs ...KV) error {
	return logger.Log(append(
		kvs,
		KV{K: levelKey, V: "warn"},
	)...)
}

func Error(logger Logger, kvs ...KV) error {
	return logger.Log(append(
		kvs,
		KV{K: levelKey, V: "error"},
	)...)
}

type KV struct {
	K string
	V interface{}
}

type kvs []KV

func kvsAsParams(kvs kvs) []interface{} {
	ps := make([]interface{}, 2*len(kvs))
	for i, kv := range kvs {
		ps[i*2], ps[i*2+1] = kv.K, normaliseLevel(kv)
	}
	return ps
}

// normaliseLevel converts the level strings accepted by this package into
// go-kit level values. Besides allowing the level filter to see them, this
// prevents a default info level from being injected alongside an explicit one.
func normaliseLevel(kv KV) interface{} {
	if kv.K != levelKey {
		return kv.V
	}

	levelName, ok := kv.V.(string)
	if !ok {
		return kv.V
	}

	switch levelName {
	case "debug":
		return level.DebugValue()
	case "info":
		return level.InfoValue()
	case "warn":
		return level.WarnValue()
	case "error":
		return level.ErrorValue()
	default:
		return kv.V
	}
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	return ok && (isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd()))
}

// colorWriter adds ANSI colour to level values after go-kit has encoded the
// rest of the record. Applying colour at this point keeps logfmt's escaping
// intact and leaves all non-level fields unchanged.
type colorWriter struct {
	io.Writer
}

func (w colorWriter) Write(p []byte) (int, error) {
	colored := colorizeLevels(p)
	n, err := w.Writer.Write(colored)
	if n == len(colored) {
		return len(p), err
	}
	if err == nil {
		err = io.ErrShortWrite
	}
	return 0, err
}

var levelColors = map[string]string{
	"debug": "\x1b[36m",
	"info":  "\x1b[32m",
	"warn":  "\x1b[33m",
	"error": "\x1b[31m",
}

const colorReset = "\x1b[0m"

func colorizeLevels(record []byte) []byte {
	var output bytes.Buffer
	output.Grow(len(record) + 16)

	fieldStart := true
	inQuote := false
	escaped := false
	for i := 0; i < len(record); {
		if fieldStart {
			for levelName, color := range levelColors {
				prefix := "level=" + levelName
				if bytes.HasPrefix(record[i:], []byte(prefix)) && isFieldEnd(record, i+len(prefix)) {
					output.WriteString("level=")
					output.WriteString(color)
					output.WriteString(levelName)
					output.WriteString(colorReset)
					i += len(prefix)
					fieldStart = false
					break
				}
			}
			if !fieldStart {
				continue
			}
		}

		c := record[i]
		output.WriteByte(c)
		i++

		if inQuote {
			if escaped {
				escaped = false
			} else if c == '\\' {
				escaped = true
			} else if c == '"' {
				inQuote = false
			}
			continue
		}

		switch c {
		case '"':
			inQuote = true
			fieldStart = false
		case ' ', '\n':
			fieldStart = true
		default:
			fieldStart = false
		}
	}
	return output.Bytes()
}

func isFieldEnd(record []byte, index int) bool {
	return index == len(record) || record[index] == ' ' || record[index] == '\n'
}

func Message(msg string) KV {
	return KV{
		K: "message",
		V: msg,
	}
}

func ErrorMessage(err error) KV {
	return KV{
		K: "error",
		V: err,
	}
}
