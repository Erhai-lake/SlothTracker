package status

import (
	"encoding/hex"
	"errors"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	psnet "github.com/shirou/gopsutil/v3/net"
)

type NetworkStatus struct {
	WifiConnected     int     `json:"wifi_connected"`      // 是否连接 WiFi(1: 连接, 2: 未连接)
	WifiSSId          string  `json:"wifi_ssid"`           // 当前连接的 WiFi 名称
	MobileDataActive  int     `json:"mobile_data_active"`  // 是否启用流量(1: 启用, 2: 未启用)
	MobileSignalDbm   int     `json:"mobile_signal_dbm"`   // 移动网络信号强度(单位 dBm)
	NetworkType       string  `json:"network_type"`        // 当前网络类型(如: WiFi, 4G, 5G, Ethernet)
	TrafficUsedMB     float64 `json:"traffic_used_mb"`     // 当日流量使用量(单位 MB)
	UploadSpeedKbps   int     `json:"upload_speed_kbps"`   // 上传速度(单位 Kbps, 可选)
	DownloadSpeedKbps int     `json:"download_speed_kbps"` // 下载速度(单位 Kbps, 可选)
}

func GetWifiInfo() (*NetworkStatus, error) {
	status := &NetworkStatus{
		MobileDataActive: 3,
		MobileSignalDbm:  -1,
		TrafficUsedMB:    0,
	}
	status.NetworkType = detectNetworkType()
	var platformInfo *NetworkStatus
	var err error
	switch runtime.GOOS {
	case "windows":
		platformInfo, err = getWindowsWifiInfo()
	case "darwin":
		platformInfo, err = getMacWifiInfo()
	case "linux":
		platformInfo, err = getLinuxWifiInfo()
	default:
		return nil, errors.New("unsupported platform")
	}
	if err != nil {
		return nil, err
	}
	status.WifiConnected = platformInfo.WifiConnected
	status.WifiSSId = platformInfo.WifiSSId
	status.NetworkType = detectNetworkType()
	speedInfo, err := GetSpeedInfo()
	if err == nil {
		status.UploadSpeedKbps = speedInfo.UploadSpeedKbps
		status.DownloadSpeedKbps = speedInfo.DownloadSpeedKbps
	}
	return status, nil
}

func detectNetworkType() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "unknown"
	}
	for _, i := range interfaces {
		// 只检查已启动且非回环的接口
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 转换为小写方便比较
		name := strings.ToLower(i.Name)
		switch {
		case isWifiInterface(name):
			return "wifi"
		case isEthernetInterface(name):
			return "ethernet"
		}
	}
	return "unknown"
}

func isWifiInterface(name string) bool {
	// 常见WiFi接口名称模式
	wifiPatterns := []string{
		"wi",   // 通用匹配 (wifi, wireless)
		"wlan", // Linux常见
		"wlp",  // Linux systemd命名
		"ath",  // Atheros无线网卡
		"ra",   // Ralink无线网卡
		"en",   // macOS可能用于WiFi (en0, en1)
		"awdl", // Apple Wireless Direct Link
		"p2p",  // WiFi Direct
		"anpi", // Android
	}
	for _, pattern := range wifiPatterns {
		if strings.Contains(name, pattern) {
			return true
		}
	}
	return false
}

func isEthernetInterface(name string) bool {
	// 常见有线网接口名称模式
	ethernetPatterns := []string{
		"eth",    // Linux传统命名
		"enp",    // Linux systemd命名 (PCI)
		"ens",    // Linux systemd命名 (SATA)
		"enx",    // Linux systemd命名 (MAC地址)
		"em",     // Intel网卡
		"p0p",    // PCI
		"en",     // macOS有线网 (en0, en1)
		"usb",    // USB以太网适配器
		"lan",    // 通用
		"local",  // 通用
		"bridge", // 桥接
		"veth",   // 虚拟以太网
	}
	for _, pattern := range ethernetPatterns {
		if strings.Contains(name, pattern) {
			return true
		}
	}
	return false
}

// 尝试将十六进制字符串转为UTF-8字符串, 如果失败就返回原始字符串
func tryDecodeSSIDHex(s string) string {
	// 长度必须是偶数, 且必须是合法十六进制字符
	if len(s)%2 != 0 {
		return s
	}
	data, err := hex.DecodeString(s)
	if err != nil {
		return s
	}
	if !utf8.Valid(data) {
		return s
	}
	return string(data)
}

func getWindowsWifiInfo() (*NetworkStatus, error) {
	// 获取当前连接的SSID
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	// 隐藏 CMD 窗口 (Windows)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	info := &NetworkStatus{WifiConnected: 2}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "SSID") && !strings.Contains(line, "BSSID") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				rawSSID := strings.TrimSpace(parts[1])
				ssid := tryDecodeSSIDHex(rawSSID)
				if ssid != "" {
					info.WifiConnected = 1
					info.WifiSSId = ssid
					break
				}
			}
		}
	}
	return info, nil
}

func getMacWifiInfo() (*NetworkStatus, error) {
	cmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-I")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	info := &NetworkStatus{WifiConnected: 2}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "SSID:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				ssid := strings.TrimSpace(parts[1])
				if ssid != "" {
					info.WifiConnected = 1
					info.WifiSSId = ssid
					break
				}
			}
		}
	}
	return info, nil
}

func getLinuxWifiInfo() (*NetworkStatus, error) {
	cmd := exec.Command("iwgetid", "-r")
	output, err := cmd.Output()
	if err != nil {
		// 可能是没有连接WiFi或者iwgetid不存在
		return &NetworkStatus{WifiConnected: 2}, nil
	}
	ssid := strings.TrimSpace(string(output))
	if ssid == "" {
		return &NetworkStatus{WifiConnected: 2}, nil
	}
	return &NetworkStatus{
		WifiConnected: 1,
		WifiSSId:      ssid,
	}, nil
}

func GetSpeedInfo() (*NetworkStatus, error) {
	// 获取初始网络IO计数
	ioCountersStart, err := psnet.IOCounters(false)
	if err != nil || len(ioCountersStart) == 0 {
		return nil, err
	}
	time.Sleep(1 * time.Second) // 采样时间间隔
	// 获取结束时的网络IO计数
	ioCountersEnd, err := psnet.IOCounters(false)
	if err != nil || len(ioCountersEnd) == 0 {
		return nil, err
	}
	start := ioCountersStart[0]
	end := ioCountersEnd[0]
	// 计算字节差值
	bytesSent := end.BytesSent - start.BytesSent
	bytesRecv := end.BytesRecv - start.BytesRecv
	// 转换为 Kbps (Bytes/sec → bits/sec → kilobits/sec)
	uploadKbps := int((bytesSent * 8) / 1024)
	downloadKbps := int((bytesRecv * 8) / 1024)
	return &NetworkStatus{
		UploadSpeedKbps:   uploadKbps,
		DownloadSpeedKbps: downloadKbps,
	}, nil
}
