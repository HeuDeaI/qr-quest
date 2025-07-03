package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"qr-quest/internal/tests"
)

func TestAdminLogin(t *testing.T) {
	router := tests.SetupTestRouter()

	t.Run("Valid credentials", func(t *testing.T) {
		form := url.Values{"password": {"admin"}} // Correct password
		req := httptest.NewRequest("POST", "/admin/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Admin login successful")
		assert.NotEmpty(t, w.Result().Cookies())
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		form := url.Values{"password": {"wrong_password"}}
		req := httptest.NewRequest("POST", "/admin/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid password")
	})
}
