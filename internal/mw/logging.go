package mw

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func LoggerHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		log.Printf("request: method: %v, url: %v,\nbody: %v\n", r.Method, r.URL.Path, string(reqBody))

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		log.Printf("response: method: %v, status: %v, body: %v\n", r.Method, rw.statusCode, rw.body.String())
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
