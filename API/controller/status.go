package controller

import (
	"net/http"
	"sloth-tracker/api/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取设备状态
func GetStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if deviceID == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		var status model.DeviceStatus
		result := db.Where("device_id = ?", deviceID).First(&status)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"code": 2, "message": "设备状态未找到"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 3, "message": "查询数据库出错"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "查询成功", "status": status})
	}
}
