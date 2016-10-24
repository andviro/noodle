package middleware

import (
	"context"
	"errors"
	"github.com/andviro/noodle"
	"net/http"
	"net/url"
)

var UnauthorizedRequest = errors.New("Unauthorized request")

// HTTPAuth is a middleware factory function that accepts the authentication realm
// and function for username and password verification. Resulting middleware injects
// username into request context if authentication successful.
func HTTPAuth(realm string, authFunc func(username, password string) bool) noodle.Middleware {
	return func(next noodle.Handler) noodle.Handler {
		return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			username, password, ok := r.BasicAuth()
			if !ok || !authFunc(username, password) {
				w.Header().Set("WWW-Authenticate", "Basic realm="+url.QueryEscape(realm))
				w.WriteHeader(http.StatusUnauthorized)
				// Provide error for logging middleware then abort chain
				return UnauthorizedRequest
			}
			// Inject user name into request context
			return next(context.WithValue(c, userKey, username), w, r)
		}
	}
}

// GetUser extract authentication information from context
func GetUser(c context.Context) string {
	res, _ := c.Value(userKey).(string)
	return res
}
