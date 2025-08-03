package status

import (
	"syscall"
	"unsafe"
)

type OtherStatus struct {
	ScreenOn         int `json:"screen_on"`           // 屏幕是否点亮(1: 点亮, 2: 未点亮)
	IsChargingViaUSB int `json:"is_charging_via_usb"` // 是否通过USB充电(1: 是, 2: 否)
	IsChargingViaAC  int `json:"is_charging_via_ac"`  // 是否通过AC插座充电(1: 是, 2: 否)
	IsLowPowerMode   int `json:"is_low_power_mode"`   // 是否开启了省电模式(1: 开启, 2: 未开启)
}

type systemPowerStatus struct {
	ACLineStatus        byte
	BatteryFlag         byte
	BatteryLifePercent  byte
	Reserved1           byte
	BatteryLifeTime     uint32
	BatteryFullLifeTime uint32
}

var (
	modkernel32              = syscall.NewLazyDLL("kernel32.dll")
	procGetSystemPowerStatus = modkernel32.NewProc("GetSystemPowerStatus")
)

func GetOtherInfo() (OtherStatus, error) {
	var status systemPowerStatus
	ret, _, err := procGetSystemPowerStatus.Call(uintptr(unsafe.Pointer(&status)))
	if ret == 0 {
		return OtherStatus{}, err
	}
	info := OtherStatus{}
	// 供电状态: 0=离线, 1=在线(AC), 255=未知
	if status.ACLineStatus == 1 {
		info.IsChargingViaAC = 1
	} else {
		info.IsChargingViaAC = 2
	}
	// 判断是否为 USB 充电(不精确, Windows 无法直接识别 USB)
	if status.ACLineStatus == 1 && status.BatteryFlag&8 != 0 {
		info.IsChargingViaUSB = 1
	} else {
		info.IsChargingViaUSB = 2
	}
	// 省电模式: BatteryFlag 中 1=High, 2=Low, 4=Critical, 8=Charging, 128=No system battery
	if status.BatteryFlag&2 != 0 || status.BatteryFlag&4 != 0 {
		info.IsLowPowerMode = 1
	} else {
		info.IsLowPowerMode = 2
	}
	// 屏幕状态 Windows 无标准 API, 只能模拟或调用 DWM/PowerSettings API
	// 默认假设屏幕亮着
	info.ScreenOn = 1
	return info, nil
}
