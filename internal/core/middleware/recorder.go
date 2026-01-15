package middleware

import "net/http"

type StatusRecorder struct {
	ResponseWriter http.ResponseWriter
	statusCode     int
}

func NewStatusRecorder(rw http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{
		ResponseWriter: rw,
		statusCode:     http.StatusOK,
	}
}

func (rw *StatusRecorder) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *StatusRecorder) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

func (rw *StatusRecorder) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
