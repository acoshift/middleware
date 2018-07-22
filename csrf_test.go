package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acoshift/middleware"
	"github.com/stretchr/testify/assert"
)

func _h(w http.ResponseWriter, r *http.Request) {}

var h = http.HandlerFunc(_h)

func TestCSRF(t *testing.T) {
	t.Run("Base", func(t *testing.T) {
		m := middleware.CSRF(middleware.CSRFConfig{
			Origins: []string{"http://example.com"},
		})

		assert.Equal(t, 200, testOrigin(m, "http://example.com"))
		assert.Equal(t, 403, testOrigin(m, "https://example.com"))
		assert.Equal(t, 403, testOrigin(m, "example.com"))
		assert.Equal(t, 403, testOrigin(m, "hacker.com"))
		assert.Equal(t, 403, testOrigin(m, "http://hacker.com"))
		assert.Equal(t, 403, testOrigin(m, "https://hacker.com"))

		assert.Equal(t, 200, testReferer(m, "http://example.com/"))
		assert.Equal(t, 403, testReferer(m, "https://example.com/"))
		assert.Equal(t, 403, testReferer(m, "example.com/"))
		assert.Equal(t, 403, testReferer(m, "hacker.com/"))
		assert.Equal(t, 403, testReferer(m, "http://hacker.com/"))
		assert.Equal(t, 403, testReferer(m, "https://hacker.com/"))
	})

	t.Run("IgnoreProto", func(t *testing.T) {
		m := middleware.CSRF(middleware.CSRFConfig{
			Origins:     []string{"http://example.com"},
			IgnoreProto: true,
		})

		assert.Equal(t, 200, testOrigin(m, "http://example.com"))
		assert.Equal(t, 200, testOrigin(m, "https://example.com"))
		assert.Equal(t, 403, testOrigin(m, "example.com"))
		assert.Equal(t, 403, testOrigin(m, "hacker.com"))
		assert.Equal(t, 403, testOrigin(m, "http://hacker.com"))
		assert.Equal(t, 403, testOrigin(m, "https://hacker.com"))

		assert.Equal(t, 200, testReferer(m, "http://example.com/"))
		assert.Equal(t, 200, testReferer(m, "https://example.com/"))
		assert.Equal(t, 403, testReferer(m, "example.com/"))
		assert.Equal(t, 403, testReferer(m, "hacker.com/"))
		assert.Equal(t, 403, testReferer(m, "http://hacker.com/"))
		assert.Equal(t, 403, testReferer(m, "https://hacker.com/"))
	})

	t.Run("IgnoreProto2", func(t *testing.T) {
		m := middleware.CSRF(middleware.CSRFConfig{
			Origins:     []string{"example.com"},
			IgnoreProto: true,
		})

		assert.Equal(t, 200, testOrigin(m, "http://example.com"))
		assert.Equal(t, 200, testOrigin(m, "https://example.com"))
		assert.Equal(t, 403, testOrigin(m, "example.com"))
		assert.Equal(t, 403, testOrigin(m, "hacker.com"))
		assert.Equal(t, 403, testOrigin(m, "http://hacker.com"))
		assert.Equal(t, 403, testOrigin(m, "https://hacker.com"))

		assert.Equal(t, 200, testReferer(m, "http://example.com/"))
		assert.Equal(t, 200, testReferer(m, "https://example.com/"))
		assert.Equal(t, 403, testReferer(m, "example.com/"))
		assert.Equal(t, 403, testReferer(m, "hacker.com/"))
		assert.Equal(t, 403, testReferer(m, "http://hacker.com/"))
		assert.Equal(t, 403, testReferer(m, "https://hacker.com/"))
	})
}

func testOrigin(m middleware.Middleware, origin string) int {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Origin", origin)
	m(h).ServeHTTP(w, r)
	return w.Code
}

func testReferer(m middleware.Middleware, referer string) int {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Referer", referer)
	m(h).ServeHTTP(w, r)
	return w.Code
}
