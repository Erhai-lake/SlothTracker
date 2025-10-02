package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist/*
var assets embed.FS

func main() {
	// 创建日志文件
	logFile, err := os.OpenFile("sloth-tracker.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Printf("=== SlothTracker 启动 ===")
	log.Printf("Go 版本: %s", runtime.Version())
	log.Printf("操作系统: %s", runtime.GOOS)
	log.Printf("架构: %s", runtime.GOARCH)

	// 获取当前运行路径(软件根目录)
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("获取可执行文件路径失败:", err)
	}
	rootDir := filepath.Dir(exePath)
	log.Printf("应用根目录: %s", rootDir)

	// 设置临时目录到根目录 tmp
	tmpDir := filepath.Join(rootDir, "tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		log.Printf("创建临时目录失败: %v", err)
	} else {
		log.Printf("临时目录: %s", tmpDir)
	}
	os.Setenv("TMP", tmpDir)
	os.Setenv("TEMP", tmpDir)

	// 自定义缓存目录
	webviewDataDir := filepath.Join(rootDir, "UserData")
	if err := os.MkdirAll(webviewDataDir, 0755); err != nil {
		log.Printf("创建用户数据目录失败: %v", err)
	} else {
		log.Printf("用户数据目录: %s", webviewDataDir)
	}

	app := NewApp()
	log.Printf("应用初始化完成")

	err = wails.Run(&options.App{
		Title:  "SlothTracker",
		Width:  1024,
		Height: 768,
		MinWidth: 450,
		MinHeight: 350,
		HideWindowOnClose: true,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []any{
			app,
		},
		Windows: &windows.Options{
			WebviewUserDataPath: webviewDataDir,
		},
		// 添加错误处理
		OnDomReady: func(ctx context.Context) {
			log.Printf("DOM 准备就绪")
		},
		OnShutdown: func(ctx context.Context) {
			log.Printf("应用关闭")
		},
	})

	if err != nil {
		log.Printf("Wails 运行错误: %v", err)
		// 显示错误对话框
		showErrorDialog(err.Error())
	} else {
		log.Printf("应用正常退出")
	}
}

// 显示错误对话框
func showErrorDialog(message string) {
	// 简单的控制台错误显示
	fmt.Printf("应用程序错误: %s\n", message)
	fmt.Println("请查看 sloth-tracker.log 文件获取详细信息")
	fmt.Println("按回车键退出...")
	fmt.Scanln()
}
