package main

import (
	"bookshare/config"
	"bookshare/models"
	"bookshare/routers"
	"log"
)

// @title BookShare API
// @version 1.0
// @description 在线书籍管理与分享平台API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	config.InitDB()    // 初始化数据库连接
	config.InitRedis() // 初始化Redis连接

	// 自动迁移模型，创建或更新表结构
	err := config.DB.AutoMigrate(&models.User{}, &models.Book{}, &models.Comment{}, &models.UserBookRelation{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
	log.Println("Database migration completed!")

	r := routers.InitRouter() // 初始化路由

	log.Println("Gin server started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start Gin server: %v", err)
	}
}