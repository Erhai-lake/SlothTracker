package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

// JSONResponse 统一的JSON响应
func JSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Success 成功响应
func Success(w http.ResponseWriter, data any) {
	JSONResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    data,
	})
}

// Error 错误响应
func Error(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]any{
		"success": false,
		"error":   message,
	})
}

// 从请求中获取路径参数
func GetPathParam(r *http.Request, position int) string {
	path := r.URL.Path
	parts := make([]string, 0)

	// 分割路径, 忽略空字符串
	for part := range strings.SplitSeq(path, "/") {
		if part != "" {
			parts = append(parts, part)
		}
	}

	if position < len(parts) {
		return parts[position]
	}
	return ""
}
