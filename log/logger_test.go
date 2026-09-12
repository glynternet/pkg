package log_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/glynternet/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaller(t *testing.T) {
	var bs bytes.Buffer
	require.NoError(t, log.NewLogger(&bs).Log())
	var capture map[string]interface{}
	require.NoError(t, json.Unmarshal(bs.Bytes(), &capture))
	assert.Equal(t, "logger_test.go:18", capture["caller"])
}

func TestLevels(t *testing.T) {
	t.Run("should inject info when no level applied", func(t *testing.T) {
		var bs bytes.Buffer
		require.NoError(t, log.NewLogger(&bs).Log())
		var capture map[string]interface{}
		require.NoError(t, json.Unmarshal(bs.Bytes(), &capture))
		assert.Equal(t, "info", capture["level"])
	})

	t.Run("should have applied level when applied", func(t *testing.T) {
		for _, tc := range []struct {
			logFn    func(logger log.Logger, kvs ...log.KV) error
			expected string
		}{{
			logFn:    log.Debug,
			expected: "debug",
		}, {
			logFn:    log.Info,
			expected: "info",
		}, {
			logFn:    log.Warn,
			expected: "warn",
		}, {
			logFn:    log.Error,
			expected: "error",
		}} {
			t.Run(tc.expected, func(t *testing.T) {
				var bs bytes.Buffer
				require.NoError(t, tc.logFn(log.NewLogger(&bs)))
				var capture map[string]interface{}
				require.NoError(t, json.Unmarshal(bs.Bytes(), &capture))
				t.Log(capture)
				assert.Equal(t, tc.expected, capture["level"])
			})
		}
	})
}

func TestNonTerminalLoggerUsesJSONWithoutANSI(t *testing.T) {
	var bs bytes.Buffer
	logger := log.NewLoggerWithOptions(&bs, log.Options{
		IsTerminal: func(io.Writer) bool { return false },
	})

	require.NoError(t, log.Error(logger,
		log.Message("request \"failed\"\nretrying"),
		log.ErrorMessage(errors.New("connection \"closed\""))))

	assert.NotContains(t, bs.String(), "\x1b[")
	var capture map[string]interface{}
	require.NoError(t, json.Unmarshal(bs.Bytes(), &capture))
	assert.Equal(t, "error", capture["level"])
	assert.Equal(t, "request \"failed\"\nretrying", capture["message"])
	assert.Equal(t, "connection \"closed\"", capture["error"])
	assert.NotEmpty(t, capture["ts"])
	assert.NotEmpty(t, capture["caller"])
}

func TestTerminalLoggerUsesColouredLogfmt(t *testing.T) {
	for _, tc := range []struct {
		name   string
		logFn  func(log.Logger, ...log.KV) error
		colour string
		level  string
	}{
		{name: "debug", logFn: log.Debug, colour: "\x1b[36m", level: "debug"},
		{name: "info", logFn: log.Info, colour: "\x1b[32m", level: "info"},
		{name: "warn", logFn: log.Warn, colour: "\x1b[33m", level: "warn"},
		{name: "error", logFn: log.Error, colour: "\x1b[31m", level: "error"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var bs bytes.Buffer
			logger := log.NewLoggerWithOptions(&bs, log.Options{
				IsTerminal: func(io.Writer) bool { return true },
			})

			require.NoError(t, tc.logFn(logger,
				log.Message("request \"failed\"\nretrying"),
				log.ErrorMessage(errors.New("connection \"closed\""))))

			output := bs.String()
			assert.Contains(t, output, "level="+tc.colour+tc.level+"\x1b[0m")
			assert.Equal(t, 1, strings.Count(output, "level="))
			assert.Equal(t, 1, strings.Count(output, "\x1b[0m"))
			assert.Contains(t, output, "message=\"request \\\"failed\\\"\\nretrying\"")
			assert.Contains(t, output, "error=\"connection \\\"closed\\\"\"")
			assert.Contains(t, output, "ts=")
			assert.Contains(t, output, "caller=")
			assert.True(t, bytes.HasSuffix(bs.Bytes(), []byte("\n")))
		})
	}
}
