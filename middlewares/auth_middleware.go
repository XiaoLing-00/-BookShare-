// bookshare/middlewares/auth_middleware.go
package middlewares

import (
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 实际的认证逻辑在此处，目前为了演示，直接放行
		// c.Set("user_id", uint(1)) // 模拟设置用户ID，以便后续控制器使用
		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 实际的管理员认证逻辑在此处，目前为了演示，直接放行
		c.Next()
	}
}
