package middleware

import (
	"log"
	"net"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func Logging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		userAgent := r.UserAgent()

		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         0,
		}

		h(recorder, r)

		duration := time.Since(start).Milliseconds()

		log.Println(ip, r.Method, r.RequestURI, r.Proto, recorder.Status, duration, userAgent)
	}
}
