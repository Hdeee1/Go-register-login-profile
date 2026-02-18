package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips	map[string]*rate.Limiter
	mu	*sync.RWMutex
	r	rate.Limit
	b	int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:	 &sync.RWMutex{},
		r:	 r,
		b:	 b,
	}
}

func (limiter *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	limiter.mu.Lock()

	limit, exist := limiter.ips[ip]
	if !exist {
		limit = rate.NewLimiter(limiter.r, limiter.b)
		limiter.ips[ip] = limit
	}

	limiter.mu.Unlock()

	return limit
}

func RateLimiterMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		rateLimiter := limiter.getLimiter(ip)
		if !rateLimiter.Allow() {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many request"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}