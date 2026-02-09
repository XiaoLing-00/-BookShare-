package controllers

import (
	"bookshare/config"
	"bookshare/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddComment godoc
// @Summary 添加书籍评论
// @Description 为指定书籍添加评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param comment body models.Comment true "评论信息"
// @Success 201 {object} models.Comment
// @Failure 400 {object} gin.H "请求参数错误"
// @Failure 500 {object} gin.H "添加评论失败"
// @Router /comments [post]
func AddComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := config.DB.Create(&comment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

// GetCommentsByBookID godoc
// @Summary 获取书籍评论
// @Description 根据书籍ID获取所有评论
// @Tags 评论
// @Produce json
// @Param book_id path int true "书籍ID"
// @Success 200 {array} models.Comment
// @Router /books/{book_id}/comments [get]
func GetCommentsByBookID(c *gin.Context) {
	bookID := c.Param("book_id")
	var comments []models.Comment
	if result := config.DB.Preload("User").Where("book_id = ?", bookID).Find(&comments); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// DeleteComment godoc
// @Summary 删除评论
// @Description 根据评论ID删除评论
// @Tags 评论
// @Produce json
// @Param id path int true "评论ID"
// @Success 204 "删除成功"
// @Failure 404 {object} gin.H "评论未找到"
// @Router /comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	var comment models.Comment
	if err := config.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}
	config.DB.Delete(&comment) // 软删除
	c.Status(http.StatusNoContent)
}