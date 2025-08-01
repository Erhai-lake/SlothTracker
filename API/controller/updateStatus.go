package controller

import (
	"net/http"
	"sloth-tracker/api/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 更新设备状态
func UpdateStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if deviceID == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "device_id 参数不能为空",
			})
			return
		}

		var req model.DeviceStatus
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}

		// 强制将 URL 参数中的 deviceID 设置进 req
		req.Id = deviceID

		var existing model.DeviceStatus
		if err := db.Where("id = ?", deviceID).First(&existing).Error; err == nil {
			// 仅更新非零字段, 避免被空值覆盖
			if err := db.Model(&existing).Select("*").Omit("DeviceID").Updates(&req).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "更新失败"})
				return
			}
		} else {
			if err := db.Create(&req).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "创建失败"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "状态更新成功"})
	}
}
