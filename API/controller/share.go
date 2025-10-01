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

// 申请共享 POST
func ApplyShare(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		var req struct {
			DeviceId string `json:"deviceId"`
			ViewerId string `json:"viewerId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 检查设备是否存在
		var device model.Device
		gormDB.Where("id = ?", req.DeviceId).First(&device)
		if device.Id == "" {
			utils.Error(w, http.StatusOK, "设备不存在")
			return
		}

		// 检查用户是否存在
		var viewer model.User
		gormDB.Where("id = ?", req.ViewerId).First(&viewer)
		if viewer.Id == "" {
			utils.Error(w, http.StatusOK, "用户不存在")
			return
		}

		// 禁止申请自己的设备
		if device.OwnerId == req.ViewerId {
			utils.Error(w, http.StatusOK, "禁止申请自己的设备")
			return
		}

		// 检查是否已存在授权记录
		var existingShared model.SharedDevice
		gormDB.Where("device_id = ? AND viewer_id = ?", req.DeviceId, req.ViewerId).First(&existingShared)
		if existingShared.Id != "" {
			utils.Error(w, http.StatusOK, "已存在授权记录")
			return
		}

		// 创建授权记录
		shared := model.SharedDevice{
			Id:            uuid.New().String(),
			DeviceId:      req.DeviceId,
			ViewerId:      req.ViewerId,
			Authorization: 2, // 2表示待授权
			CreatedAt:     time.Now(),
		}
		gormDB.Create(&shared)

		utils.Success(w, map[string]any{
			"message": "申请分享成功, 等待设备所有者授权",
		})
	}
}

// 获取用户申请的授权 GET
func GetUserApplications(db any) http.HandlerFunc {
	type ApplicationInfo struct {
		Id         string    `json:"id"`
		DeviceId   string    `json:"device_id"`
		Status     int       `json:"status"`
		UserName   string    `json:"user_name"`
		DeviceName string    `json:"device_name"`
		CreatedAt  time.Time `json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		// 从路径参数获取user_id
		userId := utils.GetPathParam(r, 2)
		if userId == "" {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 查询用户申请的授权
		var sharedDevice []model.SharedDevice
		gormDB.Where("viewer_id = ?", userId).Find(&sharedDevice)

		// 构建返回结果
		var result []ApplicationInfo
		for _, auth := range sharedDevice {
			// 获取申请人信息
			var user model.User
			gormDB.Where("id = ?", auth.ViewerId).First(&user)

			// 获取设备信息
			var device model.Device
			gormDB.Where("id = ?", auth.DeviceId).First(&device)

			result = append(result, ApplicationInfo{
				Id:         auth.Id,
				DeviceId:   auth.DeviceId,
				Status:     auth.Authorization,
				UserName:   user.Name,
				DeviceName: device.Name,
				CreatedAt:  auth.CreatedAt,
			})
		}

		utils.Success(w, map[string]any{
			"message":      "获取申请的授权成功",
			"applications": result,
		})
	}
}

// 获取共享授权列表 GET
func GetSharedAuthorizations(db any) http.HandlerFunc {
	type AuthorizationInfo struct {
		Id         string    `json:"id"`
		Status     int       `json:"status"`
		UserName   string    `json:"user_name"`
		DeviceName string    `json:"device_name"`
		CreatedAt  time.Time `json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		// 从路径参数获取user_id
		userId := utils.GetPathParam(r, 3)
		if userId == "" {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 获取用户的所有设备ID
		var deviceIds []string
		gormDB.Model(&model.Device{}).Where("owner_id = ?", userId).Pluck("id", &deviceIds)

		// 如果没有设备, 返回空数组
		if len(deviceIds) == 0 {
			utils.Success(w, map[string]any{
				"message":        "获取共享授权列表成功",
				"authorizations": []AuthorizationInfo{},
			})
			return
		}

		// 查询这些设备的共享授权申请
		var sharedDevice []model.SharedDevice
		gormDB.Where("device_id IN ?", deviceIds).Find(&sharedDevice)

		// 构建返回结果
		var result []AuthorizationInfo
		for _, auth := range sharedDevice {
			// 获取申请人信息
			var user model.User
			gormDB.Where("id = ?", auth.ViewerId).First(&user)

			// 获取设备信息
			var device model.Device
			gormDB.Where("id = ?", auth.DeviceId).First(&device)

			result = append(result, AuthorizationInfo{
				Id:         auth.Id,
				Status:     auth.Authorization,
				UserName:   user.Name,
				DeviceName: device.Name,
				CreatedAt:  auth.CreatedAt,
			})
		}

		utils.Success(w, map[string]any{
			"message":        "获取共享授权列表成功",
			"authorizations": result,
		})
	}
}

// 授权 PUT
func AuthorizeDevice(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		var req struct {
			AccessId string `json:"id"`
			Status   int    `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 检查授权记录是否存在
		var shared model.SharedDevice
		gormDB.Where("id = ?", req.AccessId).First(&shared)
		if shared.Id == "" {
			utils.Error(w, http.StatusOK, "授权记录不存在")
			return
		}

		// 检查状态参数
		if req.Status != 1 && req.Status != 2 {
			utils.Error(w, http.StatusOK, "参数错误")
			return
		}

		// 更新授权状态
		shared.Authorization = req.Status
		gormDB.Save(&shared)

		utils.Success(w, map[string]interface{}{
			"message": "授权操作成功",
		})
	}
}

// 删除共享申请 DELETE
func DeleteShare(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		var req struct {
			AccessId string `json:"id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 检查授权记录是否存在
		var shared model.SharedDevice
		gormDB.Where("id = ?", req.AccessId).First(&shared)
		if shared.Id == "" {
			utils.Error(w, http.StatusOK, "授权记录不存在")
			return
		}

		// 删除共享申请
		gormDB.Delete(&shared)

		utils.Success(w, map[string]any{
			"message": "删除共享申请成功",
		})
	}
}
