package status

import (
	"syscall"
	"unsafe"
)

type OtherInfo struct {
	ScreenOn         int `json:"screenOn"`         // 屏幕是否点亮
	IsChargingViaUsb int `json:"isChargingViaUSB"` // 是否通过USB充电
	IsChargingViaAC  int `json:"isChargingViaAC"`  // 是否通过AC充电
	IsLowPowerMode   int `json:"isLowPowerMode"`   // 是否省电模式
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

func GetOtherInfo() (OtherInfo, error) {
	var status systemPowerStatus
	ret, _, err := procGetSystemPowerStatus.Call(uintptr(unsafe.Pointer(&status)))
	if ret == 0 {
		return OtherInfo{}, err
	}
	info := OtherInfo{}
	// 供电状态: 0=离线, 1=在线(AC), 255=未知
	if status.ACLineStatus == 1 {
		info.IsChargingViaAC = 2
	} else {
		info.IsChargingViaAC = 1
	}
	// 判断是否为 USB 充电(不精确, Windows 无法直接识别 USB)
	if status.ACLineStatus == 1 && status.BatteryFlag&8 != 0 {
		info.IsChargingViaUsb = 2
	} else {
		info.IsChargingViaUsb = 1
	}
	// 省电模式: BatteryFlag 中 1=High, 2=Low, 4=Critical, 8=Charging, 128=No system battery
	if status.BatteryFlag&2 != 0 || status.BatteryFlag&4 != 0 {
		info.IsLowPowerMode = 2
	}
	// 屏幕状态 Windows 无标准 API, 只能模拟或调用 DWM/PowerSettings API
	// 默认假设屏幕亮着
	info.ScreenOn = 2
	return info, nil
}
