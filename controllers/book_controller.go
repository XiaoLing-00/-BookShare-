package controllers

import (
	"bookshare/config"
	"bookshare/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateBook godoc
// @Summary 创建新书籍
// ... (完整的 CreateBook 函数)
func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := config.DB.Create(&book); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}
	c.JSON(http.StatusCreated, book)
}

// GetBookByID godoc
// @Summary 获取书籍详情
// ... (完整的 GetBookByID 函数)
func GetBookByID(c *gin.Context) {
	id := c.Param("id")
	bookID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	cacheKey := "book:" + id
	val, err := config.RDB.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var book models.Book
		if err := json.Unmarshal([]byte(val), &book); err == nil {
			c.JSON(http.StatusOK, book)
			return
		}
	}

	var book models.Book
	if err := config.DB.Preload("User").First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	bookJSON, err := json.Marshal(book)
	if err == nil {
		config.RDB.Set(config.Ctx, cacheKey, bookJSON, 1*time.Hour)
	}

	c.JSON(http.StatusOK, book)
}

// GetAllBooks godoc
// @Summary 获取所有书籍
// ... (完整的 GetAllBooks 函数)
func GetAllBooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	category := c.Query("category")

	offset := (page - 1) * pageSize
	var books []models.Book
	query := config.DB.Model(&models.Book{}).Preload("User")

	if keyword != "" {
		search := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR author LIKE ? OR description LIKE ?", search, search, search)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if result := query.Limit(pageSize).Offset(offset).Find(&books); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

// UpdateBook godoc
// @Summary 更新书籍信息
// ... (完整的 UpdateBook 函数)
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := config.DB.Model(&book).Updates(updatedBook); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// DeleteBook godoc
// @Summary 删除书籍
// ... (完整的 DeleteBook 函数)
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	config.DB.Delete(&book)
	c.Status(http.StatusNoContent)
}

// GetBooksByUser godoc
// @Summary 获取用户上传的书籍
// ... (完整的 GetBooksByUser 函数)
func GetBooksByUser(c *gin.Context) {
	// Router uses :id as the path parameter
	userID := c.Param("id")
	var books []models.Book
	if result := config.DB.Where("user_id = ?", userID).Find(&books); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

// GetBooksByCategory godoc
// @Summary 获取指定分类的书籍
// ... (完整的 GetBooksByCategory 函数)
func GetBooksByCategory(c *gin.Context) {
	category := c.Param("category")
	var books []models.Book
	if result := config.DB.Where("category = ?", category).Find(&books); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books by category"})
		return
	}
	c.JSON(http.StatusOK, books)
}
