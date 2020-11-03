package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	ghttp "github.com/glynternet/pkg/http"
	"github.com/glynternet/pkg/log"
	"github.com/stretchr/testify/assert"
)

type mockRequestAuthoriser struct {
	err     error
	request *http.Request
}

func (ra *mockRequestAuthoriser) Authorise(r *http.Request) error {
	ra.request = r
	return ra.err
}

type mockHandler struct {
	request *http.Request
	writer  http.ResponseWriter
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.writer = w
	h.request = r
}

func TestWithAuthoriser(t *testing.T) {
	t.Run("unauthorised", func(t *testing.T) {
		var next mockHandler
		err := errors.New("auth error")
		ra := mockRequestAuthoriser{err: err}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://any/", nil)

		logger := mockLogger{}
		hf := ghttp.WithAuthoriser(&logger, &ra, &next)
		hf(w, r)
		assert.Equal(t, r, ra.request)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, logger, 2)
		assert.Equal(t, log.Message("Unauthorised request"), logger[0])
		assert.Equal(t, log.Error(err), logger[1])
		assert.Nil(t, next.writer)
		assert.Nil(t, next.request)
	})

	t.Run("authorised", func(t *testing.T) {
		logger := mockLogger{}
		var ra mockRequestAuthoriser
		next := mockHandler{}
		hf := ghttp.WithAuthoriser(&logger, &ra, &next)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://any/", nil)
		hf(w, r)
		assert.Len(t, logger, 0)
		assert.Equal(t, w, next.writer)
		assert.Equal(t, r, next.request)
	})
}

type mockLogger []log.KV

func (m *mockLogger) Log(kv ...log.KV) error {
	*m = append(*m, kv...)
	return nil
}
