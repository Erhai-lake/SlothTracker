package model

import (
	"time"
)

type User struct {
	Id           string    `gorm:"primaryKey;column:id" json:"id"` // 用户ID
	Name         string    `json:"name"`                           // 用户名
	Password     string    `json:"password"`                       // 密码
	RegisteredAt time.Time `json:"registered_at"`                  // 注册时间
}

type SharedDevice struct {
	Id            string    `gorm:"primaryKey;column:id" json:"id"` // 唯一标识
	DeviceId      string    `json:"device_id"`                      // 被访问的设备
	ViewerId      string    `json:"viewer_id"`                      // 被授权的用户ID
	Authorization int       `json:"authorization"`                  // 是否授权(1: 已授权, 2: 未授权)
	CreatedAt     time.Time `json:"created_at"`                     // 创建时间
}

type Device struct {
	Id           string    `gorm:"primaryKey;column:id" json:"id"` // 设备ID
	OwnerId      string    `json:"owner_id"`                       // 所属用户ID
	Name         string    `json:"name"`                           // 设备名称
	Platform     string    `json:"platform"`                       // 设备平台(如: Android, iOS)
	Description  string    `json:"description"`                    // 设备描述
	RegisteredAt time.Time `json:"registered_at"`                  // 注册时间
}

type DeviceStatus struct {
	Id         string           `gorm:"primaryKey;column:id" json:"id"`                        // 唯一标识设备
	DeviceId   string           `json:"device_id"`                                             // 设备ID
	Timestamp  int64            `json:"timestamp"`                                             // 上报时间戳(秒)
	Battery    BatteryStatus    `gorm:"embedded;embeddedPrefix:battery_" json:"battery"`       // 电池状态
	Network    NetworkStatus    `gorm:"embedded;embeddedPrefix:network_" json:"network"`       // 网络状态
	Foreground ForegroundStatus `gorm:"embedded;embeddedPrefix:foreground_" json:"foreground"` // 前台应用状态
	Other      OtherStatus      `gorm:"embedded;embeddedPrefix:other_" json:"other"`           // 其他状态
}

type BatteryStatus struct {
	Charging    int     `json:"charging"`    // 是否充电中(1: 充电中, 2: 未充电, 3: 已充满)
	Level       int     `json:"level"`       // 电池电量百分比(0~100)
	Temperature float64 `json:"temperature"` // 电池温度(摄氏度)
	Capacity    int     `json:"capacity"`    // 电池设计容量或总容量(单位mAh, 可选)
}

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

type ForegroundStatus struct {
	AppName        string `json:"app_name"`        // 当前前台应用包名
	AppTitle       string `json:"app_title"`       // 当前应用窗口标题
	SpeakerPlaying int    `json:"speaker_playing"` // 是否有扬声器音频播放(1: 播放, 2: 未播放, 3: 未知)
}

type OtherStatus struct {
	ScreenOn         int `json:"screen_on"`           // 屏幕是否点亮(1: 点亮, 2: 未点亮)
	IsChargingViaUSB int `json:"is_charging_via_usb"` // 是否通过USB充电(1: 是, 2: 否)
	IsChargingViaAC  int `json:"is_charging_via_ac"`  // 是否通过AC插座充电(1: 是, 2: 否)
	IsLowPowerMode   int `json:"is_low_power_mode"`   // 是否开启了省电模式(1: 开启, 2: 未开启)
}
