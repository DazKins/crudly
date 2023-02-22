package middleware

import (
	"fmt"
	"io"
	"net/http"
)

type logger struct {
	writer io.Writer
}

func NewLogger(writer io.Writer) logger {
	return logger{
		writer,
	}
}

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

func (p logger) Attach(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ww := wrapWriter(w)

		h(ww, r)

		log := fmt.Sprintf(
			"[%s %s] %d\nResponse Body: %s\n",
			r.Method,
			r.URL,
			ww.loggerDetails.status,
			string(ww.loggerDetails.body),
		)

		if ww.loggerDetails.err != nil {
			log += fmt.Sprintf("Error: %s\n", ww.loggerDetails.err)
		}

		log += "\n"

		p.writer.Write([]byte(log))
	}
}

func AttachError(w http.ResponseWriter, err error) {
	wrapped, ok := w.(wrappedWriter)

	if !ok {
		return
	}

	wrapped.loggerDetails.err = err
}
