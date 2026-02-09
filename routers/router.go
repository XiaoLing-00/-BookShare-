package routers

import (
	"bookshare/controllers"
	"bookshare/middlewares" // 导入中间件包
	"net/http"              // 用于 CORS 配置或避免“imported and not used”警告

	"github.com/gin-gonic/gin"
	// 如果需要 Swagger 文档，取消以下注释并安装依赖：
	// go get -u github.com/swaggo/gin-swagger
	// go get -u github.com/swaggo/swag/cmd/swag
	// _ "bookshare/docs" // swagger 文档生成后会有一个 docs 目录
	// ginSwagger "github.com/swaggo/gin-swagger"
	// swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

// InitRouter 初始化所有路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// --- CORS 配置 ---

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // 处理 OPTIONS 预检请求
			return
		}

		c.Next() // 继续处理请求
	})
	// --- CORS 配置结束 ---

	// 为了避免 "net/http" imported and not used 警告，如果不需要在其他地方使用
	_ = http.StatusNoContent

	// 用户相关接口 (无需认证)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// User Group - 应用认证中间件
	userRoutes := r.Group("/users")
	userRoutes.Use(middlewares.AuthMiddleware()) // 启用认证
	{
		// 1. 先注册更具体的、包含额外路径段的路由
		// 将 :user_id 改为 :id 以保持参数名一致
		// Gin 会将路径中的 :id 匹配到对应的路径参数
		userRoutes.GET("/:id/books", controllers.GetBooksByUser)
		userRoutes.GET("/:id/relations", controllers.GetUserRelations)
		userRoutes.GET("/:id/relations/:type", controllers.GetUserRelationByType)

		// 2. 然后再注册只包含单个通配符的通用路由
		// 所有参数都使用 :id
		userRoutes.GET("/:id", controllers.GetUserProfile)
		userRoutes.PUT("/:id", controllers.UpdateUserProfile)
		userRoutes.DELETE("/:id", controllers.DeleteUser)
	}

	// Book Group
	bookRoutes := r.Group("/books")
	// 示例：书籍创建/更新/删除需要认证，但浏览通常不需要。
	// 这里假设所有书籍操作都需要认证，你可以根据实际需求调整。
	bookRoutes.Use(middlewares.AuthMiddleware())
	{
		bookRoutes.POST("", controllers.CreateBook)
		bookRoutes.GET("", controllers.GetAllBooks)
		bookRoutes.GET("/:id", controllers.GetBookByID)
		bookRoutes.PUT("/:id", controllers.UpdateBook)
		bookRoutes.DELETE("/:id", controllers.DeleteBook)
		bookRoutes.GET("/category/:category", controllers.GetBooksByCategory)
	}

	// Comment Group
	commentRoutes := r.Group("/comments")
	commentRoutes.Use(middlewares.AuthMiddleware()) // 评论操作需要认证
	{
		commentRoutes.POST("", controllers.AddComment)
		// 注意：路由路径应该是 /comments/book/:book_id，而不是 /comments/:book_id/comments
		commentRoutes.GET("/book/:book_id", controllers.GetCommentsByBookID)
		commentRoutes.DELETE("/:id", controllers.DeleteComment)
	}

	// Relation Group (用户收藏/阅读记录)
	relationRoutes := r.Group("/relations")
	relationRoutes.Use(middlewares.AuthMiddleware()) // 关系操作需要认证
	{
		relationRoutes.POST("", controllers.AddUserBookRelation)
		relationRoutes.DELETE("/:id", controllers.DeleteUserBookRelation)
	}

	// Admin Group - 应用管理员认证中间件
	adminRoutes := r.Group("/admin/stats")
	adminRoutes.Use(middlewares.AdminAuthMiddleware()) // 启用管理员认证
	{
		adminRoutes.GET("/users/count", controllers.GetUserCount)
		adminRoutes.GET("/books/count", controllers.GetBookCount)
		adminRoutes.GET("/comments/count", controllers.GetCommentCount)
		adminRoutes.GET("/users/latest", controllers.GetLatestUsers)
		adminRoutes.GET("/books/latest", controllers.GetLatestBooks)
		adminRoutes.GET("/books/popular", controllers.GetPopularBooks)
	}

	// --- Swagger Docs 配置 (可选) ---
	// 确保已安装 github.com/swaggo/gin-swagger 和 github.com/swaggo/swag/cmd/swag
	// 1. 在项目根目录运行 `swag init` 生成 docs 目录
	// 2. 取消以下注释，访问 http://localhost:8080/swagger/index.html
	/*
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	*/
	// --- Swagger Docs 配置结束 ---

	return r
}
