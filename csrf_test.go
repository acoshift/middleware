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

		assert.Equal(t, 200, test(m, "", ""))

		assert.Equal(t, 200, test(m, "http://example.com", ""))
		assert.Equal(t, 200, test(m, "http://example.com", "http://example.com/"))
		assert.Equal(t, 403, test(m, "http://example.com", "http://example.com"))
		assert.Equal(t, 403, test(m, "http://example.com", "https://example.com"))
		assert.Equal(t, 403, test(m, "http://example.com", "https://example.com/"))
		assert.Equal(t, 403, test(m, "https://example.com", ""))
		assert.Equal(t, 403, test(m, "https://example.com", "http://example.com/"))
		assert.Equal(t, 403, test(m, "https://example.com", "http://example.com"))
		assert.Equal(t, 403, test(m, "example.com", ""))
		assert.Equal(t, 403, test(m, "hacker.com", ""))
		assert.Equal(t, 403, test(m, "http://hacker.com", ""))
		assert.Equal(t, 403, test(m, "http://hacker.com", "http://hacker.com/"))
		assert.Equal(t, 403, test(m, "http://hacker.com", "http://example.com/"))
		assert.Equal(t, 403, test(m, "https://hacker.com", ""))
		assert.Equal(t, 403, test(m, "https://hacker.com", "https://hacker.com/"))
		assert.Equal(t, 403, test(m, "https://hacker.com", "http://example.com/"))

		assert.Equal(t, 200, test(m, "", "http://example.com/"))
		assert.Equal(t, 403, test(m, "", "https://example.com/"))
		assert.Equal(t, 403, test(m, "", "example.com/"))
		assert.Equal(t, 403, test(m, "", "hacker.com/"))
		assert.Equal(t, 403, test(m, "", "http://hacker.com/"))
		assert.Equal(t, 403, test(m, "", "https://hacker.com/"))
	})

	t.Run("Force", func(t *testing.T) {
		m := middleware.CSRF(middleware.CSRFConfig{
			Origins: []string{"http://example.com"},
			Force:   true,
		})

		assert.Equal(t, 403, test(m, "", ""))

		assert.Equal(t, 403, test(m, "http://example.com", ""))
		assert.Equal(t, 403, test(m, "http://example.com", "http://example.com"))
		assert.Equal(t, 200, test(m, "http://example.com", "http://example.com/"))
		assert.Equal(t, 200, test(m, "http://example.com", "http://example.com/page1"))
		assert.Equal(t, 200, test(m, "http://example.com", "http://example.com/page1/page2"))
		assert.Equal(t, 403, test(m, "http://example.com", "https://example.com/"))
		assert.Equal(t, 403, test(m, "http://example.com", "https://example.com/page1"))
		assert.Equal(t, 403, test(m, "http://example.com", "https://example.com/page1/page2"))
		assert.Equal(t, 403, test(m, "https://example.com", ""))
		assert.Equal(t, 403, test(m, "https://example.com", ""))
		assert.Equal(t, 403, test(m, "example.com", ""))
		assert.Equal(t, 403, test(m, "hacker.com", ""))
		assert.Equal(t, 403, test(m, "http://hacker.com", ""))
		assert.Equal(t, 403, test(m, "http://hacker.com", "http://hacker.com"))
		assert.Equal(t, 403, test(m, "http://hacker.com", "http://example.com/"))
		assert.Equal(t, 403, test(m, "https://hacker.com", ""))

		assert.Equal(t, 403, test(m, "", "http://example.com/"))
		assert.Equal(t, 403, test(m, "", "https://example.com/"))
		assert.Equal(t, 403, test(m, "", "example.com/"))
		assert.Equal(t, 403, test(m, "", "hacker.com/"))
		assert.Equal(t, 403, test(m, "", "http://hacker.com/"))
		assert.Equal(t, 403, test(m, "", "https://hacker.com/"))
	})

	t.Run("IgnoreProto", func(t *testing.T) {
		m := middleware.CSRF(middleware.CSRFConfig{
			Origins:     []string{"http://example.com"},
			IgnoreProto: true,
		})

		assert.Equal(t, 200, test(m, "http://example.com", ""))
		assert.Equal(t, 200, test(m, "https://example.com", ""))
		assert.Equal(t, 403, test(m, "example.com", ""))
		assert.Equal(t, 403, test(m, "hacker.com", ""))
		assert.Equal(t, 403, test(m, "http://hacker.com", ""))
		assert.Equal(t, 403, test(m, "https://hacker.com", ""))

		assert.Equal(t, 200, test(m, "", "http://example.com/"))
		assert.Equal(t, 200, test(m, "", "https://example.com/"))
		assert.Equal(t, 403, test(m, "", "example.com/"))
		assert.Equal(t, 403, test(m, "", "hacker.com/"))
		assert.Equal(t, 403, test(m, "", "http://hacker.com/"))
		assert.Equal(t, 403, test(m, "", "https://hacker.com/"))
	})

	t.Run("IgnoreProto2", func(t *testing.T) {
		m := middleware.CSRF(middleware.CSRFConfig{
			Origins:     []string{"example.com"},
			IgnoreProto: true,
		})

		assert.Equal(t, 200, test(m, "http://example.com", ""))
		assert.Equal(t, 200, test(m, "https://example.com", ""))
		assert.Equal(t, 403, test(m, "example.com", ""))
		assert.Equal(t, 403, test(m, "hacker.com", ""))
		assert.Equal(t, 403, test(m, "http://hacker.com", ""))
		assert.Equal(t, 403, test(m, "https://hacker.com", ""))

		assert.Equal(t, 200, test(m, "", "http://example.com/"))
		assert.Equal(t, 200, test(m, "", "https://example.com/"))
		assert.Equal(t, 403, test(m, "", "example.com/"))
		assert.Equal(t, 403, test(m, "", "hacker.com/"))
		assert.Equal(t, 403, test(m, "", "http://hacker.com/"))
		assert.Equal(t, 403, test(m, "", "https://hacker.com/"))
	})

}

func test(m middleware.Middleware, origin, referer string) int {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Origin", origin)
	r.Header.Set("Referer", referer)
	r.Header.Set("X-Forwarded-Proto", "https")
	m(h).ServeHTTP(w, r)
	return w.Code
}
