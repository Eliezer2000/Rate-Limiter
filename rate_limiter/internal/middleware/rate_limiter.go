package middleware

import (
	"net/http"
	"rate_limite/internal/limiter"
	"strings"
	"time"
)

type RateLimiterConfig struct {
	Limiter 			*limiter.RateLimiter
	TokenLimits 		map[string]int
	TokenBlockTimes 	map[string]time.Duration
	TokenWindows 		map[string]time.Duration
	DefaultLimit 		int
	DefaultBlockTime 	time.Duration
	DefaultWindow 		time.Duration
}

func RateLimiterMiddleware(config *RateLimiterConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("API_KEY")
			var key string
			var limit int 
			var window, blockTime time.Duration

			if token != "" {
				key = "token:" + token
				limit = config.DefaultLimit
				window = config.DefaultWindow
				blockTime = config.DefaultBlockTime

				if l, ok := config.TokenLimits[token]; ok {
					limit = l
				}
				if w, ok := config.TokenWindows[token]; ok {
					window = w
				}
				if b, ok := config.TokenBlockTimes[token]; ok {
					blockTime = b
				}
			} else {
				ip := r.RemoteAddr

				if strings.Contains(ip, ":") {
					ip = strings.Split(ip, ":")[0]
				}
				key = "ip" + ip
				limit = config.DefaultLimit
				window = config.DefaultWindow
				blockTime = config.DefaultBlockTime
			}

			allowed, _, err := config.Limiter.Allow(key, limit, window, blockTime)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return 
			}
			if !allowed {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
				return 
			}
			next.ServeHTTP(w, r)
		})
	}
}