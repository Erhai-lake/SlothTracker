package router

import (
	"net/http"
	"sloth-tracker/api/controller"
	"sloth-tracker/api/middleware"
	"strings"
)

func SetupRouter(db any) http.Handler {
	mux := http.NewServeMux()

	// 基础路由
	mux.HandleFunc("GET /api/ping", controller.Ping(db))

	// 用户相关路由
	mux.HandleFunc("POST /api/user/register", controller.RegisterUser(db))
	mux.HandleFunc("POST /api/user/login", controller.LoginUser(db))
	mux.HandleFunc("PUT /api/user/reset_name", controller.ResetUsername(db))
	mux.HandleFunc("PUT /api/user/reset_password", controller.ResetPassword(db))
	mux.HandleFunc("GET /api/user/query/{user_id}", controller.GetUserInfo(db))
	mux.HandleFunc("DELETE /api/user/delete", controller.DeleteUser(db))

	// 共享相关路由
	mux.HandleFunc("POST /api/share/apply", controller.ApplyShare(db))
	mux.HandleFunc("GET /api/share/query/{user_id}", controller.GetUserApplications(db))
	mux.HandleFunc("GET /api/share/authorizations/{user_id}", controller.GetSharedAuthorizations(db))
	mux.HandleFunc("PUT /api/share/authorize", controller.AuthorizeDevice(db))
	mux.HandleFunc("DELETE /api/share/delete", controller.DeleteShare(db))

	// 设备相关路由
	mux.HandleFunc("POST /api/device/register", controller.RegisterDevice(db))
	mux.HandleFunc("PUT /api/device/update", controller.UpdateDeviceInfo(db))
	mux.HandleFunc("GET /api/devices/query/{user_id}", controller.GetDeviceList(db))
	mux.HandleFunc("GET /api/devices/shared/query/{user_id}", controller.GetSharedDeviceList(db))
	mux.HandleFunc("GET /api/device/query/{device_id}", controller.GetDeviceInfo(db))
	mux.HandleFunc("DELETE /api/device/delete", controller.DeleteDevice(db))

	// 状态相关路由
	mux.HandleFunc("PUT /api/status/update/{user_id}/{device_id}", controller.UpdateStatus(db))
	mux.HandleFunc("GET /api/status/query/{user_id}/{device_id}", controller.GetStatus(db))

	// 添加中间件
	handler := middleware.CORS(mux)
	handler = middleware.Logger(handler)

	return handler
}

// 路径参数提取辅助函数
func GetPathParam(r *http.Request, param string) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	// 在路径中查找参数位置
	for i, part := range parts {
		if part == param {
			if i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}
	return ""
}
