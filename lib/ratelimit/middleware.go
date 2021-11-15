package ratelimit

import (
	"io"
	"net/http"
	"strconv"
	"time"
)

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lim := TempSession.ByAddress(r.RemoteAddr)
		if lim.Limited() {
			w.Header().Set("X-RateLimit-Global", "true")
			w.Header().Set("Retry-After", strconv.Itoa(int(time.Until(lim.Expiration()).Seconds())))
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(MaxRequestsPerMinute))
			w.WriteHeader(http.StatusTooManyRequests)

			io.WriteString(w, "You are being rate limited.")
			return
		}
		if lim.Expired() {
			lim.Reset()
		}
		lim.AddAmount()
		h.ServeHTTP(w, r)
	})
}
