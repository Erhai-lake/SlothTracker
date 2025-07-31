package model

import (
	"time"
)

type Device struct {
	ID           string    `gorm:"primaryKey;column:id"` // 唯一标识设备
	Name         string    `json:"name"`                 // 设备名称
	Platform     string    `json:"platform"`             // 设备平台(如: Android, iOS)
	Description  string    `json:"description"`          // 设备描述
	RegisteredAt time.Time `json:"registered_at"`        // 注册时间
}

type DeviceStatus struct {
	ID        string `gorm:"primaryKey;column:id"` // 唯一标识设备
	Timestamp int64  `json:"timestamp"`            // 上报时间戳(秒)

	// 电池与充电
	BatteryCharging    int     `json:"batteryCharging"`    // 是否充电中
	BatteryCurrent     float64 `json:"batteryCurrent"`     // 充电电流(毫安, 单位 mA, 正为充电, 负为放电)
	BatteryLevel       int     `json:"batteryLevel"`       // 电池电量百分比(0~100)
	BatteryTemperature float64 `json:"batteryTemperature"` // 电池温度(摄氏度)
	BatteryCapacity    int     `json:"batteryCapacity"`    // 电池设计容量或总容量(单位mAh, 可选)

	// 网络状态
	WifiConnected     int     `json:"wifiConnected"`     // 是否连接 WiFi
	WifiSSID          string  `json:"wifiSSID"`          // 当前连接的 WiFi 名称
	MobileDataActive  int     `json:"mobileDataActive"`  // 是否启用流量
	MobileSignalDbm   int     `json:"mobileSignalDbm"`   // 移动网络信号强度(单位 dBm)
	NetworkType       string  `json:"networkType"`       // 当前网络类型(如: WiFi, 4G, 5G, Ethernet)
	UploadSpeedKbps   int     `json:"uploadSpeedKbps"`   // 上传速度(单位 Kbps, 可选)
	DownloadSpeedKbps int     `json:"downloadSpeedKbps"` // 下载速度(单位 Kbps, 可选)
	TrafficUsedMB     float64 `json:"trafficUsedMB"`     // 当日流量使用量(单位 MB)

	// 前台状态
	AppName        string `json:"appName"`        // 当前前台应用包名
	AppTitle       string `json:"appTitle"`       // 当前应用窗口标题
	SpeakerPlaying int    `json:"speakerPlaying"` // 是否有扬声器音频播放(如音乐/视频)

	// 设备状态
	ScreenOn         int `json:"screenOn"`         // 屏幕是否点亮
	IsChargingViaUSB int `json:"isChargingViaUSB"` // 是否通过USB充电(可选)
	IsChargingViaAC  int `json:"isChargingViaAC"`  // 是否通过AC插座充电(可选)
	IsLowPowerMode   int `json:"isLowPowerMode"`   // 是否开启了省电模式(可选)
}
