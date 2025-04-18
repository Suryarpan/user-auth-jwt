package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Suryarpan/user-auth-jwt/utils"
)

func ReqLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := NewWrapResponseWriter(w, r.ProtoMajor)
		s := time.Now()

		next.ServeHTTP(ww, r)

		d := time.Since(s)
		scheme := utils.If(r.TLS != nil, "https", "http")
		slog.Info("handled request",
			"status", ww.Status(),
			"method", r.Method,
			"path", fmt.Sprintf("%s://%s%s %s", scheme, r.Host, r.RequestURI, r.Proto),
			"time", fmt.Sprintf("%dms", d.Milliseconds()),
			"size", fmt.Sprintf("%db", ww.BytesWritten()),
		)
	})
}
