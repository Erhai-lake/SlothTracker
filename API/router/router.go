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

	// 设备注册
	r.POST("/api/register", controller.Register(db))
	// 状态更新
	r.POST("/api/update/:device_id", controller.UpdateStatus(db))
	// 删除设备
	r.DELETE("/api/delete/:device_id", controller.DeleteDevice(db))
	// 获取设备列表
	r.GET("/api/devices", controller.ListDevice(db))
	// 获取状态列表
	r.GET("/api/list", controller.ListStatus(db))
	// 获取设备状态
	r.GET("/api/status/:device_id", controller.GetStatus(db))

	return r
}
