package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mock-openai-api/pkg/config"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cnf := config.GetConfig()
		authorization := ctx.Request.Header.Get("Authorization")
		confAuthorization := fmt.Sprintf("Bearer %s", cnf.Http.AccessToken)
		if authorization != confAuthorization {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}
