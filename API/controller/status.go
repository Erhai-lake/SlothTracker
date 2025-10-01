package controller

import (
	"encoding/json"
	"net/http"
	"sloth-tracker/api/model"
	"sloth-tracker/api/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 获取设备状态
func GetStatus(db any) http.HandlerFunc {
	type DeviceStatusWithSource struct {
		Source string `json:"source"`
		model.DeviceStatus
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		// 从路径参数获取参数
		userID := utils.GetPathParam(r, 2)
		deviceID := utils.GetPathParam(r, 3)

		// 校验参数
		if deviceID == "" || userID == "" {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)
		source := "账户"

		// 检查设备是否归属用户
		var device model.Device
		result := gormDB.Where("id = ? AND owner_id = ?", deviceID, userID).First(&device)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// 不是所有者, 检查是否为已授权的共享设备
				var sharedDevice model.SharedDevice
				result = gormDB.Where("device_id = ? AND viewer_id = ? AND authorization = 1", deviceID, userID).First(&sharedDevice)
				source = "共享"
				if result.Error != nil {
					if result.Error == gorm.ErrRecordNotFound {
						// 既不是设备所有者, 也不是授权用户, 权限不足
						utils.Error(w, http.StatusForbidden, "无权获取该设备状态")
					} else {
						utils.Error(w, http.StatusInternalServerError, "数据库查询错误")
					}
					return
				}
			} else {
				utils.Error(w, http.StatusInternalServerError, "查询数据库出错")
				return
			}
		}

		// 查询设备状态
		var status DeviceStatusWithSource
		status.DeviceStatus = model.DeviceStatus{}
		result = gormDB.Where("device_id = ?", deviceID).First(&status.DeviceStatus)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				utils.Error(w, http.StatusNotFound, "设备状态未找到")
			} else {
				utils.Error(w, http.StatusInternalServerError, "查询数据库出错")
			}
			return
		}

		status.Source = source

		utils.Success(w, map[string]any{
			"message": "查询成功",
			"status":  status,
		})
	}
}

// 更新设备状态 PUT
func UpdateStatus(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		// 从路径参数获取参数
		userID := utils.GetPathParam(r, 3)
		deviceID := utils.GetPathParam(r, 4)

		// 校验参数
		if deviceID == "" || userID == "" {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		var req model.DeviceStatus
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 检查设备是否归属用户
		var device model.Device
		result := gormDB.Where("id = ? AND owner_id = ?", deviceID, userID).First(&device)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				utils.Error(w, http.StatusForbidden, "无权更新该设备状态")
				return
			} else {
				utils.Error(w, http.StatusInternalServerError, "查询数据库出错")
				return
			}
		}

		var existing model.DeviceStatus
		err := gormDB.Where("device_id = ?", deviceID).First(&existing).Error

		// now时间戳获取到毫秒
		now := time.Now().UnixNano() / 1e6

		if err == nil {
			// 更新现有记录
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

			if err := gormDB.Model(&existing).Updates(updateData).Error; err != nil {
				utils.Error(w, http.StatusInternalServerError, "更新失败")
				return
			}
		} else {
			// 不存在, 创建新记录
			req.Id = uuid.New().String()
			req.DeviceId = deviceID
			req.Timestamp = now
			if err := gormDB.Create(&req).Error; err != nil {
				utils.Error(w, http.StatusInternalServerError, "创建失败")
				return
			}
		}

		utils.Success(w, map[string]any{
			"message": "状态更新成功",
		})
	}
}
