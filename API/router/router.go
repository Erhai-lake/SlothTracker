package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sloth-tracker/api/controller"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 添加CORS配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ping
	r.GET("/api/ping", controller.Ping(db))

	// 用户注册
	r.POST("/api/user/register", controller.RegisterUser(db))
	// 用户登录
	r.POST("/api/user/login", controller.LoginUser(db))
	// 重置用户名
	r.PUT("/api/user/reset_name", controller.ResetUsername(db))
	// 重置密码
	r.PUT("/api/user/reset_password", controller.ResetPassword(db))
	// 获取用户信息
	r.GET("/api/user/:user_id", controller.GetUserInfo(db))
	// 注销用户
	r.DELETE("/api/user/delete", controller.DeleteUser(db))

	// 共享申请
	r.POST("/api/share/apply", controller.ApplyShare(db))
	// 获取用户申请的授权
	r.GET("/api/share/:user_id", controller.GetUserApplications(db))
	// 获取共享授权列表
	r.GET("/api/share/authorizations/:user_id", controller.GetSharedAuthorizations(db))
	// 授权设备
	r.PUT("/api/share/authorize", controller.AuthorizeDevice(db))
	// 删除共享申请
	r.DELETE("/api/share/delete", controller.DeleteShare(db))

	// 设备注册
	r.POST("/api/device/register", controller.RegisterDevice(db))
	// 修改设备信息
	r.PUT("/api/device/update", controller.UpdateDeviceInfo(db))
	// 获取设备列表
	r.GET("/api/devices/:user_id", controller.GetDeviceList(db))
	// 获取共享设备列表
	r.GET("/api/devices/shared/:user_id", controller.GetSharedDeviceList(db))
	// 获取设备信息
	r.GET("/api/device/:device_id", controller.GetDeviceInfo(db))
	// 删除设备
	r.DELETE("/api/device/delete", controller.DeleteDevice(db))

	// 状态更新
	r.PUT("/api/status/update/:device_id", controller.UpdateStatus(db))
	// 获取设备状态
	r.GET("/api/status/:device_id", controller.GetStatus(db))

	return r
}
