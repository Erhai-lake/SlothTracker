package controller

import (
	"net/http"
	"sloth-tracker/api/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 获取设备状态
func GetStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		userID := c.Param("user_id")
		// 校验参数
		if deviceID == "" || userID == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 检查设备是否归属用户
		var device model.Device
		err := db.Where("id = ? AND owner_id = ?", deviceID, userID).First(&device).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// 不是所有者, 检查是否为已授权的共享设备
				var sharedDevice model.SharedDevice
				err = db.Where("device_id = ? AND viewer_id = ? AND authorization = 1", deviceID, userID).
					First(&sharedDevice).Error
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						// 既不是设备所有者, 也不是授权用户, 权限不足
						c.JSON(http.StatusForbidden, gin.H{"code": 4, "message": "无权获取该设备状态"})
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"code": 3, "message": "数据库查询错误"})
					}
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 3, "message": "查询数据库出错"})
				return
			}
		}
		// 查询设备状态
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

// 更新设备状态 PUT
func UpdateStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if deviceID == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "device_id 参数不能为空"})
			return
		}
		var req model.DeviceStatus
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}

		var existing model.DeviceStatus
		err := db.Where("device_id = ?", deviceID).First(&existing).Error
		// now时间戳获取到毫秒
		now := time.Now().UnixNano() / 1e6
		if err == nil {
			updateData := map[string]any{
				"timestamp":                   now,
				"battery_charging":            req.Battery.Charging,
				"battery_level":               req.Battery.Level,
				"battery_temperature":         req.Battery.Temperature,
				"battery_capacity":            req.Battery.Capacity,
				"network_wifi_connected":      req.Network.WifiConnected,
				"network_wifi_ss_id":          req.Network.WifiSSId,
				"network_mobile_data_active":  req.Network.MobileDataActive,
				"network_mobile_signal_dbm":   req.Network.MobileSignalDbm,
				"network_network_type":        req.Network.NetworkType,
				"network_upload_speed_kbps":   req.Network.UploadSpeedKbps,
				"network_download_speed_kbps": req.Network.DownloadSpeedKbps,
				"network_traffic_used_mb":     req.Network.TrafficUsedMB,
				"foreground_app_name":         req.Foreground.AppName,
				"foreground_app_title":        req.Foreground.AppTitle,
				"foreground_speaker_playing":  req.Foreground.SpeakerPlaying,
				"other_screen_on":             req.Other.ScreenOn,
				"other_is_charging_via_usb":   req.Other.IsChargingViaUSB,
				"other_is_charging_via_ac":    req.Other.IsChargingViaAC,
				"other_is_low_power_mode":     req.Other.IsLowPowerMode,
			}
			if err := db.Model(&existing).Updates(updateData).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "更新失败"})
				return
			}
		} else {
			// 不存在, 创建新记录
			req.Id = uuid.New().String()
			req.DeviceId = deviceID
			req.Timestamp = now
			if err := db.Create(&req).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "创建失败"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "状态更新成功"})
	}
}
