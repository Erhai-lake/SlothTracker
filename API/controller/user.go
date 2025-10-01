package controller

import (
	"encoding/json"
	"net/http"
	"sloth-tracker/api/model"
	"sloth-tracker/api/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 注册用户 POST
func RegisterUser(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		var req struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}

		// 解析JSON参数
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 检查用户名是否已存在
		var existingUser model.User
		if err := gormDB.Where("name = ?", req.Name).First(&existingUser).Error; err == nil {
			utils.Error(w, http.StatusOK, "用户名已存在")
			return
		}

		userId := uuid.New().String()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "密码加密失败")
			return
		}

		user := model.User{
			Id:           userId,
			Name:         req.Name,
			Password:     string(hashedPassword),
			RegisteredAt: time.Now(),
		}
		gormDB.Create(&user)

		utils.Success(w, map[string]any{
			"message": "注册成功",
			"user_id": userId,
		})
	}
}

// 登录用户 POST
func LoginUser(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		var req struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 检查用户名是否存在
		var user model.User
		if err := gormDB.Where("name = ?", req.Name).First(&user).Error; err != nil {
			utils.Error(w, http.StatusOK, "用户名或密码错误")
			return
		}

		// 检查密码是否正确
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			utils.Error(w, http.StatusOK, "用户名或密码错误")
			return
		}

		// 登录成功
		utils.Success(w, map[string]any{
			"message": "登录成功",
			"user_id": user.Id,
		})
	}
}

// 重置用户名 PUT
func ResetUsername(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		var req struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 检查用户名是否已存在
		var existingUser model.User
		if err := gormDB.Where("name = ?", req.Name).First(&existingUser).Error; err == nil {
			utils.Error(w, http.StatusOK, "用户名已存在")
			return
		}

		// 更新用户名
		var user model.User
		if err := gormDB.First(&user, "id = ?", req.Id).Error; err != nil {
			utils.Error(w, http.StatusOK, "用户不存在")
			return
		}

		user.Name = req.Name
		gormDB.Save(&user)

		utils.Success(w, map[string]any{
			"message": "用户名重置成功",
		})
	}
}

// 重置密码 PUT
func ResetPassword(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		var req struct {
			Id          string `json:"id"`
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 检查用户是否存在
		var user model.User
		if err := gormDB.Where("id = ?", req.Id).First(&user).Error; err != nil {
			utils.Error(w, http.StatusOK, "用户名或密码错误")
			return
		}

		// 检查旧密码是否正确
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			utils.Error(w, http.StatusOK, "旧密码错误")
			return
		}

		// 检查新密码是否和旧密码一致
		if req.OldPassword == req.NewPassword {
			utils.Error(w, http.StatusOK, "新密码不能和旧密码相同")
			return
		}

		// 加密新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "新密码加密失败")
			return
		}

		// 更新密码
		user.Password = string(hashedPassword)
		gormDB.Save(&user)

		utils.Success(w, map[string]any{
			"message": "密码重置成功",
		})
	}
}

// 获取用户信息 GET
func GetUserInfo(db any) http.HandlerFunc {
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

		var user model.User
		if err := gormDB.First(&user, "id = ?", userId).Error; err != nil {
			utils.Error(w, http.StatusOK, "用户不存在")
			return
		}

		utils.Success(w, map[string]any{
			"message": "获取用户信息成功",
			"user": map[string]any{
				"id":            user.Id,
				"name":          user.Name,
				"registered_at": user.RegisteredAt,
			},
		})
	}
}

// 注销用户 DELETE
func DeleteUser(db any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "方法不允许")
			return
		}

		var req struct {
			Id       string `json:"id"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "参数错误")
			return
		}

		gormDB := db.(*gorm.DB)

		// 检查用户是否存在
		var user model.User
		if err := gormDB.First(&user, "id = ?", req.Id).Error; err != nil {
			utils.Error(w, http.StatusOK, "用户不存在")
			return
		}

		// 检查密码是否正确
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			utils.Error(w, http.StatusOK, "密码错误")
			return
		}

		// 启动事务
		tx := gormDB.Begin()

		// 删除用户
		if err := tx.Delete(&user).Error; err != nil {
			tx.Rollback()
			utils.Error(w, http.StatusInternalServerError, "用户注销失败-删除用户失败")
			return
		}

		// 获取用户有关的所有设备ID
		var deviceIds []string
		if err := tx.Model(&model.Device{}).
			Where("owner_id = ?", req.Id).
			Pluck("id", &deviceIds).Error; err != nil {
			tx.Rollback()
			utils.Error(w, http.StatusInternalServerError, "用户注销失败-获取设备ID失败")
			return
		}

		// 删除用户所有设备状态
		if len(deviceIds) > 0 {
			if err := tx.Where("id IN ?", deviceIds).Delete(&model.DeviceStatus{}).Error; err != nil {
				tx.Rollback()
				utils.Error(w, http.StatusInternalServerError, "用户注销失败-删除设备状态失败")
				return
			}

			// 删除用户所有设备
			if err := tx.Where("id IN ?", deviceIds).Delete(&model.Device{}).Error; err != nil {
				tx.Rollback()
				utils.Error(w, http.StatusInternalServerError, "用户注销失败-删除设备失败")
				return
			}
		}

		// 提交事务
		tx.Commit()

		utils.Success(w, map[string]any{
			"message": "用户注销成功",
		})
	}
}
