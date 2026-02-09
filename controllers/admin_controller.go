package controllers

import (
	"bookshare/config"
	"bookshare/models"
	"encoding/json"
	"fmt" // 用于Sprintf
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserCount godoc
// @Summary 获取用户总数
// @Description 获取平台注册用户总数
// @Tags 后台数据统计
// @Produce json
// @Success 200 {object} gin.H "用户总数"
// @Router /admin/stats/users/count [get]
func GetUserCount(c *gin.Context) {
	var count int64
	config.DB.Model(&models.User{}).Count(&count)
	c.JSON(http.StatusOK, gin.H{"user_count": count})
}

// GetBookCount godoc
// @Summary 获取书籍总数
// @Description 获取平台书籍总数
// @Tags 后台数据统计
// @Produce json
// @Success 200 {object} gin.H "书籍总数"
// @Router /admin/stats/books/count [get]
func GetBookCount(c *gin.Context) {
	var count int64
	config.DB.Model(&models.Book{}).Count(&count)
	c.JSON(http.StatusOK, gin.H{"book_count": count})
}

// GetCommentCount godoc
// @Summary 获取评论总数
// @Description 获取平台评论总数
// @Tags 后台数据统计
// @Produce json
// @Success 200 {object} gin.H "评论总数"
// @Router /admin/stats/comments/count [get]
func GetCommentCount(c *gin.Context) {
	var count int64
	config.DB.Model(&models.Comment{}).Count(&count)
	c.JSON(http.StatusOK, gin.H{"comment_count": count})
}

// GetLatestUsers godoc
// @Summary 获取最新注册用户
// @Description 获取最新注册的用户列表
// @Tags 后台数据统计
// @Produce json
// @Param limit query int false "限制数量" default(5)
// @Success 200 {array} models.User
// @Router /admin/stats/users/latest [get]
func GetLatestUsers(c *gin.Context) {
	limit := c.DefaultQuery("limit", "5")
	var users []models.User
	config.DB.Limit(toInt(limit)).Order("created_at desc").Find(&users)
	for i := range users {
		users[i].Password = "" // 不返回密码
	}
	c.JSON(http.StatusOK, users)
}

// GetLatestBooks godoc
// @Summary 获取最新上传书籍
// @Description 获取最新上传的书籍列表
// @Tags 后台数据统计
// @Produce json
// @Param limit query int false "限制数量" default(5)
// @Success 200 {array} models.Book
// @Router /admin/stats/books/latest [get]
func GetLatestBooks(c *gin.Context) {
	limit := c.DefaultQuery("limit", "5")
	var books []models.Book
	config.DB.Limit(toInt(limit)).Order("created_at desc").Find(&books)
	c.JSON(http.StatusOK, books)
}

// GetPopularBooks godoc
// @Summary 获取热门书籍
// @Description 根据被收藏或被阅读次数（简化为被评论次数）获取热门书籍，使用Redis缓存
// @Tags 后台数据统计
// @Produce json
// @Param limit query int false "限制数量" default(5)
// @Success 200 {array} models.Book
// @Router /admin/stats/books/popular [get]
func GetPopularBooks(c *gin.Context) {
	limit := c.DefaultQuery("limit", "5")
	limitInt := toInt(limit)
	cacheKey := fmt.Sprintf("popular_books:limit:%d", limitInt)

	// 1. 尝试从Redis获取
	val, err := config.RDB.Get(config.Ctx, cacheKey).Result()
	if err == nil { // Redis中有缓存
		var popularBooks []models.Book
		if err := json.Unmarshal([]byte(val), &popularBooks); err == nil {
			c.JSON(http.StatusOK, popularBooks)
			return
		}
	}

	// 2. Redis中没有或解析失败，从数据库获取
	var popularBooks []models.Book
	// 注意：这里的Raw查询需要确保GORM能正确Scan到Book结构体
	config.DB.Raw("SELECT b.* FROM books b LEFT JOIN comments c ON b.id = c.book_id GROUP BY b.id ORDER BY COUNT(c.id) DESC LIMIT ?", limitInt).Scan(&popularBooks)

	// 3. 将结果存入Redis，设置过期时间
	bookJSON, err := json.Marshal(popularBooks)
	if err == nil {
		config.RDB.Set(config.Ctx, cacheKey, bookJSON, 5*time.Minute) // 缓存5分钟
	}

	c.JSON(http.StatusOK, popularBooks)
}

// Helper function to convert string to int
func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0 // default or handle error
	}
	return i
}
