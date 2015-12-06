package middleware

import (
	"github.com/andviro/noodle"
	"golang.org/x/net/context"
	"log"
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

func init() {
	writers.New = func() interface{} {
		return &logWriter{}
	}
}

// Logger is a middleware that logs requests, along with
// request URI, HTTP status code, handler return value and request timing
func Logger(next noodle.Handler) noodle.Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) (err error) {
		lw := writers.Get().(*logWriter)
		lw.ResponseWriter = w
		lw.code = 0
		lw.headerWritten = false
		defer writers.Put(lw)

		start := time.Now()
		err = next(c, lw, r)
		end := time.Now()
		log.Printf("%s %s (%d) from %s [%s] error = %v", r.Method, r.URL.String(), lw.Code(), r.RemoteAddr, end.Sub(start), err)
		return
	}
}
