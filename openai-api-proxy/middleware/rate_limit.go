package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

func RateLimit(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)
	return func(ctx *gin.Context) {
		if !limiter.Allow() {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
