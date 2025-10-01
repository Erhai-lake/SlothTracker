package main

import (
	"runtime"
	"runtime/debug"
	"sloth-tracker/api/router"
	"sloth-tracker/api/storage"
)

func main() {
	// 限制CPU使用(1核)
	runtime.GOMAXPROCS(1)
	// 设置内存限制(500MB)
	debug.SetMemoryLimit(500 * 1024 * 1024)

	db := storage.InitDB()
	r := router.SetupRouter(db)
	r.Run(":8080")
}
