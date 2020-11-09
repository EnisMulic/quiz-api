package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// GzipMiddleware zips files
func GzipMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// create a gziped response
			wrw := NewWrappedResponseWriter(rw)
			wrw.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(wrw, r)
			defer wrw.Flush()

			return
		}

		// handle normal
		next.ServeHTTP(rw, r)
	})
}

// WrappedReponseWriter a wrapper for GzipMiddleware ResponseWriter
type WrappedReponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

// NewWrappedResponseWriter ctor for WrappedResponseWriter
func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedReponseWriter {
	gw := gzip.NewWriter(rw)

	return &WrappedReponseWriter{rw: rw, gw: gw}
}

// Header implementation for WrappedResponseWriter
func (wr *WrappedReponseWriter) Header() http.Header {
	return wr.rw.Header()
}

// Write implementation for WrappedResponseWriter
func (wr *WrappedReponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}

// WriteHeader implementation for WrappedResponseWriter
func (wr *WrappedReponseWriter) WriteHeader(statuscode int) {
	wr.rw.WriteHeader(statuscode)
}

// Flush implementation for WrappedResponseWriter
func (wr *WrappedReponseWriter) Flush() {
	wr.gw.Flush()
	wr.gw.Close()
}
