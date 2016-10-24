package middleware

import (
	"context"
	"github.com/andviro/noodle"
	"net/http"
	"strings"
)

var realIPKey int = 0

// clientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
// This is almost unmodified code from Gin framework and all credits and my deepest thanks go to Gin developers.
func clientIP(r *http.Request) string {
	clientIP := strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = r.Header.Get("X-Forwarded-For")
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	return strings.TrimSpace(r.RemoteAddr)
}

// RealIP is a middleware that injects client IP parsed from request headers into context
func RealIP(next noodle.Handler) noodle.Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) (err error) {
		return next(context.WithValue(c, realIPKey, clientIP(r)), w, r)
	}
}

// GetRealIP extracts real client IP from handler context
func GetRealIP(c context.Context) string {
	res, _ := c.Value(realIPKey).(string)
	return res
}
