package middleware

import (
	"AggregationService/internal/pkg/logger"
	"fmt"
	"net/http"
	"time"
)

func LoggerMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		latency := time.Since(start)
		logger.FromContext(r.Context()).
			Info(fmt.Sprintf("--- %s --- %s %s | %v", r.Method, r.RequestURI, r.RemoteAddr, latency))
	})
}
