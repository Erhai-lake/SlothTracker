package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"sloth-tracker/api/model"
	"time"
)

// 注册设备
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string `json:"deviceName"`
			Platform    string `json:"platform"`
			Description string `json:"description"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "参数错误"})
			return
		}

		deviceID := uuid.New().String()
		device := model.Device{
			ID:           deviceID,
			Name:         req.Name,
			Platform:     req.Platform,
			Description:  req.Description,
			RegisteredAt: time.Now(),
		}
		db.Create(&device)

		c.JSON(http.StatusOK, gin.H{"code": 0, "deviceId": deviceID})
	}
}
