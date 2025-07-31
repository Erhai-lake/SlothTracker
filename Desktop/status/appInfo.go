package status

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type AppInfo struct {
	AppName        string `json:"appName"`        // 应用名称
	AppTitle       string `json:"appTitle"`       // 应用标题
	SpeakerPlaying int    `json:"speakerPlaying"` // 扬声器是否正在播放
}

var (
	user32                       = windows.NewLazySystemDLL("user32.dll")
	kernel32                     = windows.NewLazySystemDLL("kernel32.dll")
	psapi                        = windows.NewLazySystemDLL("psapi.dll")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procOpenProcess              = kernel32.NewProc("OpenProcess")
	procGetModuleBaseNameW       = psapi.NewProc("GetModuleBaseNameW")
	procCloseHandle              = kernel32.NewProc("CloseHandle")
)

const PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
const PROCESS_QUERY_INFORMATION = 0x0400
const PROCESS_VM_READ = 0x0010

func GetAppInfo() (AppInfo, error) {
	switch runtime.GOOS {
	case "windows":
		return getWindowsAppInfo()
	case "darwin":
		return getMacAppInfo()
	case "linux":
		return getLinuxAppInfo()
	default:
		return AppInfo{}, errors.New("unsupported platform")
	}
}

func getWindowsAppInfo() (AppInfo, error) {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return AppInfo{}, errors.New("无法获取前台窗口句柄")
	}

	// 获取窗口标题
	var buf [256]uint16
	ret, _, _ := procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	windowTitle := "无法获取窗口标题"
	windowTitle = syscall.UTF16ToString(buf[:ret])

	// 获取进程ID
	var pid uint32
	_, _, _ = procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid == 0 {
		return AppInfo{}, errors.New("无法获取进程ID")
	}

	// 打开进程句柄
	handle, _, _ := procOpenProcess.Call(PROCESS_QUERY_INFORMATION|PROCESS_VM_READ, 0, uintptr(pid))
	if handle == 0 {
		errCode := windows.GetLastError()
		return AppInfo{}, fmt.Errorf("无法打开进程句柄, 错误码: %d", errCode)
	}
	defer procCloseHandle.Call(handle)

	// 获取进程名
	var procNameBuf [260]uint16
	ret, _, _ = procGetModuleBaseNameW.Call(handle, 0, uintptr(unsafe.Pointer(&procNameBuf[0])), uintptr(len(procNameBuf)))
	if ret == 0 {
		return AppInfo{}, errors.New("无法获取进程名")
	}
	processName := syscall.UTF16ToString(procNameBuf[:ret])

	return AppInfo{
		AppName:        processName,
		AppTitle:       windowTitle,
		SpeakerPlaying: -1,
	}, nil
}

func getMacAppInfo() (AppInfo, error) {
	// 使用osascript获取前台应用名和窗口标题
	cmdName := exec.Command("osascript", "-e",
		`tell application "System Events" to get name of first application process whose frontmost is true`)
	outName, err := cmdName.Output()
	if err != nil {
		return AppInfo{}, err
	}
	appName := strings.TrimSpace(string(outName))

	cmdTitle := exec.Command("osascript", "-e",
		`tell application "System Events" to get value of attribute "AXTitle" of front window of first application process whose frontmost is true`)
	outTitle, err := cmdTitle.Output()
	if err != nil {
		return AppInfo{}, err
	}
	appTitle := strings.TrimSpace(string(outTitle))

	return AppInfo{
		AppName:        appName,
		AppTitle:       appTitle,
		SpeakerPlaying: -1,
	}, nil
}

func getLinuxAppInfo() (AppInfo, error) {
	// 依赖xdotool, 获取活动窗口的窗口类和标题
	cmdName := exec.Command("xdotool", "getactivewindow", "getwindowname")
	outTitle, err := cmdName.Output()
	if err != nil {
		return AppInfo{}, err
	}
	title := strings.TrimSpace(string(outTitle))

	// 取窗口标题做为appTitle, appName这里简单赋值为unknown
	return AppInfo{
		AppName:        "unknown",
		AppTitle:       title,
		SpeakerPlaying: -1,
	}, nil
}
