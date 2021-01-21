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
	var capture struct {
		Caller string
	}
	require.NoError(t, json.Unmarshal(bs.Bytes(), &capture))
	assert.Equal(t, "logger_test.go:14", capture.Caller)
}
