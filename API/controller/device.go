package controller

import (
	"net/http"
	"sloth-tracker/api/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取设备列表
func ListDevice(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var devices []model.Device
		db.Find(&devices)

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"devices": devices,
		})
	}
}
