package middleware

import (
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/test"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", nil)
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	config := JWTAuthConfig{}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"

	// No signing key provided
	assert.Panics(t, func() {
		JWTAuthWithConfig(config)
	})

	// Unexpected signing method
	config.SigningKey = []byte("secret")
	config.SigningMethod = "RS256"
	h := JWTAuthWithConfig(config)(handler)
	he := h(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, he.Code)

	// Invalid key
	auth := bearer + " " + token
	req.Header().Set(echo.HeaderAuthorization, auth)
	config.SigningKey = []byte("invalid-key")
	h = JWTAuthWithConfig(config)(handler)
	he = h(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, he.Code)

	// Valid JWT
	h = JWTAuth([]byte("secret"))(handler)
	if assert.NoError(t, h(c)) {
		user := c.Get("user").(*jwt.Token)
		assert.Equal(t, user.Claims["name"], "John Doe")
	}

	// Invalid Authorization header
	req.Header().Set(echo.HeaderAuthorization, "invalid-auth")
	h = JWTAuth([]byte("secret"))(handler)
	he = h(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, he.Code)
}
