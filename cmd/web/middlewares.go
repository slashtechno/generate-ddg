package web

import (
	"context"
	"net/http"
)

// https://www.alexedwards.net/blog/making-and-using-middleware

type contextKey string

const authedContextKey contextKey = "authed"

func (app *application) attachAuthStatus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// WithValue wraps r.Context() in a new, unexported context.valueCtx
		// holding {authedContextKey: bool}; the old context is untouched.
		// Value() checks this key first, else delegates to the wrapped parent.
		// Done()/Err()/Deadline() are promoted straight through from the parent.
		ctx := context.WithValue(r.Context(), authedContextKey, app.isAuthenticated(r))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			redirect(w, r, "/", http.StatusSeeOther)
			return // critical to stop the request from continuing to the next handler
		}
		next.ServeHTTP(w, r)
	})
}
