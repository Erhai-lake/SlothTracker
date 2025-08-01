package controller

import (
	"net/http"
	"sloth-tracker/api/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 删除设备
func DeleteDevice(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if deviceID == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "device_id 参数不能为空",
			})
			return
		}

		// 开启事务
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("id = ?", deviceID).Delete(&model.Device{}).Error; err != nil {
				return err
			}
			if err := tx.Where("id = ?", deviceID).Delete(&model.DeviceStatus{}).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    2,
				"message": "删除设备失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "设备删除成功",
		})
	}
}
