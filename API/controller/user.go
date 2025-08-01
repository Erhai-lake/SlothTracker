package controller

import (
	"net/http"
	"sloth-tracker/api/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 注册用户 POST
func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		// 校验参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 检查用户名是否已存在
		var existingUser model.User
		if err := db.Where("name = ?", req.Name).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{"code": 2, "message": "用户名已存在"})
			return
		}
		userId := uuid.New().String()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 3, "message": "密码加密失败"})
			return
		}
		user := model.User{
			Id:           userId,
			Name:         req.Name,
			Password:     string(hashedPassword),
			RegisteredAt: time.Now(),
		}
		db.Create(&user)

		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "注册成功", "user_id": userId})
	}
}

// 登录用户 POST
func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		// 校验参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 检查用户名是否存在
		var user model.User
		if err := db.Where("name = ?", req.Name).First(&user).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 2, "message": "用户名或密码错误"})
			return
		}
		// 检查密码是否正确
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 3, "message": "用户名或密码错误"})
			return
		}
		// 登录成功
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "登录成功", "user_id": user.Id})
	}
}

// 重置用户名 POST
func ResetUsername(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}
		// 校验参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 检查用户名是否已存在
		var existingUser model.User
		if err := db.Where("name = ?", req.Name).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{"code": 2, "message": "用户名已存在"})
			return
		}
		// 更新用户名
		var user model.User
		if err := db.First(&user, "id = ?", req.Id).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 3, "message": "用户不存在"})
			return
		}
		user.Name = req.Name
		db.Save(&user)
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "用户名重置成功"})
	}
}

// 重置密码 POST
func ResetPassword(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Id          string `json:"id"`
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}
		// 校验参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 检查用户名是否存在
		var user model.User
		if err := db.Where("id = ?", req.Id).First(&user).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 2, "message": "用户名或密码错误"})
			return
		}
		// 检查旧密码是否正确
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 3, "message": "旧密码错误"})
			return
		}
		// 加密新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 4, "message": "新密码加密失败"})
			return
		}
		// 更新密码
		user.Password = string(hashedPassword)
		db.Save(&user)
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "密码重置成功"})
	}
}

// 获取用户信息 GET
func GetUserInfo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		// 校验参数
		if userId == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		var user model.User
		if err := db.First(&user, "id = ?", userId).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 2, "message": "用户不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "user": gin.H{
			"id":           user.Id,
			"name":         user.Name,
			"registeredAt": user.RegisteredAt,
		}})
	}
}

// 注销用户 DELETE
func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Id       string `json:"id"`
			Password string `json:"password"`
		}
		// 校验参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		// 检查用户是否存在
		var user model.User
		if err := db.First(&user, "id = ?", req.Id).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 2, "message": "用户不存在"})
			return
		}
		// 检查密码是否正确
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 3, "message": "密码错误"})
			return
		}
		// 启动事务
		tx := db.Begin()
		// 删除用户
		if err := tx.Delete(&user).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 4, "message": "用户注销失败"})
			return
		}
		// 删除用户相关的设备和状态
		// if err := tx.Where("user_id = ?", req.Id).Delete(&model.Device{}).Error; err != nil {
		// 	tx.Rollback()
		// 	c.JSON(http.StatusInternalServerError, gin.H{"code": 5, "message": "用户注销失败"})
		// 	return
		// }
		// if err := tx.Where("user_id = ?", req.Id).Delete(&model.DeviceStatus{}).Error; err != nil {
		// 	tx.Rollback()
		// 	c.JSON(http.StatusInternalServerError, gin.H{"code": 6, "message": "用户注销失败"})
		// 	return
		// }
		// 提交事务
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "用户注销成功"})
	}
}
