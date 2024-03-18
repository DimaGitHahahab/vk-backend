package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func Logging(h http.Handler, log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		log.Infof("request: %s %s", request.Method, request.URL.Path)
		logRespWriter := &loggingResponseWriter{originalRW: responseWriter, header: http.Header{}, err: nil}
		h.ServeHTTP(logRespWriter, request)
		if logRespWriter.err != nil {
			log.Errorf("response: %d, error: %s", logRespWriter.code, logRespWriter.err)
		} else {
			log.Infof("response: %d", logRespWriter.code)
		}

	})
}

type loggingResponseWriter struct {
	originalRW http.ResponseWriter
	header     http.Header
	code       int
	err        error
}

func (rw *loggingResponseWriter) Header() http.Header {
	rw.header = rw.originalRW.Header()
	return rw.originalRW.Header()
}

func (rw *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := rw.originalRW.Write(data)
	if err != nil {
		rw.err = err
	}
	return size, err
}

func (rw *loggingResponseWriter) WriteHeader(code int) {
	rw.code = code
	rw.originalRW.WriteHeader(code)
}
