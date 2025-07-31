package status

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

type WifiInfo struct {
	WifiConnected int    `json:"wifiConnected"` // 1: 否, 2: 是
	WifiSSID      string `json:"wifiSSID"`      // SSID 名称
}

func GetWifiInfo() (*WifiInfo, error) {
	switch runtime.GOOS {
	case "windows":
		return getWindowsWifiInfo()
	case "darwin":
		return getMacWifiInfo()
	case "linux":
		return getLinuxWifiInfo()
	default:
		return nil, errors.New("unsupported platform")
	}
}

func getWindowsWifiInfo() (*WifiInfo, error) {
	// 获取当前连接的SSID
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	info := &WifiInfo{WifiConnected: 1}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "SSID") && !strings.Contains(line, "BSSID") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				ssid := strings.TrimSpace(parts[1])
				if ssid != "" {
					info.WifiConnected = 2
					info.WifiSSID = ssid
					break
				}
			}
		}
	}
	return info, nil
}

func getMacWifiInfo() (*WifiInfo, error) {
	cmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-I")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	info := &WifiInfo{WifiConnected: 1}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "SSID:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				ssid := strings.TrimSpace(parts[1])
				if ssid != "" {
					info.WifiConnected = 2
					info.WifiSSID = ssid
					break
				}
			}
		}
	}
	return info, nil
}

func getLinuxWifiInfo() (*WifiInfo, error) {
	cmd := exec.Command("iwgetid", "-r")
	output, err := cmd.Output()
	if err != nil {
		// 可能是没有连接WiFi或者iwgetid不存在
		return &WifiInfo{WifiConnected: 1}, nil
	}
	ssid := strings.TrimSpace(string(output))
	if ssid == "" {
		return &WifiInfo{WifiConnected: 1}, nil
	}
	return &WifiInfo{
		WifiConnected: 2,
		WifiSSID:      ssid,
	}, nil
}
