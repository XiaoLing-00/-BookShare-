package controllers

import (
	"bookshare/config"
	"bookshare/models"
	"net/http"
	//"strconv"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
)

// AddUserBookRelation godoc
// @Summary 添加用户书籍关系
// @Description 收藏或标记书籍已读
// @Tags 关系
// @Accept json
// @Produce json
// @Param relation body models.UserBookRelation true "用户书籍关系信息"
// @Success 201 {object} models.UserBookRelation
// @Failure 400 {object} gin.H "请求参数错误"
// @Failure 409 {object} gin.H "关系已存在"
// @Failure 500 {object} gin.H "创建关系失败"
// @Router /relations [post]
func AddUserBookRelation(c *gin.Context) {
	var relation models.UserBookRelation
	if err := c.ShouldBindJSON(&relation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否已存在相同的关系
	var existingRelation models.UserBookRelation
	if config.DB.Where("user_id = ? AND book_id = ? AND relation_type = ?", relation.UserID, relation.BookID, relation.RelationType).First(&existingRelation).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Relation already exists"})
		return
	}

	if result := config.DB.Create(&relation); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add relation"})
		return
	}
	c.JSON(http.StatusCreated, relation)
}

// GetUserRelations godoc
// @Summary 获取用户的所有书籍关系
// @Description 根据用户ID获取其所有收藏或已读的书籍
// @Tags 关系
// @Produce json
// @Param user_id path int true "用户ID"
// @Success 200 {array} models.UserBookRelation
// @Router /users/{user_id}/relations [get]
func GetUserRelations(c *gin.Context) {
	// Router uses :id as the path parameter
	userID := c.Param("id")
	var relations []models.UserBookRelation
	if result := config.DB.Preload("Book").Preload("User").Where("user_id = ?", userID).Find(&relations); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user relations"})
		return
	}
	c.JSON(http.StatusOK, relations)
}

// GetUserRelationByType godoc
// @Summary 获取用户特定类型的书籍关系
// @Description 根据用户ID和关系类型获取书籍（收藏或已读）
// @Tags 关系
// @Produce json
// @Param user_id path int true "用户ID"
// @Param type path string true "关系类型 (e.g., collected, read)"
// @Success 200 {array} models.UserBookRelation
// @Router /users/{user_id}/relations/{type} [get]
func GetUserRelationByType(c *gin.Context) {
	// Router uses :id as the path parameter
	userID := c.Param("id")
	relationType := c.Param("type")
	var relations []models.UserBookRelation
	if result := config.DB.Preload("Book").Preload("User").Where("user_id = ? AND relation_type = ?", userID, relationType).Find(&relations); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user relations by type"})
		return
	}
	c.JSON(http.StatusOK, relations)
}

// DeleteUserBookRelation godoc
// @Summary 删除用户书籍关系
// @Description 根据关系ID删除用户书籍关系
// @Tags 关系
// @Produce json
// @Param id path int true "关系ID"
// @Success 204 "删除成功"
// @Failure 404 {object} gin.H "关系未找到"
// @Router /relations/{id} [delete]
func DeleteUserBookRelation(c *gin.Context) {
	id := c.Param("id")
	var relation models.UserBookRelation
	if err := config.DB.First(&relation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Relation not found"})
		return
	}
	config.DB.Delete(&relation) // 软删除
	c.Status(http.StatusNoContent)
}
