package middleware

import (
	"strings"

	"github.com/BenedictKing/ccx/internal/config"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS 中间件
func CORSMiddleware(envCfg *config.EnvConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			setCORSHeaders(c, envCfg, origin)
		}

		// OPTIONS 预检请求不携带自定义认证头，始终直接返回 CORS 响应。
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// 如果未启用 CORS，非跨域请求保持原有行为。
		if !envCfg.EnableCORS && origin == "" {
			c.Next()
			return
		}

		c.Next()
	}
}

func setCORSHeaders(c *gin.Context, envCfg *config.EnvConfig, origin string) {
	allowOrigin := envCfg.CORSOrigin
	if envCfg.IsDevelopment() && strings.Contains(origin, "localhost") {
		allowOrigin = origin
	}
	if !envCfg.EnableCORS {
		allowOrigin = origin
	}
	c.Header("Access-Control-Allow-Origin", allowOrigin)
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, x-api-key, x-goog-api-key")
	if allowOrigin != "*" {
		c.Header("Access-Control-Allow-Credentials", "true")
	}
}
