package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"sloth-tracker/api/model"
	"time"
)

// 注册设备 POST
func RegisterDevice(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			OwnerId     string `json:"ownerId"`
			DeviceName  string `json:"deviceName"`
			Platform    string `json:"platform"`
			Description string `json:"description"`
		}
		// 校验参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		deviceId := uuid.New().String()
		device := model.Device{
			Id:           deviceId,
			OwnerId:      req.OwnerId,
			Name:         req.DeviceName,
			Platform:     req.Platform,
			Description:  req.Description,
			RegisteredAt: time.Now(),
		}
		db.Create(&device)

		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "注册成功", "device_id": deviceId})
	}
}

// 修改设备信息 PUT
func UpdateDeviceInfo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			DeviceId    string `json:"deviceId"`
			Name        string `json:"name"`
			Platform    string `json:"platform"`
			Description string `json:"description"`
		}
		// 校验参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 更新设备信息
		db.Where("id = ?", req.DeviceId).Updates(&model.Device{
			Name:        req.Name,
			Platform:    req.Platform,
			Description: req.Description,
		})
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "设备信息更新成功"})
	}
}

// 获取设备列表 GET
func GetDeviceList(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userId := c.Param("user_id")
		if userId == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 查询设备列表
		var devices []model.Device
		db.Where("owner_id = ?", userId).Find(&devices)
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "查询成功", "devices": devices})
	}
}

// 获取共享设备列表 GET
func GetSharedDeviceList(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userId := c.Param("user_id")
		if userId == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 获取共享给用户的设备
		var deviceAccesses []model.DeviceAccess
		db.Where("viewer_id = ?", userId).Find(&deviceAccesses)
		// 提取设备ID
		var deviceIds []string
		for _, access := range deviceAccesses {
			deviceIds = append(deviceIds, access.DeviceId)
		}
		// 查询设备信息
		var devices []model.Device
		db.Where("id IN ?", deviceIds).Find(&devices)
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "查询成功", "devices": devices})
	}
}

// 获取设备信息 GET
func GetDeviceInfo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取设备ID
		deviceId := c.Param("device_id")
		if deviceId == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 查询设备
		var device model.Device
		if err := db.Where("id = ?", deviceId).First(&device).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "设备不存在"})
			return
		}
		// 返回设备信息
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "查询成功", "device": device})
	}
}

// 注销设备 DELETE
func DeleteDevice(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Id string `json:"id"`
		}
		// 校验参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 查询设备
		var device model.Device
		if err := db.Where("id = ?", req.Id).First(&device).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "设备不存在"})
			return
		}
		tx := db.Begin()
		// 删除设备
		if err := tx.Delete(&device).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 4, "message": "设备注销失败-删除设备失败"})
			return
		}
		// 删除设备状态
		if err := tx.Where("id = ?", req.Id).Delete(&model.DeviceStatus{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 5, "message": "设备注销失败-删除设备状态失败"})
			return
		}
		// 提交事务
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "注销成功"})
	}
}
