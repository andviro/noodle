package middleware

import (
	"bufio"
	"github.com/andviro/noodle"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var writers sync.Pool

// logWriter mimics http.ResponseWriter functionality while storing
// HTTP status code for later logging
type logWriter struct {
	code          int
	headerWritten bool
	http.ResponseWriter
}

func (l *logWriter) WriteHeader(code int) {
	l.ResponseWriter.WriteHeader(code)
	if !l.headerWritten {
		l.code = code
		l.headerWritten = true
	}
}

func (l *logWriter) Write(buf []byte) (int, error) {
	l.headerWritten = true
	return l.ResponseWriter.Write(buf)
}

func (l *logWriter) Code() int {
	if l.code == 0 {
		return 200
	}
	return l.code
}

// provide other typical ResponseWriter methods
func (l *logWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return l.ResponseWriter.(http.Hijacker).Hijack()
}

func (l *logWriter) CloseNotify() <-chan bool {
	return l.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (l *logWriter) Flush() {
	l.ResponseWriter.(http.Flusher).Flush()
}

func init() {
	writers.New = func() interface{} {
		return &logWriter{}
	}
}

// Logger is a middleware that logs requests, along with
// request URI, HTTP status code, handler return value and request timing
func Logger(next noodle.Handler) noodle.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		lw := writers.Get().(*logWriter)
		lw.ResponseWriter = w
		lw.code = 0
		lw.headerWritten = false
		defer writers.Put(lw)

		url := r.URL.String() // further calls may modify request URL
		start := time.Now()
		next(lw, r)
		end := time.Now()
		remoteAddr := GetRealIP(r) // try to get client address from middleware
		if remoteAddr == "" {
			remoteAddr = r.RemoteAddr
		}
		log.Printf("%s %s (%d) from %s [%s]", r.Method, url, lw.Code(), remoteAddr, end.Sub(start))
		return
	}
}
