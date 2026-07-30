package web

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-NewWithClaims-Hmac
func newSessionToken(secret string) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = "v1"
	secretBytes := []byte(secret)
	return token.SignedString(secretBytes)
}

func verifySessionToken(secret string, tokenString string) bool {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok || kid != "v1" {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"})) // pin the algorithm

	if err != nil || !token.Valid {
		return false // covers bad signature and expired token
	}
	return true
}

func (app *application) isAuthenticated(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil { // http.ErrNoCookie if absent
		return false
	}
	app.logger.Debug("session cookie found")
	return verifySessionToken(app.jwtSecret, c.Value)
}

// https://www.alexedwards.net/blog/working-with-cookies-in-go
var cookieTemplate = http.Cookie{
	Name:     "session",
	Value:    "",
	Path:     "/", // must match across every route that sets/clears this cookie, or login/logout end up scoped to different cookies
	HttpOnly: true, // JS (document.cookie) can't read it — mitigates XSS token theft
	// Secure:   true,                  // only sent over HTTPS — assumes you're behind TLS in prod
	SameSite: http.SameSiteLaxMode, // mitigates CSRF; Lax still allows top-level GET navigation (if landing on site from an external link)
}
