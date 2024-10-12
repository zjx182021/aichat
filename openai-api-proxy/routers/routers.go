package routers

import (
	"github.com/gin-gonic/gin"
	"openai-api-proxy/health"
	"openai-api-proxy/middleware"
	"openai-api-proxy/proxy"
)

func InitRouters(r *gin.Engine) {
	r.GET("/health", health.Health)
	r.Use(middleware.Auth(), middleware.RateLimit(10, 10))
	initProxyRouter(r)
}
func initProxyRouter(r *gin.Engine) {
	p := proxy.NewProxy()
	v1 := r.Group("/v1")
	// *relativePath 是一个通配符，用于匹配任意的相对路径
	v1.Any("/*relativePath", p.ChatProxy)
}
