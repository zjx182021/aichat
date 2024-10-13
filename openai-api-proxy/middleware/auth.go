package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"openai-api-proxy/pkg/config"
	"time"

	"github.com/gin-gonic/gin"
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
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Intn(len(cnf.Chat.APIKeys))
		apiKey := cnf.Chat.APIKeys[randIndex]
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
		ctx.Next()
	}
}
