package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"learningAssistant-backend/database"
	"learningAssistant-backend/middleware"
	"learningAssistant-backend/models"
)

func registerNotificationRoutes(r *gin.RouterGroup) {
	r.Use(middleware.AuthMiddleware())

	r.GET("", listNotifications)
	r.GET("/unread-count", getUnreadNotificationCount)
	r.PUT("/:id/read", markNotificationAsRead)
	r.PUT("/read-all", markAllNotificationsAsRead)
	r.DELETE("/clear-all", clearAllNotifications)
	r.DELETE("/:id", deleteNotification)
}

func listNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var notifications []models.Notification
	var total int64

	db := database.GetDB().Model(&models.Notification{}).Where("user_id = ?", uid)
	db.Count(&total)

	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":     notifications,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func getUnreadNotificationCount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	var count int64
	database.GetDB().Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", uid, false).Count(&count)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": count})
}

func markNotificationAsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	id := c.Param("id")

	if err := database.GetDB().Model(&models.Notification{}).Where("id = ? AND user_id = ?", id, uid).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "已标记为已读"})
}

func markAllNotificationsAsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	if err := database.GetDB().Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", uid, false).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "全部标记为已读"})
}

func deleteNotification(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	id := c.Param("id")

	if err := database.GetDB().Where("id = ? AND user_id = ?", id, uid).Delete(&models.Notification{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
}

func clearAllNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	if err := database.GetDB().Where("user_id = ?", uid).Delete(&models.Notification{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清空失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "已清空所有通知"})
}
