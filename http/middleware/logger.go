package middleware

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type loggerDetails struct {
	status int
	body   []byte
	err    error
}

type wrappedWriter struct {
	writer        http.ResponseWriter
	loggerDetails *loggerDetails
}

func wrapWriter(w http.ResponseWriter) wrappedWriter {
	return wrappedWriter{
		writer: w,

		loggerDetails: &loggerDetails{
			status: 200,
			body:   []byte(""),
		},
	}
}

func (w wrappedWriter) WriteHeader(statusCode int) {
	w.loggerDetails.status = statusCode

	w.writer.WriteHeader(statusCode)
}

func (w wrappedWriter) Header() http.Header {
	return w.writer.Header()
}

func (w wrappedWriter) Write(b []byte) (int, error) {
	w.loggerDetails.body = b

	return w.writer.Write(b)
}

func NewLogger(writer io.Writer) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := wrapWriter(w)

			t0 := time.Now()

			h.ServeHTTP(ww, r)

			t1 := time.Now()

			body := strings.TrimSuffix(string(ww.loggerDetails.body), "\n")

			requestTimeMs := t1.Sub(t0).Milliseconds()

			log := fmt.Sprintf(
				"[%s %s] %d\nRequest Time: %dms\nResponse Body: %s\n",
				r.Method,
				r.URL,
				ww.loggerDetails.status,
				requestTimeMs,
				body,
			)

			if ww.loggerDetails.err != nil {
				log += fmt.Sprintf("Error: %s\n", ww.loggerDetails.err)
			}

			log += "\n"

			writer.Write([]byte(log))
		})
	}
}

func AttachError(w http.ResponseWriter, err error) {
	wrapped, ok := w.(wrappedWriter)

	if !ok {
		return
	}

	wrapped.loggerDetails.err = err
}
