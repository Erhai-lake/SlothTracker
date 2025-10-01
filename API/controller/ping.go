package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// Ping 测试接口, 返回延迟时间
func Ping(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		time.Sleep(10 * time.Millisecond)
		latency := time.Since(start).Milliseconds()
		latencyStr := strconv.FormatInt(latency, 10) + "ms"
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": latencyStr,
		})
	}
}
