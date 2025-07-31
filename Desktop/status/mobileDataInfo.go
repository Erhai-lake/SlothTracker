package status

import (
	"net"
	"strings"
)

type MobileDataInfo struct {
	MobileDataActive int    `json:"mobileDataActive"` // 1: 否, 2: 是
	MobileSignalDbm  int    `json:"mobileSignalDbm"`  // 信号强度
	NetworkType      string `json:"networkType"`      // 网络类型
	TrafficUsedMB    int    `json:"trafficUsedMB"`    // 流量使用量
}

func GetMobileDataInfo() MobileDataInfo {
	networkType := detectNetworkType()
	return MobileDataInfo{
		MobileDataActive: 1,
		MobileSignalDbm:  -1,
		NetworkType:      networkType,
		TrafficUsedMB:    -1,
	}
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

// isWifiInterface 判断是否是WiFi接口
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

// isEthernetInterface 判断是否是有线网接口
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
