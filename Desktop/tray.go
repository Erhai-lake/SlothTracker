package main

import (
	"context"
	_ "embed"
	"log"

	"github.com/getlantern/systray"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 嵌入图标文件

//go:embed assets/icon.ico
var trayIcon []byte

type TrayManager struct {
	ctx context.Context
	app *App
}

func NewTrayManager(app *App) *TrayManager {
	return &TrayManager{
		app: app,
	}
}

func (t *TrayManager) Startup(ctx context.Context) {
	t.ctx = ctx
}

// 启动托盘
func (t *TrayManager) StartTray() {
	// 在协程中启动托盘
	go func() {
		systray.Run(t.onReady, t.onExit)
	}()
}

// 托盘准备就绪
func (t *TrayManager) onReady() {
	// 设置托盘图标
	iconData := getTrayIcon()
	if iconData != nil {
		systray.SetIcon(iconData)
	}

	// 设置托盘标题和提示
	systray.SetTitle("SlothTracker")
	systray.SetTooltip("SlothTracker - 树懒")

	// 添加菜单项
	showWindow := systray.AddMenuItem("显示窗口", "显示主窗口")
	systray.AddSeparator()
	quit := systray.AddMenuItem("退出", "退出应用程序")

	// 监听菜单点击事件
	go func() {
		for {
			select {
			case <-showWindow.ClickedCh:
				t.showWindow()
			case <-quit.ClickedCh:
				t.quitApp()
				return
			}
		}
	}()
}

// 托盘退出
func (t *TrayManager) onExit() {
	log.Println("托盘退出")
}

// 显示窗口
func (t *TrayManager) showWindow() {
	if t.ctx != nil && t.app != nil {
		t.app.ShowWindow()
	}
}

// 隐藏窗口
func (t *TrayManager) quitApp() {
	log.Println("用户从托盘退出应用")
	systray.Quit()
	if t.ctx != nil {
		wailsRuntime.Quit(t.ctx)
	}
}

// 获取托盘图标数据
func getTrayIcon() []byte {
	return trayIcon
}
