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

// 注册设备 POST
func RegisterDevice(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		var req struct {
			OwnerId     string `json:"ownerId"`
			DeviceName  string `json:"deviceName"`
			Platform    string `json:"platform"`
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)
		deviceId := uuid.New().String()
		device := model.Device{
			Id:           deviceId,
			OwnerId:      req.OwnerId,
			Name:         req.DeviceName,
			Platform:     req.Platform,
			Description:  req.Description,
			RegisteredAt: time.Now(),
		}
		gormDB.Create(&device)

		utils.Success(w, map[string]any{
			"message":   "注册成功",
			"device_id": deviceId,
		})
	}
}

// 修改设备信息 PUT
func UpdateDeviceInfo(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		var req struct {
			DeviceId    string `json:"deviceId"`
			Name        string `json:"name"`
			Platform    string `json:"platform"`
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 更新设备信息
		result := gormDB.Where("id = ?", req.DeviceId).Updates(&model.Device{
			Name:        req.Name,
			Platform:    req.Platform,
			Description: req.Description,
		})

		if result.RowsAffected == 0 {
			utils.Error(w, http.StatusOK, "设备不存在")
			return
		}

		utils.Success(w, map[string]any{
			"message": "设备信息更新成功",
		})
	}
}

// 获取设备列表 GET
func GetDeviceList(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		// 从查询参数获取user_id
		userId := utils.GetQueryParam(r, "user_id")
		if userId == "" {
			utils.Error(w, http.StatusBadRequest, "参数错误: user_id 不能为空")
			return
		}

		gormDB := db.(*gorm.DB)

		// 查询设备列表
		var devices []model.Device
		gormDB.Where("owner_id = ?", userId).Find(&devices)

		utils.Success(w, map[string]any{
			"message": "查询成功",
			"devices": devices,
		})
	}
}

// 获取共享设备列表 GET
func GetSharedDeviceList(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		// 从查询参数获取user_id
		userId := utils.GetQueryParam(r, "user_id")
		if userId == "" {
			utils.Error(w, http.StatusBadRequest, "参数错误: user_id 不能为空")
			return
		}

		gormDB := db.(*gorm.DB)

		// 获取共享给用户的设备(已授权的设备)
		var sharedDevice []model.SharedDevice
		gormDB.Where("viewer_id = ? AND Authorization = ?", userId, 1).Find(&sharedDevice)

		// 提取设备ID
		var deviceIds []string
		for _, access := range sharedDevice {
			deviceIds = append(deviceIds, access.DeviceId)
		}

		// 如果没有共享设备, 返回空数组
		if len(deviceIds) == 0 {
			utils.Success(w, map[string]interface{}{
				"message": "查询成功",
				"devices": []model.Device{},
			})
			return
		}

		// 查询设备信息
		var devices []model.Device
		gormDB.Where("id IN ?", deviceIds).Find(&devices)

		utils.Success(w, map[string]any{
			"message": "查询成功",
			"devices": devices,
		})
	}
}

// 获取设备信息 GET
func GetDeviceInfo(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		// 从查询参数获取device_id
		deviceId := utils.GetQueryParam(r, "device_id")
		if deviceId == "" {
			utils.Error(w, http.StatusBadRequest, "参数错误: device_id 不能为空")
			return
		}

		gormDB := db.(*gorm.DB)

		// 查询设备
		var device model.Device
		if err := gormDB.Where("id = ?", deviceId).First(&device).Error; err != nil {
			utils.Error(w, http.StatusOK, "设备不存在")
			return
		}

		utils.Success(w, map[string]any{
			"message": "查询成功",
			"device":  device,
		})
	}
}

// 注销设备 DELETE
func DeleteDevice(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		var req struct {
			Id string `json:"id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 查询设备
		var device model.Device
		if err := gormDB.Where("id = ?", req.Id).First(&device).Error; err != nil {
			utils.Error(w, http.StatusOK, "设备不存在")
			return
		}

		tx := gormDB.Begin()

		// 删除设备
		if err := tx.Delete(&device).Error; err != nil {
			tx.Rollback()
			utils.Error(w, http.StatusInternalServerError, "设备注销失败-删除设备失败")
			return
		}

		// 删除设备状态
		if err := tx.Where("id = ?", req.Id).Delete(&model.DeviceStatus{}).Error; err != nil {
			tx.Rollback()
			utils.Error(w, http.StatusInternalServerError, "设备注销失败-删除设备状态失败")
			return
		}

		// 提交事务
		tx.Commit()

		utils.Success(w, map[string]any{
			"message": "注销成功",
		})
	}
}
