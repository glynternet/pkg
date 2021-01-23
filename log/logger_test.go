package log_test

import (
	"bytes"
	"encoding/json"
	"github.com/glynternet/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCaller(t *testing.T) {
	var bs bytes.Buffer
	require.NoError(t, log.NewLogger(&bs).Log())
	var capture map[string]interface{}
	require.NoError(t, json.Unmarshal(bs.Bytes(), &capture))
	assert.Equal(t, "logger_test.go:14", capture["caller"])
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
