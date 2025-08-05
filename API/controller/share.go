package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"sloth-tracker/api/model"
	"time"
)

// 申请贡献 POST
func ApplyShare(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			DeviceId string `json:"deviceId"`
			ViewerId string `json:"viewerId"`
		}
		// 校验参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 检查设备是否存在
		var device model.Device
		db.Where("id = ?", req.DeviceId).First(&device)
		if device.Id == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "设备不存在"})
			return
		}
		// 检查用户是否存在
		var viewer model.User
		db.Where("id = ?", req.ViewerId).First(&viewer)
		if viewer.Id == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "用户不存在"})
			return
		}
		// 禁止申请自己的设备
		if device.OwnerId == req.ViewerId {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "禁止申请自己的设备"})
			return
		}
		// 检查是否已存在授权记录
		var existingShared model.SharedDevice
		db.Where("device_id = ? AND viewer_id = ?", req.DeviceId, req.ViewerId).First(&existingShared)
		if existingShared.Id != "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "已存在授权记录"})
			return
		}
		// 创建授权记录
		shared := model.SharedDevice{
			Id:            uuid.New().String(),
			DeviceId:      req.DeviceId,
			ViewerId:      req.ViewerId,
			Authorization: 2,
			CreatedAt:     time.Now(),
		}
		db.Create(&shared)
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "申请分享成功, 等待设备所有者授权"})
	}
}

// 获取用户申请的授权 GET
func GetUserApplications(db *gorm.DB) gin.HandlerFunc {
	type ApplicationInfo struct {
		Id         string    `json:"id"`
		DeviceId   string    `json:"device_id"`
		Status     int       `json:"status"`
		UserName   string    `json:"user_name"`
		DeviceName string    `json:"device_name"`
		CreatedAt  time.Time `json:"created_at"`
	}
	return func(c *gin.Context) {
		// 获取用户ID
		userId := c.Param("user_id")
		if userId == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 查询用户申请的授权
		var sharedDevice []model.SharedDevice
		db.Where("viewer_id = ?", userId).Find(&sharedDevice)
		// 构建返回结果
		var result []ApplicationInfo
		for _, auth := range sharedDevice {
			// 获取申请人信息
			var user model.User
			db.Where("id = ?", auth.ViewerId).First(&user)
			// 获取设备信息
			var device model.Device
			db.Where("id = ?", auth.DeviceId).First(&device)
			result = append(result, ApplicationInfo{
				Id:         auth.Id,
				DeviceId:   auth.DeviceId,
				Status:     auth.Authorization,
				UserName:   user.Name,
				DeviceName: device.Name,
				CreatedAt:  auth.CreatedAt,
			})
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "获取申请的授权成功", "applications": result})
	}
}

// 获取共享授权列表 GET
func GetSharedAuthorizations(db *gorm.DB) gin.HandlerFunc {
	type AuthorizationInfo struct {
		Id         string    `json:"id"`
		Status     int       `json:"status"`
		UserName   string    `json:"user_name"`
		DeviceName string    `json:"device_name"`
		CreatedAt  time.Time `json:"created_at"`
	}
	return func(c *gin.Context) {
		// 获取用户ID
		userId := c.Param("user_id")
		if userId == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 获取用户的所有设备ID
		var deviceIds []string
		db.Model(&model.Device{}).Where("owner_id = ?", userId).Pluck("id", &deviceIds)
		// 查询这些设备的共享授权申请
		var sharedDevice []model.SharedDevice
		db.Where("device_id IN ?", deviceIds).Find(&sharedDevice)
		// 构建返回结果
		var result []AuthorizationInfo
		for _, auth := range sharedDevice {
			// 获取申请人信息
			var user model.User
			db.Where("id = ?", auth.ViewerId).First(&user)
			// 获取设备信息
			var device model.Device
			db.Where("id = ?", auth.DeviceId).First(&device)
			result = append(result, AuthorizationInfo{
				Id:         auth.Id,
				Status:     auth.Authorization,
				UserName:   user.Name,
				DeviceName: device.Name,
				CreatedAt:  auth.CreatedAt,
			})
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "获取共享授权列表成功", "authorizations": result})
	}
}

// 授权 PUT
func AuthorizeDevice(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			AccessId string `json:"id"`
			Status   int    `json:"status"`
		}
		// 校验参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 检查授权记录是否存在
		var shared model.SharedDevice
		db.Where("id = ?", req.AccessId).First(&shared)
		if shared.Id == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "授权记录不存在"})
			return
		}
		// 检查状态参数
		if req.Status != 1 && req.Status != 2 {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 更新授权状态
		shared.Authorization = req.Status
		db.Save(&shared)
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "授权操作成功"})
	}
}

// 删除共享申请 DELETE
func DeleteShare(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			AccessId string `json:"id"`
		}
		// 校验参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 检查授权记录是否存在
		var shared model.SharedDevice
		db.Where("id = ?", req.AccessId).First(&shared)
		if shared.Id == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "授权记录不存在"})
			return
		}
		// 删除共享申请
		db.Delete(&shared)
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除共享申请成功"})
	}
}
